package main

import "github.com/iits-consulting/otc-prometheus-exporter/grafana"

// alarmDashboard returns the hand-crafted dashboard for CES alarms.
// The metric otc_alarm_state is synthesized by the exporter from the OTC Alarm API
// (not a CES metric), so it has no documentation source and cannot be auto-generated.
func alarmDashboard() grafana.DashboardConfig {
	return grafana.DashboardConfig{
		Title: "ALARM - CES Alarms",
		UID:   "otc-alarm",
		Sections: []grafana.PanelSection{
			{Title: "Overview", Panels: []grafana.PanelConfig{
				{Metric: "otc_alarm_state", Title: "Firing Alarms", Unit: "short", Type: grafana.Stat,
					Expr:   `count(otc_alarm_state == 1) or vector(0)`,
					Legend: "Firing"},
			}},
			{Title: "Alarm States", Panels: []grafana.PanelConfig{
				{Metric: "otc_alarm_state", Title: "All Alarms", Unit: "short", Type: grafana.Table},
			}},
		},
	}
}
