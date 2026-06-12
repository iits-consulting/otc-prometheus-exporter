package provider

import (
	"context"
	"errors"
	"testing"

	otcMetricData "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	otcMetrics "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metrics"
)

func TestConvertBatchDataToFamilies(t *testing.T) {
	data := []otcMetricData.BatchMetricData{
		{
			Namespace:  "SYS.ECS",
			MetricName: "cpu_util",
			Dimensions: []otcMetricData.MetricsDimension{
				{Name: "instance_id", Value: "id-001"},
			},
			Unit: "%",
			Datapoints: []otcMetricData.DatapointForBatchMetric{
				{Average: 50.0, Timestamp: 1000},
				{Average: 75.0, Timestamp: 2000},
			},
		},
		{
			Namespace:  "SYS.ECS",
			MetricName: "mem_util",
			Dimensions: []otcMetricData.MetricsDimension{
				{Name: "instance_id", Value: "id-002"},
			},
			Unit: "%",
			Datapoints: []otcMetricData.DatapointForBatchMetric{
				{Average: 60.0, Timestamp: 1500},
			},
		},
	}

	families := ConvertBatchDataToFamilies(data)

	if len(families) != 2 {
		t.Fatalf("expected 2 families, got %d", len(families))
	}

	// Build a lookup by family name for stable assertions.
	famByName := make(map[string]int)
	for i, f := range families {
		famByName[f.GetName()] = i
	}

	// ecs_cpu_util has 2 datapoints -- only the latest (ts=2000) should be emitted.
	cpuIdx, ok := famByName["ecs_cpu_util"]
	if !ok {
		t.Fatal("expected family ecs_cpu_util not found")
	}
	cpuFam := families[cpuIdx]
	if len(cpuFam.Metric) != 1 {
		t.Fatalf("expected 1 metric in ecs_cpu_util (latest datapoint), got %d", len(cpuFam.Metric))
	}
	if cpuFam.Metric[0].GetGauge().GetValue() != 75.0 {
		t.Errorf("expected value 75.0 (latest), got %f", cpuFam.Metric[0].GetGauge().GetValue())
	}
	// No explicit timestamp -- Prometheus uses scrape time to avoid staleness gaps.
	if cpuFam.Metric[0].TimestampMs != nil {
		t.Errorf("expected no timestamp, got %d", cpuFam.Metric[0].GetTimestampMs())
	}

	// Verify labels.
	labels := make(map[string]string)
	for _, lp := range cpuFam.Metric[0].Label {
		labels[lp.GetName()] = lp.GetValue()
	}
	if labels["unit"] != "%" {
		t.Errorf("expected unit '%%', got %q", labels["unit"])
	}
	if labels["resource_id"] != "id-001" {
		t.Errorf("expected resource_id 'id-001', got %q", labels["resource_id"])
	}
	if labels["instance_id"] != "id-001" {
		t.Errorf("expected instance_id dimension label 'id-001', got %q", labels["instance_id"])
	}

	// ecs_mem_util has 1 datapoint.
	memIdx, ok := famByName["ecs_mem_util"]
	if !ok {
		t.Fatal("expected family ecs_mem_util not found")
	}
	memFam := families[memIdx]
	if len(memFam.Metric) != 1 {
		t.Fatalf("expected 1 metric in ecs_mem_util, got %d", len(memFam.Metric))
	}
	if memFam.Metric[0].Gauge.GetValue() != 60.0 {
		t.Errorf("expected gauge value 60.0, got %f", memFam.Metric[0].Gauge.GetValue())
	}
	if memFam.Metric[0].TimestampMs != nil {
		t.Errorf("expected no timestamp, got %d", memFam.Metric[0].GetTimestampMs())
	}
}

func TestConvertBatchDataEmptyDatapoints(t *testing.T) {
	data := []otcMetricData.BatchMetricData{
		{
			Namespace:  "SYS.ECS",
			MetricName: "cpu_util",
			Dimensions: []otcMetricData.MetricsDimension{
				{Name: "instance_id", Value: "id-001"},
			},
			Unit:       "%",
			Datapoints: []otcMetricData.DatapointForBatchMetric{},
		},
	}

	families := ConvertBatchDataToFamilies(data)

	// Families with no metrics are filtered out.
	if len(families) != 0 {
		t.Fatalf("expected 0 families (empty datapoints filtered), got %d", len(families))
	}
}

func TestConvertBatchDataMultipleDimensions(t *testing.T) {
	data := []otcMetricData.BatchMetricData{
		{
			Namespace:  "SYS.ECS",
			MetricName: "cpu_util",
			Dimensions: []otcMetricData.MetricsDimension{
				{Name: "instance_id", Value: "id-001"},
			},
			Unit: "%",
			Datapoints: []otcMetricData.DatapointForBatchMetric{
				{Average: 50.0, Timestamp: 1000},
			},
		},
		{
			Namespace:  "SYS.ECS",
			MetricName: "cpu_util",
			Dimensions: []otcMetricData.MetricsDimension{
				{Name: "instance_id", Value: "id-002"},
			},
			Unit: "%",
			Datapoints: []otcMetricData.DatapointForBatchMetric{
				{Average: 80.0, Timestamp: 1000},
			},
		},
	}

	families := ConvertBatchDataToFamilies(data)

	if len(families) != 1 {
		t.Fatalf("expected 1 family (grouped by metric name), got %d", len(families))
	}
	if len(families[0].Metric) != 2 {
		t.Fatalf("expected 2 metrics in the family, got %d", len(families[0].Metric))
	}

	// Collect resource_ids from both metrics.
	resourceIDs := make(map[string]float64)
	for _, m := range families[0].Metric {
		for _, lp := range m.Label {
			if lp.GetName() == "resource_id" {
				resourceIDs[lp.GetValue()] = m.Gauge.GetValue()
			}
		}
	}

	if val, ok := resourceIDs["id-001"]; !ok || val != 50.0 {
		t.Errorf("expected resource id-001 with value 50.0, got %v (exists=%v)", val, ok)
	}
	if val, ok := resourceIDs["id-002"]; !ok || val != 80.0 {
		t.Errorf("expected resource id-002 with value 80.0, got %v (exists=%v)", val, ok)
	}
}

func TestFetchMetricDataBatchedRespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	// With a cancelled context, fetchMetricDataBatched should return
	// before attempting any batch calls.
	metrics := make([]otcMetrics.MetricInfoList, 10)
	for i := range metrics {
		metrics[i] = otcMetrics.MetricInfoList{
			Namespace:  "SYS.ECS",
			MetricName: "cpu_util",
		}
	}

	_, err := fetchMetricDataBatched(ctx, nil, metrics, cesPeriodRaw, Config.CESLookback)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
}

func TestConvertBatchDataMultiDimensionLabels(t *testing.T) {
	// Simulates ELB-style metrics where the same metric name has entries with
	// the same first dimension (instance) but different second dimensions (listener).
	// All dimensions must become labels to avoid duplicate series.
	data := []otcMetricData.BatchMetricData{
		{
			Namespace:  "SYS.ELB",
			MetricName: "m1_cps",
			Dimensions: []otcMetricData.MetricsDimension{
				{Name: "lbaas_instance_id", Value: "lb-001"},
				{Name: "lbaas_listener_id", Value: "listener-001"},
			},
			Unit:       "Count",
			Datapoints: []otcMetricData.DatapointForBatchMetric{{Average: 100.0, Timestamp: 1000}},
		},
		{
			Namespace:  "SYS.ELB",
			MetricName: "m1_cps",
			Dimensions: []otcMetricData.MetricsDimension{
				{Name: "lbaas_instance_id", Value: "lb-001"},
				{Name: "lbaas_listener_id", Value: "listener-002"},
			},
			Unit:       "Count",
			Datapoints: []otcMetricData.DatapointForBatchMetric{{Average: 200.0, Timestamp: 1000}},
		},
	}

	families := ConvertBatchDataToFamilies(data)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if len(families[0].Metric) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(families[0].Metric))
	}

	// Verify both metrics have the listener dimension label so they are distinguishable.
	for _, m := range families[0].Metric {
		labels := make(map[string]string)
		for _, lp := range m.Label {
			labels[lp.GetName()] = lp.GetValue()
		}
		if labels["lbaas_instance_id"] != "lb-001" {
			t.Errorf("expected lbaas_instance_id=lb-001, got %q", labels["lbaas_instance_id"])
		}
		if labels["lbaas_listener_id"] == "" {
			t.Error("expected lbaas_listener_id label to be present")
		}
		if labels["resource_id"] != "lb-001" {
			t.Errorf("expected resource_id=lb-001, got %q", labels["resource_id"])
		}
	}
}
