package helm

import (
	"fmt"

	"github.com/endverse/log/log"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"

	kitstrings "github.com/endverse/go-kit/strings"
)

func (c *HelmClient) install(chart, version string, values map[string]interface{}) (string, string, error) {
	client := action.NewInstall(c.cfg)
	client.ReleaseName = "a" + kitstrings.RandomString(11).ToLower().String()
	client.InsecureSkipTLSverify = true

	var manifest string

	if settings.RepositoryCache == "" {
		settings.RepositoryCache = defaultCachePath
	}

	if settings.RepositoryConfig == "" {
		settings.RepositoryConfig = defaultRepositoryConfigPath
	}

	if version == "" {
		client.Version = ">0.0.0-0"
	} else {
		client.Version = version
	}

	chart = getRepoName() + "/" + chart

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return "", "", err
	}

	p := getter.All(settings)

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return "", "", err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return "", "", err
	}

	if chartRequested.Metadata.Deprecated {
		log.Warnf("This chart [%s] is deprecated", chartRequested.Name())
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return "", "", err
				}
			} else {
				return "", "", err
			}
		}
	}

	client.Namespace = settings.Namespace()

	// Debug Run
	client.DryRun = true
	rel, err := client.Run(chartRequested, values)
	if err != nil {
		return "", "", err
	}

	manifest = rel.Manifest

	// Run
	client.DryRun = false
	if c.linting {
		err = c.lint(cp, values)
		if err != nil {
			return manifest, "", errors.WithMessage(err, "helm lint")
		}
	}

	rel, err = client.Run(chartRequested, values)
	if err != nil {
		return manifest, "", errors.WithMessage(err, "helm install run")
	}

	// replace with install succeeded release.Manifest
	manifest = rel.Manifest

	log.Infof("release installed successfully: %s/%s-%s", rel.Name, rel.Name, rel.Chart.Metadata.Version)

	return manifest, client.ReleaseName, nil
}

func (c *HelmClient) lint(chartPath string, values map[string]interface{}) error {
	client := action.NewLint()

	result := client.Run([]string{chartPath}, values)

	for _, err := range result.Errors {
		log.Infof("Error %s", err)
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf("linting for chartpath %q failed", chartPath)
	}

	return nil
}
