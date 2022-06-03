package exporter

import (
	"context"
	"encoding/json"

	"github.com/mercury200Hg/metrics-server-prometheus-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
)

// NodeMetricsData is a struct that defines structure of a Node's metric data individually specifying data of it's containers
type NodeMetricsData struct {
	Metadata struct {
		CreationTimestamp string `json:"creationTimestamp"`
		Name              string `json:"name"`
		SelfLink          string `json:"selfLink"`
	}
	Timestamp string `json:"timestamp"`
	Window    string `json:"window"`
	Usage     struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	}
}

// NodeMetrics is a struct that defines structure of the response returned by the metrics-server api
type NodeMetrics struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	MetaData   struct {
		SelfLink string `json:"selfLink"`
	}
	Items []NodeMetricsData
}

// NodeMetricCPU represents the Gauge Vector for Node wise CPU usage
var NodeMetricCPU = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "kube_metrics_server_node_cpu", Help: "Shows the current cpu usage in number of cores (cpu shares) of a given node as visible to kubelet"}, []string{
	"node_name",
})

// NodeMetricMemory represents the Gauge Vector for Node wise Memory usage
var NodeMetricMemory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "kube_metrics_server_node_mem", Help: "Shows the current memory in bytes usage of given node as visible to kubelet"}, []string{
	"node_name",
})

// GetNodeMetric returns Pod wise CPU and memory usage of the containers
func getNodeMetric() NodeMetrics {
	var data NodeMetrics
	if utils.KubeConfig == nil {
		utils.InitKubeConfig()
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(utils.KubeConfig)
	if err != nil {
		log.Error().Msg(err.Error())
	} else {
		result, err := clientset.RESTClient().Get().RequestURI("/apis/metrics.k8s.io/v1beta1/nodes").Do(context.TODO()).Raw()
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

// RecordNodeMetrics records the metrics and sets it to promhttp Handler
func RecordNodeMetrics() {
	NodeMetricCPU.Reset()
	NodeMetricMemory.Reset()
	data := getNodeMetric()
	for i := range data.Items {
		nodeName := data.Items[i].Metadata.Name
		cpuUsage, errCPU := utils.ParseCPU(data.Items[i].Usage.CPU)
		memUsage, errMem := utils.ParseMemory(data.Items[i].Usage.Memory)
		if errCPU == nil {
			NodeMetricCPU.With(prometheus.Labels{"node_name": nodeName}).Set(cpuUsage)
		} else {
			log.Err(errCPU)
		}
		if errMem == nil {
			NodeMetricMemory.With(prometheus.Labels{"node_name": nodeName}).Set(memUsage)
		} else {
			log.Err(errMem)
		}
	}
}
