package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/iits-consulting/otc-prometheus-exporter/grafana"
	"github.com/iits-consulting/otc-prometheus-exporter/otcdoc"
	"github.com/spf13/cobra"
)

// serviceTitle maps a dashboard key (namespace, optionally suffixed with subcomponent)
// to a human-readable Grafana dashboard title.
var serviceTitle = map[string]string{
	"ecs":            "ECS - Elastic Cloud Server",
	"bms":            "BMS - Bare Metal Server",
	"as":             "AS - Auto Scaling",
	"evs":            "EVS - Elastic Volume Service",
	"sfs":            "SFS - Scalable File Service",
	"efs":            "EFS - Elastic File Service (SFS Turbo)",
	"cbr":            "CBR - Cloud Backup & Recovery",
	"vpc":            "VPC - Virtual Private Cloud",
	"elb":            "ELB - Elastic Load Balancing",
	"nat":            "NAT - NAT Gateway",
	"waf":            "WAF - Web Application Firewall",
	"dms":            "DMS - Distributed Message Service",
	"dcs":            "DCS - Distributed Cache Service",
	"rds-mysql":      "RDS - MySQL",
	"rds-postgres":   "RDS - PostgreSQL",
	"rds-sqlserver":  "RDS - SQL Server",
	"dds":            "DDS - Document Database Service",
	"nosql":          "NoSQL - GaussDB NoSQL",
	"gaussdb":        "GaussDB - GaussDB for MySQL",
	"gaussdbv5":      "GaussDB v5 - GaussDB (for OpenGauss)",
	"dws":            "DWS - Data Warehouse Service",
	"css":            "CSS - Cloud Search Service",
	"obs":            "OBS - Object Storage Service",
	"dcaas":          "DCaaS - Direct Connect",
	"vpn-classic":    "VPN - Classic",
	"vpn-enterprise": "VPN - Enterprise Edition",
	"apic":           "APIC - API Gateway",
	"ddm":            "DDM - Distributed Database Middleware",
}

// formatDescription converts remaining RST artifacts in description strings
// to markdown equivalents for display in Grafana panel descriptions.
func formatDescription(s string) string {
	s = strings.ReplaceAll(s, ".. note::", "\n\n**Note:**")
	return strings.TrimSpace(s)
}

// grafanaUnit converts OTC documentation unit strings to Grafana unit identifiers.
func grafanaUnit(otcUnit string) string {
	switch strings.ToLower(strings.TrimSpace(otcUnit)) {
	case "%", "percent":
		return "percent"
	case "ms", "milliseconds", "millisecond":
		return "ms"
	case "s", "seconds", "second":
		return "s"
	case "byte/s", "bytes/s":
		return "Bps"
	case "kb/s", "kbyte/s", "kbytes/s":
		return "KBs"
	case "mb/s", "mbyte/s", "mbytes/s":
		return "MBs"
	case "bit/s", "bits/s":
		return "bps"
	case "kbit/s", "kbits/s":
		return "Kbits"
	case "mbit/s", "mbits/s":
		return "Mbits"
	case "byte", "bytes":
		return "bytes"
	case "kb", "kilobyte", "kilobytes":
		return "kbytes"
	case "mb", "megabyte", "megabytes":
		return "mbytes"
	case "gb", "gigabyte", "gigabytes":
		return "gbytes"
	case "count/s", "counts/s", "req/s", "ops/s":
		return "ops"
	default:
		return "short"
	}
}

// dashboardKey returns the lookup key for a source: namespace, or namespace-subcomponent.
func dashboardKey(ds otcdoc.DocumentationSource) string {
	if ds.SubComponent != "" {
		return ds.Namespace + "-" + strings.ToLower(ds.SubComponent)
	}
	return ds.Namespace
}

// buildAutoConfigs fetches OTC documentation for each source and builds a DashboardConfig.
func buildAutoConfigs() []grafana.DashboardConfig {
	var configs []grafana.DashboardConfig

	for _, ds := range otcdoc.DocumentationSources {
		key := dashboardKey(ds)
		uid := "otc-" + key

		title, ok := serviceTitle[key]
		if !ok {
			title = strings.ToUpper(ds.Namespace)
			if ds.SubComponent != "" {
				title += " - " + ds.SubComponent
			}
		}

		fmt.Printf("fetching %s ...\n", ds.GithubRawUrl)
		page, err := otcdoc.FetchDocumentationSource(ds)
		if err != nil {
			log.Printf("WARN: skipping %s: %v\n", key, err)
			continue
		}

		ns := otcdoc.OTCNamespace(ds, page.Namespace)
		prefix := otcdoc.PrometheusPrefix(ns)

		seen := map[string]bool{}
		var panels []grafana.PanelConfig

		for _, m := range page.Metrics {
			if m.MetricId == "" {
				continue
			}
			id := strings.ToLower(m.MetricId)
			seen[id] = true
			panelTitle := m.MetricName
			if panelTitle == "" {
				panelTitle = m.MetricId
			}
			panels = append(panels, grafana.PanelConfig{
				Metric:      prefix + "_" + id,
				Title:       panelTitle,
				Description: formatDescription(m.Description),
				Unit:        grafanaUnit(m.Unit),
				Type:        grafana.TimeSeries,
			})
		}
		fmt.Printf("  → %s: %d panels from RST\n", key, len(panels))

		// Add Huawei-only metrics (unit available, no description).
		if ds.HuaweiFallbackUrl != "" {
			fmt.Printf("fetching huawei catalog %s ...\n", ds.HuaweiFallbackUrl)
			huaweiMetrics, err := otcdoc.FetchMarkdownMetrics(ds.HuaweiFallbackUrl)
			if err != nil {
				log.Printf("WARN: huawei fallback for %s: %v\n", key, err)
			} else {
				var extra []grafana.PanelConfig
				for _, m := range huaweiMetrics {
					id := strings.ToLower(m.MetricId)
					if seen[id] || id == "" {
						continue
					}
					seen[id] = true
					panelTitle := m.MetricName
					if panelTitle == "" {
						panelTitle = m.MetricId
					}
					extra = append(extra, grafana.PanelConfig{
						Metric:      prefix + "_" + id,
						Title:       panelTitle,
						Description: formatDescription(m.Description),
						Unit:        grafanaUnit(m.Unit),
						Type:        grafana.TimeSeries,
					})
				}
				sort.Slice(extra, func(i, j int) bool { return extra[i].Metric < extra[j].Metric })
				panels = append(panels, extra...)
				fmt.Printf("  → %s: +%d panels from Huawei catalog\n", key, len(extra))
			}
		}

		if len(panels) == 0 {
			log.Printf("WARN: no metrics for %s, skipping\n", key)
			continue
		}

		configs = append(configs, grafana.DashboardConfig{
			Title: title,
			UID:   uid,
			Sections: []grafana.PanelSection{
				{Title: "Metrics", Panels: panels},
			},
		})
	}

	return configs
}

func main() {
	var outputPath string
	var rootCmd = &cobra.Command{
		Use:   "grafanadashboards",
		Short: "Generates Grafana dashboards from OTC documentation.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				log.Fatalf("Could not create output directory %s: %v\n", outputPath, err)
			}

			// Auto-generated dashboards from OTC documentation.
			configs := buildAutoConfigs()
			// Hand-crafted dashboards for metrics that aren't CES-sourced.
			configs = append(configs, exporterDashboard(), alarmDashboard(), statusDashboard())

			for _, cfg := range configs {
				board := grafana.GenerateDashboard(cfg)

				b, err := json.MarshalIndent(board, "", "  ")
				if err != nil {
					log.Fatalf("Could not marshal dashboard %s: %v\n", cfg.Title, err)
				}

				filename := strings.ToLower(strings.ReplaceAll(cfg.UID, "otc-", "")) + ".json"
				outputFile := path.Join(outputPath, filename)
				if err := os.WriteFile(outputFile, b, 0644); err != nil {
					log.Fatalf("Could not write %s: %v\n", outputFile, err)
				}
				fmt.Printf("Generated %s\n", outputFile)
			}
		},
	}
	rootCmd.Flags().StringVar(&outputPath, "output-path", "", "Directory for generated dashboards.")
	rootCmd.MarkFlagRequired("output-path") //nolint:errcheck

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
