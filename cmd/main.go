package main

import (
	"fmt"

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
		client, err := internal.NewOtcClientFromConfig(internal.Config)
		if err != nil {
			panic(err)
		}

		resourceIdToName, err := FetchResourceIdToNameMapping(client, internal.Config.Namespaces)
		if err != nil {
			panic(err)
		}

		metrics, err := client.GetMetrics()
		if err != nil {
			panic(err)
		}

		filteredMetrics := internal.FilterByNamespaces(metrics, internal.Config.Namespaces)

		internal.PrometheusMetrics = internal.SetupPrometheusMetricsFromOtcMetrics(filteredMetrics)

		for {
			fmt.Println(time.Now())

			for _, metric := range filteredMetrics {
				cloudeyeResponse, err := client.GetMetricData(metric)
				if err != nil {
					panic(err)
				}
				time.Sleep(time.Second)

				for _, d := range cloudeyeResponse.Datapoints {
					internal.PrometheusMetrics[internal.StandardPrometheusMetricName(metric)].With(
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

func FetchResourceIdToNameMapping(client *internal.OtcWrapper, namespaces []string) (map[string]string, error) {
	resourceIdToName := make(map[string]string)

	if slices.Contains(namespaces, internal.EcsNamespace) {
		result, err := client.GetEcsIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.RdsNamespace) {
		result, err := client.GetRdsIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.DmsNamespace) {
		result, err := client.GetDmsIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.NatNamespace) {
		result, err := client.GetNatIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.ElbNamespace) {
		result, err := client.GetElbIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}
	fmt.Printf("Collected %d resources\n", len(resourceIdToName))
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
