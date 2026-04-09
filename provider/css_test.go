package provider

import (
	"testing"

	cssClusters "github.com/opentelekomcloud/gophertelekomcloud/openstack/css/v1/clusters"
)

func TestConvertCSSClustersToMetrics(t *testing.T) {
	clusters := []cssClusters.Cluster{
		{ID: "css-001", Name: "es-prod", Status: "200"},
		{ID: "css-002", Name: "es-staging", Status: "100"},
		{ID: "css-003", Name: "es-broken", Status: "303"},
	}

	families := convertCSSClustersToMetrics(clusters)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if families[0].GetName() != "css_cluster_status" {
		t.Errorf("expected family name %q, got %q", "css_cluster_status", families[0].GetName())
	}
	if len(families[0].Metric) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(families[0].Metric))
	}

	// "200" -> 0.0
	if families[0].Metric[0].Gauge.GetValue() != 0.0 {
		t.Errorf("expected status 200 -> 0.0, got %f", families[0].Metric[0].Gauge.GetValue())
	}
	// "100" -> 1.0
	if families[0].Metric[1].Gauge.GetValue() != 1.0 {
		t.Errorf("expected status 100 -> 1.0, got %f", families[0].Metric[1].Gauge.GetValue())
	}
	// "303" -> 1.0
	if families[0].Metric[2].Gauge.GetValue() != 1.0 {
		t.Errorf("expected status 303 -> 1.0, got %f", families[0].Metric[2].Gauge.GetValue())
	}
}

func TestBuildCSSNameMap(t *testing.T) {
	clusters := []cssClusters.Cluster{
		{ID: "css-001", Name: "es-prod"},
		{ID: "css-002", Name: "es-staging"},
	}

	nameMap := buildCSSNameMap(clusters)

	if len(nameMap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(nameMap))
	}
	if nameMap["css-001"] != "es-prod" {
		t.Errorf("expected css-001 -> %q, got %q", "es-prod", nameMap["css-001"])
	}
}
