package provider

import (
	"context"

	sfsShares "github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs/v2/shares"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// SFSProvider collects CES metrics for the OTC Scalable File Service,
// enriches them with share names, and reports share status.
type SFSProvider struct{}

func (p *SFSProvider) Namespace() string { return "SYS.SFS" }

func (p *SFSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.SFS", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		sfsClient, err := client.SFSV2()
		if err != nil {
			return nil, err
		}
		shares, err := sfsShares.List(sfsClient, sfsShares.ListOpts{})
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildSFSNameMap(shares),
			ExtraFamilies: convertSFSSharesToMetrics(shares),
		}, nil
	})
}

// buildSFSNameMap creates a mapping from SFS share ID to share name.
func buildSFSNameMap(shares []sfsShares.Share) map[string]string {
	m := make(map[string]string, len(shares))
	for _, s := range shares {
		m[s.ID] = s.Name
	}
	return m
}

// convertSFSSharesToMetrics creates a MetricFamily "sfs_share_status" with
// a gauge metric per share. The value is 0.0 for "available" shares, 1.0 otherwise
// (OTC convention: 0=normal, 1=abnormal).
func convertSFSSharesToMetrics(shares []sfsShares.Share) []*dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(shares))
	for _, s := range shares {
		value := 1.0
		if s.Status == "available" {
			value = 0.0
		}
		metrics = append(metrics, NewGaugeMetric(value, map[string]string{
			"resource_id":   s.ID,
			"resource_name": s.Name,
			"status":        s.Status,
		}))
	}
	return []*dto.MetricFamily{NewGaugeMetricFamily("sfs_share_status", metrics)}
}
