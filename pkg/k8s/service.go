package k8s

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"github.com/pkg/errors"

	"github.com/ghodss/yaml"
)

var (
	client *k8s.Client
)

func Init(kubeconfig string) {

	//client, err := k8s.NewInClusterClient()
	var err error
	client, err = loadClient(kubeconfig)
	if err != nil {
		client, err = k8s.NewInClusterClient()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func loadClient(kubeconfigPath string) (*k8s.Client, error) {
	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

func Service(name string) (service corev1.Service, err error) {
	err = client.Get(context.Background(), "", name, &service)
	if err != nil {
		err = fmt.Errorf("get service err %v", err)
		return
	}
	return
}

func ServiceListAll() (services []*corev1.Service, err error) {
	return ServiceList("")
}

func ServiceList(ns string) (services []*corev1.Service, err error) {
	var slist corev1.ServiceList
	err = client.List(context.Background(), ns, &slist)
	if err != nil {
		err = fmt.Errorf("get secret err %v", err)
		return
	}
	services = slist.GetItems()
	return
}

func nodeList() error {
	var nodes corev1.NodeList
	if err := client.List(context.Background(), "", &nodes); err != nil {
		return errors.Wrap(err, "client list")
	}
	for _, node := range nodes.Items {
		log.Printf("name=%q schedulable=%t\n", *node.Metadata.Name, !*node.Spec.Unschedulable)
	}
	return nil
}
