package otcdoc

import (
	"testing"
)

func TestOTCNamespaceUsesParsed(t *testing.T) {
	ds := DocumentationSource{Namespace: "ecs"}
	got := OTCNamespace(ds, "SYS.ECS")
	if got != "SYS.ECS" {
		t.Errorf("expected SYS.ECS, got %q", got)
	}
}

func TestOTCNamespaceCSSException(t *testing.T) {
	ds := DocumentationSource{Namespace: "css"}
	got := OTCNamespace(ds, "")
	if got != "SYS.ES" {
		t.Errorf("expected SYS.ES for css, got %q", got)
	}
}

func TestOTCNamespaceDDMException(t *testing.T) {
	ds := DocumentationSource{Namespace: "ddm"}
	got := OTCNamespace(ds, "")
	if got != "SYS.DDMS" {
		t.Errorf("expected SYS.DDMS for ddm, got %q", got)
	}
}

func TestOTCNamespaceFallback(t *testing.T) {
	ds := DocumentationSource{Namespace: "nat"}
	got := OTCNamespace(ds, "")
	if got != "SYS.NAT" {
		t.Errorf("expected SYS.NAT fallback, got %q", got)
	}
}

func TestPrometheusPrefixStripesSYS(t *testing.T) {
	got := PrometheusPrefix("SYS.ECS")
	if got != "ecs" {
		t.Errorf("expected ecs, got %q", got)
	}
}

func TestPrometheusPrefixStripsSERVICE(t *testing.T) {
	got := PrometheusPrefix("SERVICE.BMS")
	if got != "bms" {
		t.Errorf("expected bms, got %q", got)
	}
}

func TestPrometheusPrefixLowercases(t *testing.T) {
	got := PrometheusPrefix("SYS.NAT")
	if got != "nat" {
		t.Errorf("expected nat, got %q", got)
	}
}
