package provider

import (
	"testing"

	elbLB "github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
	dto "github.com/prometheus/client_model/go"
)

func TestConvertELBToMetrics(t *testing.T) {
	lbs := []elbLB.LoadBalancer{
		{
			ID:                 "lb-001",
			Name:               "web-lb",
			ProvisioningStatus: "ACTIVE",
			OperatingStatus:    "ONLINE",
		},
		{
			ID:                 "lb-002",
			Name:               "api-lb",
			ProvisioningStatus: "ACTIVE",
			OperatingStatus:    "OFFLINE",
		},
	}

	families := convertELBToMetrics(lbs)

	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}
	if families[0].GetName() != "elb_loadbalancer_status" {
		t.Errorf("expected family name %q, got %q", "elb_loadbalancer_status", families[0].GetName())
	}
	if len(families[0].Metric) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(families[0].Metric))
	}
}

func TestEnrichELBNames(t *testing.T) {
	// Simulate CES metrics: one LB-level, one listener-level
	lbID := "lbaas_instance_id"
	lbVal := "lb-001"
	ridKey := "resource_id"
	rnKey := "resource_name"
	empty := ""
	listenerKey := "lbaas_listener_id"
	listenerVal := "listener-001"
	v100, v50 := float64(100), float64(50)

	families := []*dto.MetricFamily{
		NewGaugeMetricFamily("elb_m1_cps", []*dto.Metric{
			// LB-level (no listener ID)
			{Label: []*dto.LabelPair{
				{Name: &lbID, Value: &lbVal},
				{Name: &ridKey, Value: &lbVal},
				{Name: &rnKey, Value: &empty},
			}, Gauge: &dto.Gauge{Value: &v100}},
			// Listener-level
			{Label: []*dto.LabelPair{
				{Name: &lbID, Value: &lbVal},
				{Name: &listenerKey, Value: &listenerVal},
				{Name: &ridKey, Value: &lbVal},
				{Name: &rnKey, Value: &empty},
			}, Gauge: &dto.Gauge{Value: &v50}},
		}),
	}

	lbNames := map[string]string{"lb-001": "web-lb"}
	listenerNames := map[string]string{"listener-001": "https-443"}

	enrichELBNames(families, lbNames, listenerNames)

	metrics := families[0].GetMetric()

	// LB-level should get just the LB name
	lbMetricName := getLabelValue(metrics[0], "resource_name")
	if lbMetricName != "web-lb" {
		t.Errorf("LB-level: expected resource_name='web-lb', got %q", lbMetricName)
	}

	// Listener-level should get "lb/listener"
	listenerMetricName := getLabelValue(metrics[1], "resource_name")
	if listenerMetricName != "web-lb/https-443" {
		t.Errorf("Listener-level: expected resource_name='web-lb/https-443', got %q", listenerMetricName)
	}
}

func getLabelValue(m *dto.Metric, name string) string {
	for _, lp := range m.GetLabel() {
		if lp.GetName() == name {
			return lp.GetValue()
		}
	}
	return ""
}

func TestELBCacheIntegration(t *testing.T) {
	cache := NewNameCache()
	merged := map[string]string{
		"lb-1":       "my-loadbalancer",
		"listener-1": "my-listener",
	}
	cache.Put("SYS.ELB", merged)

	got := cache.Get("SYS.ELB")
	if len(got) != 2 {
		t.Fatalf("expected 2 cached entries, got %d", len(got))
	}
	if got["lb-1"] != "my-loadbalancer" {
		t.Errorf("expected lb-1 -> my-loadbalancer, got %q", got["lb-1"])
	}
}
