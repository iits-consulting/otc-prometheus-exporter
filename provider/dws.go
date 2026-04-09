package provider

import (
	"context"

	dwsClusters "github.com/opentelekomcloud/gophertelekomcloud/openstack/dws/v1/cluster"
	dto "github.com/prometheus/client_model/go"

	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
)

// DWSProvider collects CES metrics for the OTC Data Warehouse Service,
// enriches them with cluster names, and reports cluster status.
type DWSProvider struct{}

func (p *DWSProvider) Namespace() string { return "SYS.DWS" }

func (p *DWSProvider) Collect(ctx context.Context, client *otcclient.Client) ([]*dto.MetricFamily, error) {
	return CollectWithEnrichment(ctx, client, "SYS.DWS", func(ctx context.Context, client *otcclient.Client) (*EnrichResult, error) {
		dwsClient, err := client.DWSV1()
		if err != nil {
			return nil, err
		}
		resp, err := dwsClusters.ListClusters(dwsClient)
		if err != nil {
			return nil, err
		}
		return &EnrichResult{
			NameMap:       buildDWSNameMap(resp.Clusters),
			ExtraFamilies: convertDWSClustersToMetrics(resp.Clusters),
		}, nil
	})
}

// buildDWSNameMap creates a mapping from DWS cluster ID to cluster name.
func buildDWSNameMap(clusters []dwsClusters.ClusterInfo) map[string]string {
	m := make(map[string]string, len(clusters))
	for _, c := range clusters {
		m[c.Id] = c.Name
	}
	return m
}

// convertDWSClustersToMetrics creates a MetricFamily "dws_cluster_status" with
// a gauge metric per cluster. The value is 0.0 for AVAILABLE clusters, 1.0 otherwise
// (OTC convention: 0=normal, 1=abnormal).
func convertDWSClustersToMetrics(clusters []dwsClusters.ClusterInfo) []*dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, len(clusters))
	for _, c := range clusters {
		value := 1.0
		if c.Status == "AVAILABLE" {
			value = 0.0
		}
		metrics = append(metrics, NewGaugeMetric(value, map[string]string{
			"resource_id":   c.Id,
			"resource_name": c.Name,
			"status":        c.Status,
		}))
	}
	return []*dto.MetricFamily{NewGaugeMetricFamily("dws_cluster_status", metrics)}
}
