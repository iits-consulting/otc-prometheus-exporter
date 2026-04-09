package provider

import (
	"context"

	bmsServers "github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/servers"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// BMSProvider collects CES metrics for the OTC Bare Metal Server service,
// enriches them with server names, and reports instance status.
type BMSProvider struct{}

func (p *BMSProvider) Namespace() string { return "SYS.BMS" }

func (p *BMSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.BMS", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		computeClient, err := client.ComputeV2()
		if err != nil {
			return nil, err
		}
		servers, err := bmsServers.List(computeClient, bmsServers.ListOpts{})
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildBMSNameMap(servers),
			ExtraFamilies: convertBMSInstancesToMetrics(servers),
		}, nil
	})
}

// buildBMSNameMap creates a mapping from BMS server ID to server name.
func buildBMSNameMap(servers []bmsServers.Server) map[string]string {
	m := make(map[string]string, len(servers))
	for _, s := range servers {
		m[s.ID] = s.Name
	}
	return m
}

// convertBMSInstancesToMetrics creates a MetricFamily "bms_instance_status" with
// a gauge metric per server. The value is 0.0 for ACTIVE servers, 1.0 otherwise
// (OTC convention: 0=normal, 1=abnormal).
func convertBMSInstancesToMetrics(servers []bmsServers.Server) []*dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(servers))
	for _, s := range servers {
		value := 1.0
		if s.Status == "ACTIVE" {
			value = 0.0
		}
		metrics = append(metrics, NewGaugeMetric(value, map[string]string{
			"resource_id":   s.ID,
			"resource_name": s.Name,
			"status":        s.Status,
		}))
	}
	return []*dto.MetricFamily{NewGaugeMetricFamily("bms_instance_status", metrics)}
}
