package main

import (
	"fmt"
	"github.com/iits-consulting/otc-prometheus-exporter/otc_client"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"net/http"
	"time"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func collectMetricsInBackground() {
	go func() {
		client := otc_client.NewOtcClient(internal.Config.OtcProjectId, internal.Config.OtcProjectToken)

		resourceIdToName, err := FetchResourceIdToNameMapping(client, internal.Config.Namespaces)
		if err != nil {
			panic(err)
		}

		metrics, _ := client.GetMetricTypes()
		filteredMetrics := metrics.FilterByNamespaces(internal.Config.Namespaces)

		internal.PrometheusMetrics = internal.SetupPrometheusMetricsFromOtcMetrics(filteredMetrics)

		for {
			fmt.Println(time.Now())
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
			fmt.Println(time.Now())
			time.Sleep(internal.Config.WaitDuration)
		}
	}()
}

func FetchResourceIdToNameMapping(client otc_client.OtcClient, namespaces []string) (map[string]string, error) {
	resourceIdToName := make(map[string]string)

	if slices.Contains(namespaces, internal.EcsNamespace) {
		result, err := client.GetEcsData()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, internal.GetEcsIdToNameMapping(*result))
	}

	if slices.Contains(namespaces, internal.RdsNamespace) {
		// TODO: complete this and other remaining namespaces
	}

	return resourceIdToName, nil
}

func main() {

	collectMetricsInBackground()

	http.Handle("/metrics", promhttp.Handler())
	address := fmt.Sprintf(":%d", internal.Config.Port)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
