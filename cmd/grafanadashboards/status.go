package main

import "github.com/iits-consulting/otc-prometheus-exporter/grafana"

// statusDashboard returns the hand-crafted "Service Health Overview" dashboard.
// It aggregates synthesized *_status metrics (0=healthy, 1=abnormal) from all specific
// providers into a single view. These metrics are not in OTC documentation so they
// cannot be auto-generated from RST sources.
func statusDashboard() grafana.DashboardConfig {
	unhealthyThresholds := []grafana.Threshold{
		{Color: "green", Value: 0},
		{Color: "red", Value: 1},
	}

	// Maps 0 → "All healthy" so it's clear the service is monitored and all resources are OK.
	// Values > 0 display as-is (count of unhealthy resources).
	// No data means the service is not monitored / has no resources.
	healthyMapping := []any{
		map[string]any{
			"type": "value",
			"options": map[string]any{
				"0": map[string]any{
					"text":  "All healthy",
					"color": "green",
					"index": 0,
				},
			},
		},
	}

	// unhealthy returns a Stat panel showing count of resources with status==1.
	// Metric is intentionally empty: these custom exprs span multiple namespaces,
	// so no single resource_name template variable would be meaningful.
	unhealthy := func(metric, title string) grafana.PanelConfig {
		return grafana.PanelConfig{
			Title:      title,
			Unit:       "short",
			Type:       grafana.Stat,
			Expr:       `count(` + metric + ` == 1) or vector(0)`,
			Legend:     "Unhealthy",
			Thresholds: unhealthyThresholds,
			Mappings:   healthyMapping,
		}
	}

	return grafana.DashboardConfig{
		Title: "Service Health Overview",
		UID:   "otc-status",
		Sections: []grafana.PanelSection{
			{Title: "Compute", Panels: []grafana.PanelConfig{
				unhealthy("ecs_instance_status", "ECS Instances"),
				unhealthy("bms_instance_status", "BMS Instances"),
			}},
			{Title: "Load Balancing & Networking", Panels: []grafana.PanelConfig{
				unhealthy("elb_loadbalancer_status", "ELB Load Balancers"),
				unhealthy("nat_gateway_status", "NAT Gateways"),
				unhealthy("dcaas_virtual_gateway_status", "DCaaS Virtual Gateways"),
			}},
			{Title: "Database", Panels: []grafana.PanelConfig{
				unhealthy("rds_instance_status", "RDS Instances"),
				unhealthy("rds_node_status", "RDS Nodes"),
				unhealthy("dds_instance_status", "DDS Instances"),
				unhealthy("dds_node_status", "DDS Nodes"),
				unhealthy("dcs_instance_status", "DCS Instances"),
				unhealthy("dws_cluster_status", "DWS Clusters"),
				unhealthy("css_cluster_status", "CSS Clusters"),
			}},
			{Title: "Messaging", Panels: []grafana.PanelConfig{
				unhealthy("dms_instance_status", "DMS Instances"),
			}},
			{Title: "Storage", Panels: []grafana.PanelConfig{
				unhealthy("evs_volume_status", "EVS Volumes"),
				unhealthy("sfs_share_status", "SFS Shares"),
				unhealthy("efs_filesystem_status", "EFS Filesystems"),
			}},
			{Title: "Auto Scaling & Backup", Panels: []grafana.PanelConfig{
				unhealthy("as_group_status", "AS Groups"),
				unhealthy("cbr_backup_status", "CBR Backups"),
			}},
		},
	}
}
