package provider

import (
	"testing"

	sfsShares "github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs/v2/shares"
)

func TestConvertSFSSharesToMetrics(t *testing.T) {
	shares := []sfsShares.Share{
		{ID: "sfs-001", Name: "share-prod", Status: "available"},
		{ID: "sfs-002", Name: "share-creating", Status: "creating"},
		{ID: "sfs-003", Name: "share-error", Status: "error"},
	}

	families := convertSFSSharesToMetrics(shares)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if families[0].GetName() != "sfs_share_status" {
		t.Errorf("expected family name %q, got %q", "sfs_share_status", families[0].GetName())
	}
	if len(families[0].Metric) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(families[0].Metric))
	}

	// available -> 0.0
	if families[0].Metric[0].Gauge.GetValue() != 0.0 {
		t.Errorf("expected available -> 0.0, got %f", families[0].Metric[0].Gauge.GetValue())
	}
	// creating -> 1.0
	if families[0].Metric[1].Gauge.GetValue() != 1.0 {
		t.Errorf("expected creating -> 1.0, got %f", families[0].Metric[1].Gauge.GetValue())
	}
	// error -> 1.0
	if families[0].Metric[2].Gauge.GetValue() != 1.0 {
		t.Errorf("expected error -> 1.0, got %f", families[0].Metric[2].Gauge.GetValue())
	}
}

func TestBuildSFSNameMap(t *testing.T) {
	shares := []sfsShares.Share{
		{ID: "sfs-001", Name: "share-prod"},
		{ID: "sfs-002", Name: "share-staging"},
	}

	nameMap := buildSFSNameMap(shares)

	if len(nameMap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(nameMap))
	}
	if nameMap["sfs-001"] != "share-prod" {
		t.Errorf("expected sfs-001 -> %q, got %q", "share-prod", nameMap["sfs-001"])
	}
}
