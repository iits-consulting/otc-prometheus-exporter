package provider

import (
	"context"

	evsVolumes "github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/cloudvolumes"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// EVSProvider collects CES metrics for the OTC Elastic Volume Service,
// enriches them with volume names, and reports volume status and size.
type EVSProvider struct{}

func (p *EVSProvider) Namespace() string { return "SYS.EVS" }

func (p *EVSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.EVS", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		blockClient, err := client.BlockStorageV3()
		if err != nil {
			return nil, err
		}
		volumes, err := evsVolumes.List(blockClient, evsVolumes.ListOpts{})
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildEVSNameMap(volumes),
			ExtraFamilies: convertEVSVolumesToMetrics(volumes),
		}, nil
	})
}

// buildEVSNameMap creates a mapping from EVS volume ID to volume name.
func buildEVSNameMap(volumes []evsVolumes.Volume) map[string]string {
	m := make(map[string]string, len(volumes))
	for _, v := range volumes {
		m[v.ID] = v.Name
	}
	return m
}

// convertEVSVolumesToMetrics creates MetricFamily objects for EVS-specific metrics:
// - evs_volume_status: 0.0 if available or in-use, 1.0 for error states
// - evs_volume_size_gb: volume size in GB
func convertEVSVolumesToMetrics(volumes []evsVolumes.Volume) []*dto.MetricFamily {
	statusMetrics := make([]*dto.Metric, 0, len(volumes))
	sizeMetrics := make([]*dto.Metric, 0, len(volumes))

	for _, v := range volumes {
		labels := map[string]string{
			"resource_id":   v.ID,
			"resource_name": v.Name,
		}

		statusValue := 0.0
		if v.Status == "error" || v.Status == "error_deleting" || v.Status == "error_restoring" || v.Status == "error_extending" {
			statusValue = 1.0
		}
		statusMetrics = append(statusMetrics, NewGaugeMetric(statusValue, map[string]string{
			"resource_id":   v.ID,
			"resource_name": v.Name,
			"status":        v.Status,
		}))

		sizeMetrics = append(sizeMetrics, NewGaugeMetric(float64(v.Size), labels))
	}

	return []*dto.MetricFamily{
		NewGaugeMetricFamily("evs_volume_status", statusMetrics),
		NewGaugeMetricFamily("evs_volume_size_gb", sizeMetrics),
	}
}
