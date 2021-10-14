package helm

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/hex-techs/klog"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"go-arsenal.kanzhun.tech/arsenal/go-kit/retry"
)

func (c *HelmClient) repoAdd(name, url string) (string, string, error) {
	repoFile := settings.RepositoryConfig
	cache := settings.RepositoryCache

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return repoFile, cache, err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return repoFile, cache, err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return repoFile, cache, err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return repoFile, cache, err
	}

	if f.Has(name) {
		klog.Infof("repository name (%s) already exists", name)
		return repoFile, cache, nil
	}

	e := repo.Entry{
		Name:                  name,
		URL:                   url,
		InsecureSkipTLSverify: true,
	}

	r, err := repo.NewChartRepository(&e, getter.All(settings))
	if err != nil {
		return repoFile, cache, err
	}

	run := func() error {
		if _, err := r.DownloadIndexFile(); err != nil {
			return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		}
		return nil
	}

	err = retry.RetryFunc(5, time.Second, run)
	if err != nil {
		return repoFile, cache, err
	}

	f.Update(&e)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		klog.Error(err)
		os.Exit(3)
	}

	klog.Infof("%q has been added to your repositories", name)
	return repoFile, cache, nil
}

var errNoRepositories = errors.New("no repositories found. You must add one before updating")

type repoUpdateOptions struct {
	update               func([]*repo.ChartRepository) error
	repoFile             string
	repoCache            string
	names                []string
	failOnRepoUpdateFail bool
}

func (o *repoUpdateOptions) run() error {
	f, err := repo.LoadFile(o.repoFile)
	switch {
	case isNotExist(err):
		return errNoRepositories
	case err != nil:
		return errors.Wrapf(err, "failed loading file: %s", o.repoFile)
	case len(f.Repositories) == 0:
		return errNoRepositories
	}

	var repos []*repo.ChartRepository
	updateAllRepos := len(o.names) == 0

	if !updateAllRepos {
		// Fail early if the user specified an invalid repo to update
		if err := checkRequestedRepos(o.names, f.Repositories); err != nil {
			return err
		}
	}

	for _, cfg := range f.Repositories {
		if updateAllRepos || isRepoRequested(cfg.Name, o.names) {
			r, err := repo.NewChartRepository(cfg, getter.All(settings))
			if err != nil {
				return err
			}
			if o.repoCache != "" {
				r.CachePath = o.repoCache
			}
			repos = append(repos, r)
		}
	}

	return o.update(repos)
}

func updateCharts(repos []*repo.ChartRepository) error {
	klog.Info("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	var repoFailList []string
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			klog.Infof("updating repo (%s) %s...", re.Config.Name, re.Config.URL)
			if _, err := re.DownloadIndexFile(); err != nil {
				klog.Errorf("...Unable to get an update from the %q chart repository (%s): %v", re.Config.Name, re.Config.URL, err)
				repoFailList = append(repoFailList, re.Config.URL)
			} else {
				klog.Infof("...Successfully got an update from the %q chart repository", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()

	if len(repoFailList) > 0 {
		klog.Warningf("Failed to update the following repositories: %v", repoFailList)
	}

	klog.Info("Update Complete. ⎈Happy Helming!⎈")
	return nil
}
