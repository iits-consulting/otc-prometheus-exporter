package provider

import (
	"testing"

	bmsServers "github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/servers"
)

func TestConvertBMSInstancesToMetrics(t *testing.T) {
	servers := []bmsServers.Server{
		{ID: "bms-001", Name: "prod-bms-1", Status: "ACTIVE"},
		{ID: "bms-002", Name: "prod-bms-2", Status: "SHUTOFF"},
	}

	families := convertBMSInstancesToMetrics(servers)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if families[0].GetName() != "bms_instance_status" {
		t.Errorf("expected family name %q, got %q", "bms_instance_status", families[0].GetName())
	}
	if len(families[0].Metric) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(families[0].Metric))
	}

	// ACTIVE -> 0.0
	if families[0].Metric[0].Gauge.GetValue() != 0.0 {
		t.Errorf("expected ACTIVE status 0.0, got %f", families[0].Metric[0].Gauge.GetValue())
	}
	// SHUTOFF -> 1.0
	if families[0].Metric[1].Gauge.GetValue() != 1.0 {
		t.Errorf("expected SHUTOFF status 1.0, got %f", families[0].Metric[1].Gauge.GetValue())
	}
}

func TestBuildBMSNameMap(t *testing.T) {
	servers := []bmsServers.Server{
		{ID: "bms-001", Name: "prod-bms-1"},
		{ID: "bms-002", Name: "prod-bms-2"},
	}

	nameMap := buildBMSNameMap(servers)

	if len(nameMap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(nameMap))
	}
	if nameMap["bms-001"] != "prod-bms-1" {
		t.Errorf("expected bms-001 -> %q, got %q", "prod-bms-1", nameMap["bms-001"])
	}
	if nameMap["bms-002"] != "prod-bms-2" {
		t.Errorf("expected bms-002 -> %q, got %q", "prod-bms-2", nameMap["bms-002"])
	}
}
