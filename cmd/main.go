package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
)

func main() {

	client := internal.NewOtcClient(internal.Config.OtcProjectId, internal.Config.OtcProjectToken)

	result, _ := client.GetEcsData()

	m := make(map[string]string)

	for _, s := range result.Servers {
		m[s.Id] = s.Name
	}

	metrics, _ := client.GetMetricTypes()
	filterdMetrics := metrics.FilterByNamespaces(internal.Config.Namespaces)

	prometheusGauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "TEST",
			Help: "The total number of processed events", // TODO
		},
		[]string{"unit", "resource_id", "resource_name"},
	)

	y, err := client.GetAllMetricData(filterdMetrics)
	if err != nil {
		panic(err)
	}
	fmt.Println(y)

	for i, datapoint := range y {
		fmt.Println("i", i)

		if len(datapoint.DataPoints) == 0 {
			continue
		}

		prometheusGauge.With(
			prometheus.Labels{
				"unit":          datapoint.DataPoints[0].Unit,
				"resource_id":   strconv.Itoa(i), // TODO: fix here the resource it to the dimension value
				"resource_name": strconv.Itoa(i), // TODO: fix this todos with the translated name
			},
		).Set(datapoint.DataPoints[0].Average)
	}

	fmt.Println(prometheusGauge)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
