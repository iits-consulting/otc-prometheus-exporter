package provider

import (
	"testing"

	dcaasVgw "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/virtual-gateway"
)

func TestConvertDCaaSGatewaysToMetrics(t *testing.T) {
	gateways := []dcaasVgw.VirtualGateway{
		{ID: "vgw-001", Name: "vgw-prod", Status: "ACTIVE"},
		{ID: "vgw-002", Name: "vgw-down", Status: "DOWN"},
		{ID: "vgw-003", Name: "vgw-error", Status: "ERROR"},
	}

	families := convertDCaaSGatewaysToMetrics(gateways)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if families[0].GetName() != "dcaas_virtual_gateway_status" {
		t.Errorf("expected family name %q, got %q", "dcaas_virtual_gateway_status", families[0].GetName())
	}
	if len(families[0].Metric) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(families[0].Metric))
	}

	// ACTIVE -> 0.0
	if families[0].Metric[0].Gauge.GetValue() != 0.0 {
		t.Errorf("expected ACTIVE -> 0.0, got %f", families[0].Metric[0].Gauge.GetValue())
	}
	// DOWN -> 1.0
	if families[0].Metric[1].Gauge.GetValue() != 1.0 {
		t.Errorf("expected DOWN -> 1.0, got %f", families[0].Metric[1].Gauge.GetValue())
	}
	// ERROR -> 1.0
	if families[0].Metric[2].Gauge.GetValue() != 1.0 {
		t.Errorf("expected ERROR -> 1.0, got %f", families[0].Metric[2].Gauge.GetValue())
	}
}

func TestBuildDCaaSNameMap(t *testing.T) {
	gateways := []dcaasVgw.VirtualGateway{
		{ID: "vgw-001", Name: "vgw-prod"},
		{ID: "vgw-002", Name: "vgw-staging"},
	}

	nameMap := buildDCaaSNameMap(gateways)

	if len(nameMap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(nameMap))
	}
	if nameMap["vgw-001"] != "vgw-prod" {
		t.Errorf("expected vgw-001 -> %q, got %q", "vgw-prod", nameMap["vgw-001"])
	}
}
