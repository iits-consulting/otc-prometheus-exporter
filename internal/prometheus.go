package internal

import (
	"fmt"
	otcMetrics "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var PrometheusMetrics map[string]*prometheus.GaugeVec
var PrometheusVectorLabels = []string{"unit", "resource_id", "resource_name"}

func SetupPrometheusMetricsFromOtcMetrics(otcMetrics []otcMetrics.MetricInfoList) map[string]*prometheus.GaugeVec {
	metrics := make(map[string]*prometheus.GaugeVec)

	for _, metric := range otcMetrics {
		metricName := StandardPrometheusMetricName(metric)

		if _, ok := metrics[metricName]; !ok {
			fmt.Println("created prometheus metric", metricName)
			metrics[metricName] = promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: metricName,
					Help: "The total number of processed events", // TODO
				},
				PrometheusVectorLabels,
			)
		}
	}

	return metrics

}
