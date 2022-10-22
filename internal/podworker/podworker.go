package podworker

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8s.io/client-go/kubernetes"
)

type (
	PodWorker struct {
		kubeClient    *kubernetes.Clientset
		metricsClient *metrics.Clientset
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

func New(kubeclient *kubernetes.Clientset, metricsClient *metrics.Clientset) *PodWorker {
	return &PodWorker{kubeClient: kubeclient, metricsClient: metricsClient}
}

func (pw *PodWorker) mapPods(pods []corev1.Pod) []Pod {
	var finalPods []Pod
	fmt.Print()

	for _, p := range pods {
		podmetrics, err := pw.metricsClient.MetricsV1beta1().PodMetricses(p.Namespace).Get(context.TODO(), p.GetName(), metav1.GetOptions{})
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

func (pw *PodWorker) GetAllPodsInformation() ([]Pod, error) {
	pods, err := pw.kubeClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pw.mapPods(pods.Items), nil
}
