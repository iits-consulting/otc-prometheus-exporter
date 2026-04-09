package provider

import (
	"testing"

	sfsTurboShares "github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs_turbo/v1/shares"
)

func TestConvertEFSFileSysToMetrics(t *testing.T) {
	turbos := []sfsTurboShares.Turbo{
		{ID: "efs-001", Name: "efs-prod", Status: "200"},
		{ID: "efs-002", Name: "efs-creating", Status: "100"},
		{ID: "efs-003", Name: "efs-error", Status: "303"},
	}

	families := convertEFSFileSysToMetrics(turbos)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if families[0].GetName() != "efs_filesystem_status" {
		t.Errorf("expected family name %q, got %q", "efs_filesystem_status", families[0].GetName())
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

func TestBuildEFSNameMap(t *testing.T) {
	turbos := []sfsTurboShares.Turbo{
		{ID: "efs-001", Name: "efs-prod"},
		{ID: "efs-002", Name: "efs-staging"},
	}

	nameMap := buildEFSNameMap(turbos)

	if len(nameMap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(nameMap))
	}
	if nameMap["efs-001"] != "efs-prod" {
		t.Errorf("expected efs-001 -> %q, got %q", "efs-prod", nameMap["efs-001"])
	}
}
