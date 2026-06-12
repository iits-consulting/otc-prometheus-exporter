package provider

import (
	"context"
	"fmt"
	"strings"
	"sync"

	otcMetricData "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	otcMetrics "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metrics"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

func isOBSStorageMetric(name string) bool {
	return strings.HasPrefix(name, "object_num_") || strings.HasPrefix(name, "capacity_")
}

// OBSProvider collects CES metrics for the OTC Object Storage Service.
// OBS is a global service whose CES metrics are only visible under the
// region-level project (e.g. eu-de), not under a specific project.
type OBSProvider struct{}

func (p *OBSProvider) Namespace() string { return "SYS.OBS" }

func (p *OBSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	cesClient, err := client.RegionCESClient()
	if err != nil {
		return nil, fmt.Errorf("region ces client: %w", err)
	}

	metrics, err := listMetricsByNamespace(cesClient, "SYS.OBS")
	if err != nil {
		return nil, fmt.Errorf("list metrics: %w", err)
	}
	if len(metrics) == 0 {
		return nil, nil
	}

	var storageMetrics, requestMetrics []otcMetrics.MetricInfoList
	for _, m := range metrics {
		if isOBSStorageMetric(m.MetricName) {
			storageMetrics = append(storageMetrics, m)
		} else {
			requestMetrics = append(requestMetrics, m)
		}
	}

	var (
		requestData []otcMetricData.BatchMetricData
		storageData []otcMetricData.BatchMetricData
		requestErr  error
		storageErr  error
		wg          sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if len(requestMetrics) > 0 {
			requestData, requestErr = fetchMetricDataBatched(ctx, cesClient, requestMetrics, cesPeriodRaw, Config.CESLookback)
		}
	}()
	go func() {
		defer wg.Done()
		if len(storageMetrics) > 0 {
			// Storage metrics have a 30min monitoring period.
			// They are fetched with Period=300, to ensure that data points are returned.
			// Because for Period=1, some metrics (for example object_num_all) are only returned for 6min consecutively out of 30min (even with large lookback).
			storageData, storageErr = fetchMetricDataBatched(ctx, cesClient, storageMetrics, cesPeriod5Min, Config.CESObsStorageMetricLookback)
		}
	}()
	wg.Wait()

	if requestErr != nil {
		return nil, fmt.Errorf("fetch request obs metrics: %w", requestErr)
	}
	if storageErr != nil {
		return nil, fmt.Errorf("fetch storage obs metrics: %w", storageErr)
	}

	families := ConvertBatchDataToFamilies(append(requestData, storageData...))

	// OBS CES metrics include a bucket_name dimension -- use it as resource_name.
	FillResourceNameFromLabel(families, "bucket_name")

	return families, nil
}
