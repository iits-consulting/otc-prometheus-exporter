package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func collectMetricsInBackground() {
	go func() {
		
	}()
}

func main() {

	client := internal.NewOtcClient(internal.Config.OtcProjectId, internal.Config.OtcProjectToken)

	result, err := client.GetEcsData()
	if err != nil {
		panic(err)
	}

	resourceIdToName := make(map[string]string)
	for _, s := range result.Servers {
		resourceIdToName[s.Id] = s.Name
	}

	metrics, _ := client.GetMetricTypes()
	filteredMetrics := metrics.FilterByNamespaces(internal.Config.Namespaces)

	internal.PrometheusMetrics = internal.SetupPrometheusMetricsFromOtcMetrics(filteredMetrics)

	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Second)
	for _, metric := range filteredMetrics.Metrics {
		cloudeyeResponse, err := client.GetMetricData(
			metric.Namespace,
			metric.MetricName,
			metric.Dimensions[0].Name,
			metric.Dimensions[0].Value,
			startTime,
			endTime,
		)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)

		for _, d := range cloudeyeResponse.DataPoints {
			internal.PrometheusMetrics[metric.StandardPrometheusMetricName()].With(
				prometheus.Labels{
					"unit":          d.Unit,
					"resource_id":   metric.Dimensions[0].Value,
					"resource_name": resourceIdToName[metric.Dimensions[0].Value],
				}).Set(d.Average)
		}
	}

	http.Handle("/metrics", promhttp.Handler())
	address := fmt.Sprintf(":%d", internal.Config.Port)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
