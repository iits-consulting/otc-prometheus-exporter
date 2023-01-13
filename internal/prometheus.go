package internal

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var PrometheusMetrics map[string]*prometheus.GaugeVec
var PrometheusVectorLabels = []string{"unit", "resource_id", "resource_name"}

func SetupPrometheusMetricsFromOtcMetrics(otcMetrics MetricsResponse) map[string]*prometheus.GaugeVec {
	metrics := make(map[string]*prometheus.GaugeVec)

	for _, metric := range otcMetrics.Metrics {
		metricName := metric.StandardPrometheusMetricName()

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
