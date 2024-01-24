package helm

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/endverse/log/log"
	githuberrors "github.com/pkg/errors"
)

func TestMapToString(t *testing.T) {
	result := MapToString(map[string]string{
		"name":  "xiaoming",
		"class": "people",
	})

	fmt.Printf("result: %v \n", result)
}

func TestHelmClientList(t *testing.T) {
	helmClient := NewHelmClient("kf-partition")

	results, err := helmClient.List(map[string]string{
		"name": "ada3884db7ba2f",
	})

	if err != nil {
		t.Error(err)
		os.Exit(1)
	}

	for _, res := range results {
		log.Infof("helm release name: %v", res.Name)
		log.Infof("helm release info: %#v", *res.Info)
		log.Infof("helm release chart: %#v", *res.Chart.Metadata)
		log.Infof("helm release config: %v", res.Config)
		log.Infof("helm release hooks: %v", res.Hooks)
		log.Infof("helm release version: %v", res.Version)
		log.Infof("helm release namespace: %v", res.Namespace)
		log.Infof("helm release labels: %v", res.Labels)
	}
}

func TestHelmClientUninstall(t *testing.T) {
	helmClient := NewHelmClient("kf-partition")

	result, err := helmClient.Uninstall("ada3884db7ba2f")
	if err != nil {
		t.Error(err)
		log.Error(githuberrors.Cause(err))
		if strings.Contains(err.Error(), "release: not found") {
			log.Error("release not found, continue")
		}
		log.Error(reflect.TypeOf(err))
		os.Exit(1)
	}

	log.Infof("helm uninstall name: %v", result.Release.Name)
	log.Infof("helm uninstall info: %#v", *result.Release.Info)
}
