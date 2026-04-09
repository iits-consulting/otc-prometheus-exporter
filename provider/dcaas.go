package provider

import (
	"context"

	dcaasVgw "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/virtual-gateway"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// DCaaSProvider collects CES metrics for the OTC Direct Connect service,
// enriches them with virtual gateway names, and reports gateway status.
type DCaaSProvider struct{}

func (p *DCaaSProvider) Namespace() string { return "SYS.DCAAS" }

func (p *DCaaSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.DCAAS", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		dcaasClient, err := client.DCaaSV2()
		if err != nil {
			return nil, err
		}
		gateways, err := dcaasVgw.List(dcaasClient, dcaasVgw.ListOpts{})
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildDCaaSNameMap(gateways),
			ExtraFamilies: convertDCaaSGatewaysToMetrics(gateways),
		}, nil
	})
}

// buildDCaaSNameMap creates a mapping from virtual gateway ID to gateway name.
func buildDCaaSNameMap(gateways []dcaasVgw.VirtualGateway) map[string]string {
	m := make(map[string]string, len(gateways))
	for _, g := range gateways {
		m[g.ID] = g.Name
	}
	return m
}

// convertDCaaSGatewaysToMetrics creates a MetricFamily "dcaas_virtual_gateway_status"
// with a gauge per gateway. The value is 0.0 for ACTIVE gateways, 1.0 otherwise
// (OTC convention: 0=normal, 1=abnormal).
func convertDCaaSGatewaysToMetrics(gateways []dcaasVgw.VirtualGateway) []*dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(gateways))
	for _, g := range gateways {
		value := 1.0
		if g.Status == "ACTIVE" {
			value = 0.0
		}
		metrics = append(metrics, NewGaugeMetric(value, map[string]string{
			"resource_id":   g.ID,
			"resource_name": g.Name,
			"status":        g.Status,
		}))
	}
	return []*dto.MetricFamily{NewGaugeMetricFamily("dcaas_virtual_gateway_status", metrics)}
}
