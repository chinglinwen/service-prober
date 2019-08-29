package main

import (
	"flag"
	"os"
	"path/filepath"
	"wen/service-prober/pkg/k8s"
)

func init() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	if k := os.Getenv("KUBECONFIG"); k != "" {
		*kubeconfig = k
	}
	flag.Parse()

	k8s.Init(*kubeconfig)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
