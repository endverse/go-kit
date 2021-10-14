package helm

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hex-techs/klog"
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
		klog.Infof("helm release name: %v", res.Name)
		klog.Infof("helm release info: %#v", *res.Info)
		klog.Infof("helm release chart: %#v", *res.Chart.Metadata)
		klog.Infof("helm release config: %v", res.Config)
		klog.Infof("helm release hooks: %v", res.Hooks)
		klog.Infof("helm release version: %v", res.Version)
		klog.Infof("helm release namespace: %v", res.Namespace)
		klog.Infof("helm release labels: %v", res.Labels)
	}
}

func TestHelmClientUninstall(t *testing.T) {
	klog.InitFlags(nil)
	flag.Parse()

	helmClient := NewHelmClient("kf-partition")

	result, err := helmClient.Uninstall("ada3884db7ba2f")
	if err != nil {
		t.Error(err)
		klog.Error(githuberrors.Cause(err))
		if strings.Contains(err.Error(), "release: not found") {
			klog.Error("release not found, continue")
		}
		klog.Error(reflect.TypeOf(err))
		os.Exit(1)
	}

	klog.Infof("helm uninstall name: %v", result.Release.Name)
	klog.Infof("helm uninstall info: %#v", *result.Release.Info)
}
