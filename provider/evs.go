package provider

import (
	"context"
	"strings"

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
		blockClient, err := client.BlockStorageV2()
		if err != nil {
			return nil, err
		}
		volumes, err := evsVolumes.List(blockClient, evsVolumes.ListOpts{})
		if err != nil {
			return nil, err
		}
		nameMap := buildEVSNameMap(volumes)
		return &EnrichResult{
			NameMap:       nameMap,
			ExtraFamilies: convertEVSVolumesToMetrics(volumes),
			Enrich:        func(families []*dto.MetricFamily) { enrichEVSNames(families, nameMap) },
		}, nil
	})
}

// enrichEVSNames resolves resource_name for EVS metric families.
// Any metric still without a name falls back to its resource_id.
func enrichEVSNames(families []*dto.MetricFamily, nameMap map[string]string) {
	EnrichWithNames(families, nameMap)
	FillResourceNameFromLabel(families, "resource_id")
}

// buildEVSNameMap maps EVS resource ids to volume names.
// It keys both the plain volume ID (for volume-level metrics like evs_volume_status)
// and per-device metrics (e.g. "835e980e-...-vda")
func buildEVSNameMap(volumes []evsVolumes.Volume) map[string]string {
	m := make(map[string]string, len(volumes))
	for _, v := range volumes {
		m[v.ID] = v.Name
		for _, a := range v.Attachments {
			device := strings.TrimPrefix(a.Device, "/dev/")
			if a.ServerID != "" && device != "" {
				m[a.ServerID+"-"+device] = v.Name
			}
		}
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
