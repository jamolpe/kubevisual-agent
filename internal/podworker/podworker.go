package podworker

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type (
	PodWorker struct {
		clientSet        *kubernetes.Clientset
		metricsclientSet *metrics.Clientset
	}
	Pod struct {
		Name      string `json:"name"`
		Uid       string `json:"uid"`
		Node      string `json:"node"`
		Namespace string `json:"namespace"`
		Status    Status `json:"status"`
	}
	Status struct {
		Phase string `json:"phase"`
		Usage *Usage `json:"usage"`
	}

	Usage struct {
		Cpu    float64 `json:"cpu"`
		Memory int64   `json:"memory"`
	}
)

func New() *PodWorker {
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
	return &PodWorker{clientSet: clientset, metricsclientSet: metricsclientset}
}

func (pdc *PodWorker) mapPods(pods []corev1.Pod) []Pod {
	var finalPods []Pod
	fmt.Print()

	for _, p := range pods {
		podmetrics, err := pdc.metricsclientSet.MetricsV1beta1().PodMetricses(p.Namespace).Get(context.TODO(), p.GetName(), metav1.GetOptions{})
		pod := Pod{Name: p.Name, Uid: string(p.UID), Node: p.Spec.NodeName, Namespace: p.Namespace, Status: Status{Phase: string(p.Status.Phase)}}
		if err == nil {
			var totalMemory int64 = 0
			var totalCpu float64 = 0
			for _, cont := range podmetrics.Containers {
				memory, _ := cont.Usage.Memory().AsInt64()
				cpu := cont.Usage.Cpu().AsApproximateFloat64()
				totalMemory += memory
				totalCpu += cpu
			}
			pod.Status.Usage = &Usage{Cpu: totalCpu, Memory: totalMemory}
		}
		finalPods = append(finalPods, pod)
	}
	return finalPods
}

func (pdc *PodWorker) GetAllPodsInformation() ([]Pod, error) {

	pods, err := pdc.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pdc.mapPods(pods.Items), nil
}
