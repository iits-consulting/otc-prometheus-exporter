package provider

import (
	"testing"

	dwsClusters "github.com/opentelekomcloud/gophertelekomcloud/openstack/dws/v1/cluster"
)

func TestConvertDWSClustersToMetrics(t *testing.T) {
	clusters := []dwsClusters.ClusterInfo{
		{Id: "dws-001", Name: "dw-prod", Status: "AVAILABLE"},
		{Id: "dws-002", Name: "dw-creating", Status: "CREATING"},
		{Id: "dws-003", Name: "dw-broken", Status: "UNAVAILABLE"},
	}

	families := convertDWSClustersToMetrics(clusters)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if families[0].GetName() != "dws_cluster_status" {
		t.Errorf("expected family name %q, got %q", "dws_cluster_status", families[0].GetName())
	}
	if len(families[0].Metric) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(families[0].Metric))
	}

	// AVAILABLE -> 0.0
	if families[0].Metric[0].Gauge.GetValue() != 0.0 {
		t.Errorf("expected AVAILABLE -> 0.0, got %f", families[0].Metric[0].Gauge.GetValue())
	}
	// CREATING -> 1.0
	if families[0].Metric[1].Gauge.GetValue() != 1.0 {
		t.Errorf("expected CREATING -> 1.0, got %f", families[0].Metric[1].Gauge.GetValue())
	}
	// UNAVAILABLE -> 1.0
	if families[0].Metric[2].Gauge.GetValue() != 1.0 {
		t.Errorf("expected UNAVAILABLE -> 1.0, got %f", families[0].Metric[2].Gauge.GetValue())
	}
}

func TestBuildDWSNameMap(t *testing.T) {
	clusters := []dwsClusters.ClusterInfo{
		{Id: "dws-001", Name: "dw-prod"},
		{Id: "dws-002", Name: "dw-staging"},
	}

	nameMap := buildDWSNameMap(clusters)

	if len(nameMap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(nameMap))
	}
	if nameMap["dws-001"] != "dw-prod" {
		t.Errorf("expected dws-001 -> %q, got %q", "dw-prod", nameMap["dws-001"])
	}
}
