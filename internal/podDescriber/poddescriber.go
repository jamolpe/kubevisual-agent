package podDescriber

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type (
	PodDescriber struct {
		clientSet        *kubernetes.Clientset
		metricsclientSet *metrics.Clientset
	}
	Pod struct {
		Name   string `json:"name"`
		Uid    string `json:"uid"`
		Status Status `json:"status"`
	}
	Status struct {
		Phase string              `json:"phase"`
		Usage *v1beta1.PodMetrics `json:"usage"`
	}
)

func New() *PodDescriber {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	metricsclientset, err2 := metrics.NewForConfig(config)
	if err != nil || err2 != nil {
		panic(err.Error())
	}
	return &PodDescriber{clientSet: clientset, metricsclientSet: metricsclientset}
}

func (pdc *PodDescriber) mapPods(pods []corev1.Pod) []Pod {
	var finalPods []Pod
	fmt.Print()

	for _, p := range pods {
		podmetrics, err := pdc.metricsclientSet.MetricsV1beta1().PodMetricses(p.Namespace).Get(context.TODO(), p.GetName(), metav1.GetOptions{})
		pod := Pod{Name: p.Name, Uid: string(p.UID), Status: Status{Phase: string(p.Status.Phase)}}
		if err == nil {
			pod.Status.Usage = podmetrics
		}
		finalPods = append(finalPods, pod)
	}
	return finalPods
}

func (pdc *PodDescriber) GetAllPodsInformation() ([]Pod, error) {

	pods, err := pdc.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pdc.mapPods(pods.Items), nil
}
