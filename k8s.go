package main

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

type podRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}
type cpuStats struct {
	Time                 time.Time `json:"time"`
	UsageNanoCores       int       `json:"usageNanoCores"`
	UsageCoreNanoSeconds int       `json:"usageCoreNanoSeconds"`
}

type memoryStats struct {
	Time            time.Time `json:"time"`
	AvailableBytes  int       `json:"availableBytes"`
	UsageBytes      int       `json:"usageBytes"`
	WorkingSetBytes int       `json:"workingSetBytes"`
	RSSBytes        int       `json:"rssBytes"`
	PageFaults      int       `json:"pageFaults"`
	MajorPageFaults int       `json:"majorPageFaults"`
}

type container struct {
	Name      string      `json:"name"`
	StartTime time.Time   `json:"startTime"`
	CPU       cpuStats    `json:"cpu"`
	Memory    memoryStats `json:"memory"`
	// skip rootfs and logs for now
}
type pod struct {
	PodRef     podRef      `json:"podRef"`
	StartTime  time.Time   `json:"startTime"`
	Containers []container `json:"containers"`
	CPU        cpuStats    `json:"cpu"`
	Memory     memoryStats `json:"memory"`
	// skip everything else for now
}

type nodeSummary struct {
	// interested only in pods at the moment
	Pods []pod `json:"pods"`
}

func getClientset() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
