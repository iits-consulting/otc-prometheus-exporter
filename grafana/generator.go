package grafana

import (
	"strings"
)

// GenerateDashboard converts a DashboardConfig into a Grafana dashboard struct.
func GenerateDashboard(cfg DashboardConfig) Dashboard {
	board := NewDefaultDashboard(cfg.Title, cfg.UID)

	// Add resource_name variable scoped to all metrics for this namespace.
	// Uses a regex like {__name__=~"ecs_.*"} to match any metric from this
	// namespace, so all resources appear regardless of which specific metric
	// they have data for.
	if prefix := metricPrefix(cfg); prefix != "" {
		board.Templating.List = append(board.Templating.List, TemplatingVariable{
			Current:     Current{Selected: true, Text: "All", Value: "$__all"},
			Datasource:  &Datasource{Type: "prometheus", UID: "${datasource}"},
			Description: "Filter by resource name",
			Hide:        0,
			IncludeAll:  true,
			Label:       "Resource",
			Multi:       true,
			Name:        "resource_name",
			Options:     []any{},
			Query:       `label_values({__name__=~"` + prefix + `.*"}, resource_name)`,
			Refresh:     2,
			SkipURLSync: false,
			Type:        "query",
		})
	}

	panelID := 1
	y := 0

	for _, section := range cfg.Sections {
		// Add a row panel for the section
		board.Panels = append(board.Panels, newRowPanel(section.Title, panelID, y))
		panelID++
		y++

		// Lay out panels using a running y offset.
		// Full-width panels (Table) occupy an entire row; half-width panels
		// (all other types) are placed two per row in a 2-column grid.
		panelY := y
		col := 0  // tracks which column (0 or 1) the next half-width panel goes in
		rowH := 0 // height of the current half-width row in progress
		for _, pc := range section.Panels {
			w := 12
			h := 8
			if pc.Type == Stat || pc.Type == Gauge {
				h = 4
			}
			if pc.Type == Table {
				// Full-width: flush any open half-width row first
				if col > 0 {
					panelY += rowH
					col = 0
					rowH = 0
				}
				panel := newMetricPanel(pc, panelID, 0, panelY, 24, h)
				board.Panels = append(board.Panels, panel)
				panelID++
				panelY += h
			} else {
				x := col * 12
				if col == 0 {
					rowH = h
				}
				panel := newMetricPanel(pc, panelID, x, panelY, w, h)
				board.Panels = append(board.Panels, panel)
				panelID++
				col++
				if col == 2 {
					panelY += rowH
					col = 0
					rowH = 0
				}
			}
		}
		// Flush any remaining open half-width row
		if col > 0 {
			panelY += rowH
		}
		y = panelY
	}

	return board
}

// metricPrefix extracts the namespace prefix from the first metric name.
// e.g., "ecs_instance_status" -> "ecs_", "obs_request_count_get_per_second" -> "obs_"
func metricPrefix(cfg DashboardConfig) string {
	for _, s := range cfg.Sections {
		for _, p := range s.Panels {
			if p.Metric != "" {
				if idx := strings.Index(p.Metric, "_"); idx > 0 {
					return p.Metric[:idx+1]
				}
				return p.Metric
			}
		}
	}
	return ""
}

func newRowPanel(title string, id, y int) Panel {
	return Panel{
		Title:   title,
		Type:    "row",
		ID:      id,
		GridPos: GridPos{H: 1, W: 24, X: 0, Y: y},
	}
}

func newMetricPanel(pc PanelConfig, id, x, y, w, h int) Panel {
	return Panel{
		Datasource: Datasource{Type: "prometheus", UID: "${datasource}"},
		FieldConfig: FieldConfig{
			Defaults: Defaults{
				Color:  Color{Mode: "palette-classic"},
				Custom: defaultCustom(),
				Mappings: func() []any {
					if pc.Mappings != nil {
						return pc.Mappings
					}
					return []any{}
				}(),
				Thresholds: Thresholds{
					Mode:  "absolute",
					Steps: convertThresholds(pc.Thresholds),
				},
				Unit:      pc.Unit,
				UnitScale: true,
			},
		},
		GridPos: GridPos{H: h, W: w, X: x, Y: y},
		ID:      id,
		Options: Options{
			Legend:  Legend{Calcs: []any{}, DisplayMode: "list", Placement: "bottom", ShowLegend: true},
			Tooltip: Tooltip{Mode: "single", Sort: "none"},
		},
		Description: pc.Description,
		Targets: []Target{{
			Datasource:          Datasource{Type: "prometheus", UID: "${datasource}"},
			EditorMode:          "builder",
			Expr:                panelExpr(pc),
			IncludeNullMetadata: true,
			Instant:             pc.Type == Stat || pc.Type == Table,
			LegendFormat:        legendFormat(pc),
			Range:               pc.Type == TimeSeries || pc.Type == Gauge,
			RefID:               "A",
		}},
		Title: pc.Title,
		Type:  panelTypeString(pc.Type),
	}
}

func panelExpr(pc PanelConfig) string {
	if pc.Expr != "" {
		return pc.Expr
	}
	return pc.Metric + `{resource_name=~"$resource_name"}`
}

func legendFormat(pc PanelConfig) string {
	if pc.Legend != "" {
		return pc.Legend
	}
	return "{{resource_name}}"
}

func panelTypeString(pt PanelType) string {
	switch pt {
	case TimeSeries:
		return "timeseries"
	case Stat:
		return "stat"
	case Gauge:
		return "gauge"
	case Table:
		return "table"
	default:
		return "timeseries"
	}
}

func convertThresholds(thresholds []Threshold) []Steps {
	if len(thresholds) == 0 {
		return []Steps{{Color: "green", Value: nil}}
	}
	steps := make([]Steps, len(thresholds))
	for i, t := range thresholds {
		var val any
		if i == 0 {
			val = nil
		} else {
			val = t.Value
		}
		steps[i] = Steps{Color: t.Color, Value: val}
	}
	return steps
}

func defaultCustom() Custom {
	return Custom{
		AxisBorderShow:    false,
		AxisCenteredZero:  false,
		AxisColorMode:     "text",
		AxisPlacement:     "auto",
		DrawStyle:         "line",
		FillOpacity:       0,
		GradientMode:      "none",
		HideFrom:          HideFrom{},
		InsertNulls:       false,
		LineInterpolation: "linear",
		LineWidth:         1,
		PointSize:         5,
		ScaleDistribution: ScaleDistribution{Type: "linear"},
		ShowPoints:        "auto",
		SpanNulls:         false,
		Stacking:          Stacking{Group: "A", Mode: "none"},
		ThresholdsStyle:   ThresholdsStyle{Mode: "off"},
	}
}
