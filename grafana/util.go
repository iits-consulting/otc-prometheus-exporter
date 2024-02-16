package grafana

type PanelSettings struct {
	Expr   string
	Title  string
	Id     int
	X      int
	Y      int
	Width  int
	Height int
	Unit   string
}

func NewPanelWithSettings(s PanelSettings) Panel {
	return Panel{
		Datasource: Datasource{
			Type: "prometheus",
			UID:  "${DS_PROMETHEUS}",
		},
		FieldConfig: FieldConfig{
			Defaults: Defaults{
				Color: Color{
					Mode: "palette-classic",
				},
				Custom: Custom{
					AxisBorderShow:   false,
					AxisCenteredZero: false,
					AxisColorMode:    "text",
					AxisLabel:        "",
					AxisPlacement:    "auto",
					BarAlignment:     0,
					DrawStyle:        "line",
					FillOpacity:      0,
					GradientMode:     "none",
					HideFrom: HideFrom{
						Legend:  false,
						Tooltip: false,
						Viz:     false,
					},
					InsertNulls:       false,
					LineInterpolation: "linear",
					LineWidth:         1,
					PointSize:         5,
					ScaleDistribution: ScaleDistribution{
						Type: "linear",
					},
					ShowPoints: "auto",
					SpanNulls:  false,
					Stacking: Stacking{
						Group: "A",
						Mode:  "none",
					},
					ThresholdsStyle: ThresholdsStyle{
						Mode: "off",
					},
				},
				Mappings: []any{},
				Thresholds: Thresholds{
					Mode: "absolute",
					Steps: []Steps{
						{Color: "green", Value: nil},
						{Color: "red", Value: 80},
					},
				},
				Unit:      s.Unit,
				UnitScale: true,
			},
		},
		GridPos: GridPos{X: s.X, Y: s.Y, W: s.Width, H: s.Height},
		ID:      s.Id,
		Options: Options{
			Legend: Legend{
				Calcs:       []any{},
				DisplayMode: "list",
				Placement:   "bottom",
				ShowLegend:  true,
			},
			Tooltip: Tooltip{
				Mode: "single",
				Sort: "none",
			},
		},
		Targets: []Target{
			{
				Datasource: Datasource{
					Type: "prometheus",
					UID:  "${DS_PROMETHEUS}",
				},
				DisableTextWrap:     false,
				EditorMode:          "builder",
				Expr:                s.Expr,
				FullMetaSearch:      false,
				IncludeNullMetadata: true,
				Instant:             false,
				LegendFormat:        "__auto",
				Range:               true,
				RefID:               "A",
				UseBackend:          false,
			},
		},
		Title: s.Title,
		Type:  "timeseries",
	}
}

func NewDefaultDashboard(title, uid string) Dashboad {
	return Dashboad{
		Inputs: []Input{
			{
				Name:        "DS_PROMETHEUS",
				Label:       "Prometheus",
				Description: "",
				Type:        "datasource",
				PluginID:    "prometheus",
				PluginName:  "Prometheus",
			},
		},
		Elements: Elements{},
		Requires: []Require{
			{
				Type:    "grafana",
				ID:      "grafana",
				Name:    "Grafana",
				Version: "10.3.3",
			},
			{
				Type:    "datasource",
				ID:      "prometheus",
				Name:    "Prometheus",
				Version: "1.0.0",
			},
			{
				Type:    "panel",
				ID:      "timeseries",
				Name:    "Time series",
				Version: "",
			},
		},
		Annotations: Annotations{
			List: []AnnotationList{
				{
					BuiltIn: 1,
					Datasource: Datasource{
						Type: "grafana",
						UID:  "-- Grafana --",
					},
					Enable:    true,
					Hide:      true,
					IconColor: "rgba(0, 211, 255, 1)",
					Name:      "Annotations & Alerts",
					Type:      "dashboard",
				},
			},
		},
		Editable:             true,
		FiscalYearStartMonth: 0,
		GraphTooltip:         0,
		ID:                   nil,
		Links:                []any{},
		LiveNow:              false,
		Panels:               []Panel{},
		Refresh:              "",
		SchemaVersion:        39,
		Tags:                 []string{"test"},
		Templating: Templating{
			List: []any{},
		},
		Time: Time{
			From: "now-30m",
			To:   "now",
		},
		Timepicker: Timepicker{},
		Timezone:   "browser",
		Title:      title,
		UID:        uid,
		Version:    4,
		WeekStart:  "",
	}
}
