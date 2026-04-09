package otcdoc

import (
	"testing"
)

const rstWithNamespaceAndMetrics = `
Some preamble text.

Namespace
---------

SYS.NAT

Metrics
-------

+--------------------+------------------+----------+-------------------------+
| Metric ID          | Metric Name      | Unit     | Description             |
+====================+==================+==========+=========================+
| nat001_bytes_in    | Inbound Bytes    | Bytes/s  | Total inbound traffic.  |
+--------------------+------------------+----------+-------------------------+
| nat002_bytes_out   | Outbound Bytes   | Bytes/s  | Total outbound traffic. |
+--------------------+------------------+----------+-------------------------+
`

const rstWithDuplicateMetrics = `
Namespace
---------

SYS.NAT

+--------------------+------------------+----------+
| Metric ID          | Metric Name      | Unit     |
+====================+==================+==========+
| nat001_bytes_in    | Inbound Bytes    | Bytes/s  |
+--------------------+------------------+----------+
| nat001_bytes_in    | Inbound Bytes    | Bytes/s  |
+--------------------+------------------+----------+
`

const rstWithUnitInDescription = `
Namespace
---------

SYS.NAT

+--------------------+------------------+------------------------------------+
| Metric ID          | Metric Name      | Description                        |
+====================+==================+====================================+
| nat003_conn        | Connections      | Connection count. Unit: count      |
+--------------------+------------------+------------------------------------+
`

const rstNoNamespace = `
Some preamble without a namespace section.

+--------------------+------------------+
| Metric ID          | Metric Name      |
+====================+==================+
| nat001_bytes_in    | Inbound Bytes    |
+--------------------+------------------+
`

const rstNonMetricTable = `
Namespace
---------

SYS.NAT

+----------+-------------+
| Column A | Column B    |
+==========+=============+
| value1   | value2      |
+----------+-------------+
`

func TestParseNamespace(t *testing.T) {
	page, err := ParseDocumentationPageFromRstBytes([]byte(rstWithNamespaceAndMetrics))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Namespace != "SYS.NAT" {
		t.Errorf("expected namespace SYS.NAT, got %q", page.Namespace)
	}
}

func TestParseMetrics(t *testing.T) {
	page, err := ParseDocumentationPageFromRstBytes([]byte(rstWithNamespaceAndMetrics))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(page.Metrics))
	}
	if page.Metrics[0].MetricId != "nat001_bytes_in" {
		t.Errorf("expected metric ID nat001_bytes_in, got %q", page.Metrics[0].MetricId)
	}
	if page.Metrics[0].MetricName != "Inbound Bytes" {
		t.Errorf("expected metric name 'Inbound Bytes', got %q", page.Metrics[0].MetricName)
	}
	if page.Metrics[0].Unit != "Bytes/s" {
		t.Errorf("expected unit Bytes/s, got %q", page.Metrics[0].Unit)
	}
}

func TestDeduplicatesMetrics(t *testing.T) {
	page, err := ParseDocumentationPageFromRstBytes([]byte(rstWithDuplicateMetrics))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Metrics) != 1 {
		t.Errorf("expected 1 metric after dedup, got %d", len(page.Metrics))
	}
}

func TestExtractsUnitFromDescription(t *testing.T) {
	page, err := ParseDocumentationPageFromRstBytes([]byte(rstWithUnitInDescription))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(page.Metrics))
	}
	if page.Metrics[0].Unit != "count" {
		t.Errorf("expected unit 'count' extracted from description, got %q", page.Metrics[0].Unit)
	}
}

func TestEmptyNamespaceWhenSectionMissing(t *testing.T) {
	page, err := ParseDocumentationPageFromRstBytes([]byte(rstNoNamespace))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Namespace != "" {
		t.Errorf("expected empty namespace, got %q", page.Namespace)
	}
	if len(page.Metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(page.Metrics))
	}
}

func TestSkipsNonMetricTables(t *testing.T) {
	page, err := ParseDocumentationPageFromRstBytes([]byte(rstNonMetricTable))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Metrics) != 0 {
		t.Errorf("expected 0 metrics from non-metric table, got %d", len(page.Metrics))
	}
}

func TestEmptyInput(t *testing.T) {
	page, err := ParseDocumentationPageFromRstBytes([]byte(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Namespace != "" || len(page.Metrics) != 0 {
		t.Errorf("expected empty page, got %+v", page)
	}
}

func TestExtractUnitFromDesc(t *testing.T) {
	cases := []struct {
		desc string
		want string
	}{
		{"Connection count. Unit: count", "count"},
		{"Unit: Bytes/s more text", "Bytes/s"},
		{"No unit here", ""},
		{"Unit:", ""},
	}
	for _, c := range cases {
		got := extractUnitFromDesc(c.desc)
		if got != c.want {
			t.Errorf("extractUnitFromDesc(%q) = %q, want %q", c.desc, got, c.want)
		}
	}
}
