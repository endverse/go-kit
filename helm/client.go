package helm

import (
	"fmt"
	"log"
	"os"

	"github.com/hex-techs/klog"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

var (
	repoUrl  = "https://harbor.weizhipin.com/chartrepo/arsenal-ai"
	repoName = "arsenal-ai"
)

const (
	defaultCachePath            = "/tmp/.helmcache"
	defaultRepositoryConfigPath = "/tmp/.helmrepo"
)

type HelmClient struct {
	cfg       *action.Configuration
	list      *action.List
	linting   bool
	uninstall *action.Uninstall
}

var settings = cli.New()

func logFunc(format string, v ...interface{}) {
	log.SetFlags(log.Llongfile)
	format = fmt.Sprintf("[helm client] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}

func NewHelmClient(namespace string) *HelmClient {
	actionConfig := new(action.Configuration)
	actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), logFunc)

	hc := &HelmClient{
		cfg:       actionConfig,
		list:      action.NewList(actionConfig),
		linting:   true,
		uninstall: action.NewUninstall(actionConfig),
	}

	repoFile, cache, err := hc.RepoAdd(getRepoName(), getRepoUrl())
	if err != nil {
		klog.Fatal(err)
	}

	settings.RepositoryConfig = repoFile
	settings.RepositoryCache = cache

	return hc
}

func (c *HelmClient) List(labelMap map[string]string) ([]*release.Release, error) {
	c.list.Selector = MapToString(labelMap)

	return c.list.Run()
}

func (c *HelmClient) Install(chart, version string, values map[string]interface{}) (string, string, error) {
	return c.install(chart, version, values)
}

func (c *HelmClient) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	return c.uninstall.Run(name)
}

func (c *HelmClient) RepoAdd(name, url string) (string, string, error) {
	return c.repoAdd(name, url)
}

func (c *HelmClient) RepoUpdate(names ...string) error {
	o := &repoUpdateOptions{
		update:    updateCharts,
		repoFile:  settings.RepositoryConfig,
		repoCache: settings.RepositoryCache,
		names:     names,
	}

	return o.run()
}
