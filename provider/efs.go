package provider

import (
	"context"

	sfsTurboShares "github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs_turbo/v1/shares"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// EFSProvider collects CES metrics for the OTC Elastic File Service (SFS Turbo),
// enriches them with file system names, and reports file system status.
type EFSProvider struct{}

func (p *EFSProvider) Namespace() string { return "SYS.EFS" }

func (p *EFSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.EFS", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		efsClient, err := client.SFSTurboV1()
		if err != nil {
			return nil, err
		}
		pager := sfsTurboShares.List(efsClient, nil)
		page, err := pager.AllPages()
		if err != nil {
			return nil, err
		}
		turbos, err := sfsTurboShares.ExtractTurbos(page)
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildEFSNameMap(turbos),
			ExtraFamilies: convertEFSFileSysToMetrics(turbos),
		}, nil
	})
}

// buildEFSNameMap creates a mapping from EFS file system ID to file system name.
func buildEFSNameMap(turbos []sfsTurboShares.Turbo) map[string]string {
	m := make(map[string]string, len(turbos))
	for _, t := range turbos {
		m[t.ID] = t.Name
	}
	return m
}

// convertEFSFileSysToMetrics creates a MetricFamily "efs_filesystem_status" with
// a gauge metric per file system. The value is 0.0 for status "200" (available),
// 1.0 otherwise (OTC convention: 0=normal, 1=abnormal).
// EFS (SFS Turbo) uses HTTP-like numeric status codes as strings (e.g. "200" = healthy),
// unlike most OTC services which use word statuses ("ACTIVE", "AVAILABLE").
func convertEFSFileSysToMetrics(turbos []sfsTurboShares.Turbo) []*dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(turbos))
	for _, t := range turbos {
		value := 1.0
		if t.Status == "200" {
			value = 0.0
		}
		metrics = append(metrics, NewGaugeMetric(value, map[string]string{
			"resource_id":   t.ID,
			"resource_name": t.Name,
			"status":        t.Status,
		}))
	}
	return []*dto.MetricFamily{NewGaugeMetricFamily("efs_filesystem_status", metrics)}
}
