package helm

import (
	"os"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"
)

func getRepoUrl() string {
	url := os.Getenv("REPO_URL")
	if url != "" {
		return url
	}

	return repoUrl
}

func getRepoName() string {
	name := os.Getenv("REPO_NAME")
	if name != "" {
		return name
	}

	return repoName
}

func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func checkRequestedRepos(requestedRepos []string, validRepos []*repo.Entry) error {
	for _, requestedRepo := range requestedRepos {
		found := false
		for _, repo := range validRepos {
			if requestedRepo == repo.Name {
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("no repositories found matching '%s'.  Nothing will be updated", requestedRepo)
		}
	}
	return nil
}

func isRepoRequested(repoName string, requestedRepos []string) bool {
	for _, requestedRepo := range requestedRepos {
		if repoName == requestedRepo {
			return true
		}
	}
	return false
}

func isNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}

func MapToString(m map[string]string) string {
	var s string
	for k, v := range m {
		s += k + "=" + v + ","
	}

	return s[:len(s)-1]
}
