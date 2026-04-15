package provider

import (
	"context"

	cssClusters "github.com/opentelekomcloud/gophertelekomcloud/openstack/css/v1/clusters"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// CSSProvider collects CES metrics for the OTC Cloud Search Service,
// enriches them with cluster names, and reports cluster status.
type CSSProvider struct{}

func (p *CSSProvider) Namespace() string { return "SYS.ES" }

func (p *CSSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.ES", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		cssClient, err := client.CSSV1()
		if err != nil {
			return nil, err
		}
		clusters, err := cssClusters.List(cssClient)
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildCSSNameMap(clusters),
			ExtraFamilies: convertCSSClustersToMetrics(clusters),
		}, nil
	})
}

// buildCSSNameMap creates a mapping from CSS cluster ID to cluster name.
func buildCSSNameMap(clusters []cssClusters.Cluster) map[string]string {
	m := make(map[string]string, len(clusters))
	for _, c := range clusters {
		m[c.ID] = c.Name
	}
	return m
}

// convertCSSClustersToMetrics creates a MetricFamily "css_cluster_status" with
// a gauge metric per cluster. The value is 0.0 for status "200" (available),
// 1.0 otherwise (OTC convention: 0=normal, 1=abnormal).
// CSS uses HTTP-like numeric status codes as strings (e.g. "200" = healthy),
// unlike most OTC services which use word statuses ("ACTIVE", "AVAILABLE").
func convertCSSClustersToMetrics(clusters []cssClusters.Cluster) []*dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(clusters))
	for _, c := range clusters {
		value := 1.0
		if c.Status == "200" {
			value = 0.0
		}
		metrics = append(metrics, NewGaugeMetric(value, map[string]string{
			"resource_id":   c.ID,
			"resource_name": c.Name,
			"status":        c.Status,
		}))
	}
	return []*dto.MetricFamily{NewGaugeMetricFamily("css_cluster_status", metrics)}
}
