package main

import "github.com/iits-consulting/otc-prometheus-exporter/grafana"

// exporterDashboard returns the hand-crafted dashboard for the OTC Prometheus Exporter
// itself. It uses custom PromQL expressions (histogram_quantile, rate) that cannot be
// auto-generated from OTC documentation.
func exporterDashboard() grafana.DashboardConfig {
	return grafana.DashboardConfig{
		Title: "OTC Exporter",
		UID:   "otc-exporter",
		Sections: []grafana.PanelSection{
			{Title: "Scrape Duration", Panels: []grafana.PanelConfig{
				{Metric: "otc_scrape_duration_seconds_bucket", Title: "Scrape Duration p50", Unit: "s", Type: grafana.TimeSeries,
					Expr:   `histogram_quantile(0.50, sum by (namespace, le) (rate(otc_scrape_duration_seconds_bucket[5m])))`,
					Legend: "{{namespace}}"},
				{Metric: "otc_scrape_duration_seconds_bucket", Title: "Scrape Duration p95", Unit: "s", Type: grafana.TimeSeries,
					Expr:   `histogram_quantile(0.95, sum by (namespace, le) (rate(otc_scrape_duration_seconds_bucket[5m])))`,
					Legend: "{{namespace}}"},
			}},
			{Title: "Scrape Results", Panels: []grafana.PanelConfig{
				{Metric: "otc_scrape_families_count", Title: "Families per Namespace", Unit: "short", Type: grafana.TimeSeries,
					Expr:   `otc_scrape_families_count`,
					Legend: "{{namespace}}"},
				{Metric: "otc_scrape_metrics_count", Title: "Metrics per Namespace", Unit: "short", Type: grafana.TimeSeries,
					Expr:   `otc_scrape_metrics_count`,
					Legend: "{{namespace}}"},
			}},
			{Title: "HTTP Request Duration", Panels: []grafana.PanelConfig{
				{Metric: "otc_http_request_duration_seconds_bucket", Title: "Request Duration p50", Unit: "s", Type: grafana.TimeSeries,
					Expr:   `histogram_quantile(0.50, sum by (host, le) (rate(otc_http_request_duration_seconds_bucket[5m])))`,
					Legend: "{{host}}"},
				{Metric: "otc_http_request_duration_seconds_bucket", Title: "Request Duration p95", Unit: "s", Type: grafana.TimeSeries,
					Expr:   `histogram_quantile(0.95, sum by (host, le) (rate(otc_http_request_duration_seconds_bucket[5m])))`,
					Legend: "{{host}}"},
			}},
			{Title: "HTTP Connection Phases", Panels: []grafana.PanelConfig{
				{Metric: "otc_http_dns_duration_seconds_bucket", Title: "DNS Lookup p95", Unit: "s", Type: grafana.TimeSeries,
					Expr:   `histogram_quantile(0.95, sum by (host, le) (rate(otc_http_dns_duration_seconds_bucket[5m])))`,
					Legend: "{{host}}"},
				{Metric: "otc_http_tls_duration_seconds_bucket", Title: "TLS Handshake p95", Unit: "s", Type: grafana.TimeSeries,
					Expr:   `histogram_quantile(0.95, sum by (host, le) (rate(otc_http_tls_duration_seconds_bucket[5m])))`,
					Legend: "{{host}}"},
				{Metric: "otc_http_ttfb_duration_seconds_bucket", Title: "Time to First Byte p95", Unit: "s", Type: grafana.TimeSeries,
					Expr:   `histogram_quantile(0.95, sum by (host, method, le) (rate(otc_http_ttfb_duration_seconds_bucket[5m])))`,
					Legend: "{{host}} {{method}}"},
			}},
			{Title: "HTTP Connections", Panels: []grafana.PanelConfig{
				{Metric: "otc_http_connections_reused_total", Title: "Connections Reused", Unit: "ops", Type: grafana.TimeSeries,
					Expr:   `rate(otc_http_connections_reused_total[5m])`,
					Legend: "{{host}}"},
				{Metric: "otc_http_connections_new_total", Title: "New Connections", Unit: "ops", Type: grafana.TimeSeries,
					Expr:   `rate(otc_http_connections_new_total[5m])`,
					Legend: "{{host}}"},
			}},
			{Title: "Go Runtime", Panels: []grafana.PanelConfig{
				{Metric: "go_goroutines", Title: "Goroutines", Unit: "short", Type: grafana.TimeSeries,
					Expr: `go_goroutines`},
				{Metric: "go_memstats_alloc_bytes", Title: "Memory Allocated", Unit: "bytes", Type: grafana.TimeSeries,
					Expr: `go_memstats_alloc_bytes`},
				{Metric: "go_gc_duration_seconds", Title: "GC Duration", Unit: "s", Type: grafana.TimeSeries,
					Expr: `go_gc_duration_seconds`},
			}},
			{Title: "Process", Panels: []grafana.PanelConfig{
				{Metric: "process_cpu_seconds_total", Title: "CPU Usage", Unit: "s", Type: grafana.TimeSeries,
					Expr: `rate(process_cpu_seconds_total[5m])`},
				{Metric: "process_resident_memory_bytes", Title: "Resident Memory", Unit: "bytes", Type: grafana.TimeSeries,
					Expr: `process_resident_memory_bytes`},
				{Metric: "process_open_fds", Title: "Open File Descriptors", Unit: "short", Type: grafana.TimeSeries,
					Expr: `process_open_fds`},
			}},
		},
	}
}
