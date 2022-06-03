package exporter

import (
	"context"
	"encoding/json"

	"github.com/mercury200Hg/metrics-server-prometheus-exporter/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
)

// ContainerMetricsData is a struct that defines structure of container metrics data
type ContainerMetricsData struct {
	Name  string `json:"name"`
	Usage struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	}
}

// PodMetricsData is a struct that defines structure of a Pod's metric data individually specifying data of it's containers
type PodMetricsData struct {
	Timestamp string `json:"timestamp"`
	Window    string `json:"window"`
	Metadata  struct {
		Name              string `json:"name"`
		Namespace         string `json:"namespace"`
		SelfLink          string `json:"selfLink"`
		CreationTimestamp string `json:"creationTimestamp"`
	}
	Containers []ContainerMetricsData
}

// PodMetrics is a struct that defines structure of the response returned by the metrics-server api
type PodMetrics struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	MetaData   struct {
		SelfLink string `json:"selfLink"`
	}
	Items []PodMetricsData
}

// PodMetricCPU represents the Gauge Vector for Pod wise CPU usage of all containers
var PodMetricCPU = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "kube_metrics_server_pod_cpu", Help: "Shows the current cpu usage in number of cores (cpu shares) of a given container for a given pod"}, []string{
	"pod_name",
	"pod_namespace",
	"container_name",
})

// PodMetricMemory represents the Gauge Vector for Pod wise Memory usage of all containers
var PodMetricMemory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "kube_metrics_server_pod_mem", Help: "Shows the current memory usage of a given container for a given pod"}, []string{
	"pod_name",
	"pod_namespace",
	"container_name",
})

// GetPodMetric returns Pod wise CPU and memory usage of the containers
func getPodMetric() PodMetrics {
	var data PodMetrics
	if utils.KubeConfig == nil {
		utils.InitKubeConfig()
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(utils.KubeConfig)
	if err != nil {
		log.Error().Msg(err.Error())
	} else {
		result, err := clientset.RESTClient().Get().RequestURI("/apis/metrics.k8s.io/v1beta1/pods").Do(context.TODO()).Raw()
		if err != nil {
			log.Err(err)
		} else {
			dataString := string(result)                    // JSON String of data
			err = json.Unmarshal([]byte(dataString), &data) // JSON Data
			if err != nil {
				log.Error().Err(err)
			}
		}
	}
	return data
}
func setGaugeMetric(name string, help string, label string, labelvalue string) prometheus.Gauge {
	var gaugeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        name,
		Help:        help,
		ConstLabels: prometheus.Labels{label: labelvalue},
	})
	return gaugeMetric
}

// RecordPodMetrics records the metrics and sets it to promhttp Handler
func RecordPodMetrics() {
	data := getPodMetric()
	PodMetricCPU.Reset()
	PodMetricMemory.Reset()
	for i := range data.Items {
		podName := data.Items[i].Metadata.Name
		podNamespace := data.Items[i].Metadata.Namespace
		for j := range data.Items[i].Containers {
			containerName := data.Items[i].Containers[j].Name
			cpuUsage, errCPU := utils.ParseCPU(data.Items[i].Containers[j].Usage.CPU)
			memUsage, errMem := utils.ParseMemory(data.Items[i].Containers[j].Usage.Memory)
			if errCPU == nil {
				PodMetricCPU.With(prometheus.Labels{"pod_name": podName, "pod_namespace": podNamespace, "container_name": containerName}).Set(cpuUsage)
			} else {
				log.Err(errCPU)
			}
			if errMem == nil {
				PodMetricMemory.With(prometheus.Labels{"pod_name": podName, "pod_namespace": podNamespace, "container_name": containerName}).Set(memUsage)
			} else {
				log.Err(errMem)
			}
		}
	}
}
