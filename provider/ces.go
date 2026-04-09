package provider

import (
	"context"
	"fmt"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	otcMetricData "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	otcMetrics "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metrics"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// CollectCESMetrics fetches all CES metrics for the given namespace using the
// project-scoped CES client. Most providers should use this.
func CollectCESMetrics(ctx context.Context, client *otcclient.Client, namespace string) ([]*dto.MetricFamily, error) {
	cesClient, err := client.CESClient()
	if err != nil {
		return nil, fmt.Errorf("ces client: %w", err)
	}
	return collectCESMetricsWithClient(ctx, cesClient, namespace)
}

// CollectCESMetricsRegionScoped fetches CES metrics using the region-level
// project scope. Use this for global services like OBS whose metrics are only
// visible under the region project (e.g. eu-de).
func CollectCESMetricsRegionScoped(ctx context.Context, client *otcclient.Client, namespace string) ([]*dto.MetricFamily, error) {
	cesClient, err := client.RegionCESClient()
	if err != nil {
		return nil, fmt.Errorf("region ces client: %w", err)
	}
	return collectCESMetricsWithClient(ctx, cesClient, namespace)
}

func collectCESMetricsWithClient(ctx context.Context, cesClient *golangsdk.ServiceClient, namespace string) ([]*dto.MetricFamily, error) {
	metrics, err := listMetricsByNamespace(cesClient, namespace)
	if err != nil {
		return nil, fmt.Errorf("list metrics: %w", err)
	}
	if len(metrics) == 0 {
		return nil, nil
	}

	data, err := fetchMetricDataBatched(ctx, cesClient, metrics)
	if err != nil {
		return nil, fmt.Errorf("fetch metric data: %w", err)
	}

	return ConvertBatchDataToFamilies(data), nil
}

// listMetricsByNamespace returns all CES metric definitions for the given namespace.
func listMetricsByNamespace(cesClient *golangsdk.ServiceClient, namespace string) ([]otcMetrics.MetricInfoList, error) {
	limit := 1000
	pages, err := otcMetrics.ListMetrics(cesClient, otcMetrics.ListMetricsRequest{
		Namespace: namespace,
		Limit:     &limit,
	}).AllPages()
	if err != nil {
		return nil, err
	}

	resp, err := otcMetrics.ExtractMetrics(pages)
	if err != nil {
		return nil, err
	}
	return resp.Metrics, nil
}

// fetchMetricDataBatched retrieves actual metric values for a list of CES
// metrics. It splits the request into batches of 500 via the batch API.
// The OTC SDK documents a limit of 10 per batch, but the API accepts 500
// in practice (matching the Huawei Cloud documentation).
func fetchMetricDataBatched(ctx context.Context, cesClient *golangsdk.ServiceClient, metrics []otcMetrics.MetricInfoList) ([]otcMetricData.BatchMetricData, error) {
	window, err := internal.NewSliceWindow(metrics, Config.CESBatchSize)
	if err != nil {
		return nil, fmt.Errorf("slice window: %w", err)
	}

	var result []otcMetricData.BatchMetricData
	for window.HasNext() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		batch := window.Window()
		data, err := fetchBatch(cesClient, batch)
		if err != nil {
			return nil, err
		}
		result = append(result, data...)
		window.NextWindow()
	}
	return result, nil
}

// fetchBatch sends a single BatchListMetricData request for a slice of metrics.
func fetchBatch(cesClient *golangsdk.ServiceClient, metrics []otcMetrics.MetricInfoList) ([]otcMetricData.BatchMetricData, error) {
	reqMetrics := make([]otcMetricData.Metric, len(metrics))
	for i, m := range metrics {
		dims := make([]otcMetricData.MetricsDimension, len(m.Dimensions))
		for j, d := range m.Dimensions {
			dims[j] = otcMetricData.MetricsDimension{
				Name:  d.Name,
				Value: d.Value,
			}
		}
		reqMetrics[i] = otcMetricData.Metric{
			Namespace:  m.Namespace,
			MetricName: m.MetricName,
			Dimensions: dims,
		}
	}

	now := time.Now()
	from := now.Add(-Config.CESLookback)

	return otcMetricData.BatchListMetricData(cesClient, otcMetricData.BatchListMetricDataOpts{
		Metrics: reqMetrics,
		From:    from.UnixMilli(),
		To:      now.UnixMilli(),
		Filter:  "average",
		Period:  "1",
	})
}

// ConvertBatchDataToFamilies converts CES BatchMetricData responses into
// Prometheus MetricFamily objects. Metrics are grouped by their Prometheus name
// (derived from namespace + metric_name). The latest datapoint value is used,
// but no timestamp is attached — Prometheus records scrape time instead. This
// avoids staleness gaps caused by CES only updating every ~5 minutes.
func ConvertBatchDataToFamilies(data []otcMetricData.BatchMetricData) []*dto.MetricFamily {
	familyMap := make(map[string]*dto.MetricFamily)

	for _, entry := range data {
		promName := PrometheusMetricName(entry.Namespace, entry.MetricName)

		fam, exists := familyMap[promName]
		if !exists {
			fam = NewGaugeMetricFamily(promName, nil)
			familyMap[promName] = fam
		}

		if len(entry.Datapoints) == 0 {
			continue
		}

		// Build labels from all dimensions to ensure uniqueness.
		labels := map[string]string{
			"unit":          entry.Unit,
			"resource_name": "",
		}
		for _, dim := range entry.Dimensions {
			labels[dim.Name] = dim.Value
		}
		if len(entry.Dimensions) > 0 {
			labels["resource_id"] = entry.Dimensions[0].Value
		}

		// Use latest datapoint (highest timestamp). No explicit timestamp --
		// Prometheus uses the scrape time. CES updates every ~5 minutes, so
		// attaching the real timestamp causes Prometheus staleness (>5min old
		// samples are dropped from instant queries, creating gaps).
		latest := entry.Datapoints[0]
		for _, dp := range entry.Datapoints[1:] {
			if dp.Timestamp > latest.Timestamp {
				latest = dp
			}
		}
		m := NewGaugeMetric(latest.Average, labels)
		fam.Metric = append(fam.Metric, m)
	}

	families := make([]*dto.MetricFamily, 0, len(familyMap))
	for _, fam := range familyMap {
		if len(fam.GetMetric()) > 0 {
			families = append(families, fam)
		}
	}
	return families
}
