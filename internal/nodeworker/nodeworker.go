package nodeworker

import (
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type (
	NodeWorker struct {
		kubeClient    *kubernetes.Clientset
		metricsClient *metrics.Clientset
	}
)

func New(kubeclient *kubernetes.Clientset, metricsClient *metrics.Clientset) *NodeWorker {
	return &NodeWorker{kubeClient: kubeclient, metricsClient: metricsClient}
}
