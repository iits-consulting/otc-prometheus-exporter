package provider

import (
	"context"
	"strconv"

	dmsInstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/instances/lifecycle"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// DMSProvider collects CES metrics for the OTC Distributed Message Service,
// enriches them with instance names, and reports storage and partition metrics.
type DMSProvider struct{}

func (p *DMSProvider) Namespace() string { return "SYS.DMS" }

func (p *DMSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.DMS", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		dmsClient, err := client.DMSV2()
		if err != nil {
			return nil, err
		}
		resp, err := dmsInstances.List(dmsClient, dmsInstances.ListOpts{})
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildDMSNameMap(resp.Instances),
			ExtraFamilies: convertDMSInstancesToMetrics(resp.Instances),
		}, nil
	})
}

// buildDMSNameMap creates a mapping from DMS instance ID to instance name.
func buildDMSNameMap(instances []dmsInstances.Instance) map[string]string {
	m := make(map[string]string, len(instances))
	for _, inst := range instances {
		m[inst.InstanceID] = inst.Name
	}
	return m
}

// convertDMSInstancesToMetrics creates MetricFamily objects for DMS-specific metrics:
// - dms_instance_status: 0.0 if RUNNING, 1.0 otherwise (OTC convention: 0=normal, 1=abnormal)
// - dms_instance_storage_used_gb: used storage space in GB
// - dms_instance_storage_total_gb: total storage space in GB
// - dms_instance_partitions: number of partitions
func convertDMSInstancesToMetrics(instances []dmsInstances.Instance) []*dto.MetricFamily {
	statusMetrics := make([]*dto.Metric, 0, len(instances))
	usedMetrics := make([]*dto.Metric, 0, len(instances))
	totalMetrics := make([]*dto.Metric, 0, len(instances))
	partitionMetrics := make([]*dto.Metric, 0, len(instances))

	for _, inst := range instances {
		labels := map[string]string{
			"resource_id":   inst.InstanceID,
			"resource_name": inst.Name,
		}

		statusValue := 1.0
		if inst.Status == "RUNNING" {
			statusValue = 0.0
		}
		statusMetrics = append(statusMetrics, NewGaugeMetric(statusValue, map[string]string{
			"resource_id":   inst.InstanceID,
			"resource_name": inst.Name,
			"status":        inst.Status,
		}))

		usedMetrics = append(usedMetrics, NewGaugeMetric(float64(inst.UsedStorageSpace), labels))
		totalMetrics = append(totalMetrics, NewGaugeMetric(float64(inst.TotalStorageSpace), labels))

		partitions := 0.0
		if inst.PartitionNum != "" {
			if v, err := strconv.ParseFloat(inst.PartitionNum, 64); err == nil {
				partitions = v
			}
		}
		partitionMetrics = append(partitionMetrics, NewGaugeMetric(partitions, labels))
	}

	return []*dto.MetricFamily{
		NewGaugeMetricFamily("dms_instance_status", statusMetrics),
		NewGaugeMetricFamily("dms_instance_storage_used_gb", usedMetrics),
		NewGaugeMetricFamily("dms_instance_storage_total_gb", totalMetrics),
		NewGaugeMetricFamily("dms_instance_partitions", partitionMetrics),
	}
}
