package provider

import (
	"testing"

	evsVolumes "github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/cloudvolumes"
)

func TestConvertEVSVolumesToMetrics(t *testing.T) {
	volumes := []evsVolumes.Volume{
		{ID: "evs-001", Name: "data-disk-1", Status: "available", Size: 100},
		{ID: "evs-002", Name: "data-disk-2", Status: "in-use", Size: 50},
		{ID: "evs-003", Name: "data-disk-3", Status: "error", Size: 200},
	}

	families := convertEVSVolumesToMetrics(volumes)

	if len(families) != 2 {
		t.Fatalf("expected 2 families, got %d", len(families))
	}

	expectedNames := []string{"evs_volume_status", "evs_volume_size_gb"}
	for i, name := range expectedNames {
		if families[i].GetName() != name {
			t.Errorf("expected family[%d] name %q, got %q", i, name, families[i].GetName())
		}
	}

	statusFam := families[0]
	if len(statusFam.Metric) != 3 {
		t.Fatalf("expected 3 status metrics, got %d", len(statusFam.Metric))
	}
	// available -> 0.0
	if statusFam.Metric[0].Gauge.GetValue() != 0.0 {
		t.Errorf("expected available status 0.0, got %f", statusFam.Metric[0].Gauge.GetValue())
	}
	// in-use -> 0.0
	if statusFam.Metric[1].Gauge.GetValue() != 0.0 {
		t.Errorf("expected in-use status 0.0, got %f", statusFam.Metric[1].Gauge.GetValue())
	}
	// error -> 1.0
	if statusFam.Metric[2].Gauge.GetValue() != 1.0 {
		t.Errorf("expected error status 1.0, got %f", statusFam.Metric[2].Gauge.GetValue())
	}

	sizeFam := families[1]
	if sizeFam.Metric[0].Gauge.GetValue() != 100.0 {
		t.Errorf("expected size 100.0, got %f", sizeFam.Metric[0].Gauge.GetValue())
	}
	if sizeFam.Metric[2].Gauge.GetValue() != 200.0 {
		t.Errorf("expected size 200.0, got %f", sizeFam.Metric[2].Gauge.GetValue())
	}
}

func TestBuildEVSNameMap(t *testing.T) {
	volumes := []evsVolumes.Volume{
		{ID: "evs-001", Name: "data-disk-1"},
		{ID: "evs-002", Name: "data-disk-2"},
	}

	nameMap := buildEVSNameMap(volumes)

	if len(nameMap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(nameMap))
	}
	if nameMap["evs-001"] != "data-disk-1" {
		t.Errorf("expected evs-001 -> %q, got %q", "data-disk-1", nameMap["evs-001"])
	}
}
