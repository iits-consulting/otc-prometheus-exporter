package grafana

import (
	"fmt"
	"strings"

	"github.com/iits-consulting/otc-prometheus-exporter/otcdoc"
)

var metricMapping = map[string]string{
	"percent":     "percent",
	"Percent":     "percent",
	"%":           "percent",
	"percent (%)": "percent",

	"count/s":         "cps",
	"Count/s":         "cps",
	"packages/second": "pps",
	"Packet/s":        "pps",
	"Packets/s":       "pps",
	"request/s":       "reqps",
	"Request/s":       "reqps",
	"Requests/s":      "reqps",
	"Query/s":         "ops",
	"bit/s":           "bps",
	"byte/s":          "Bps",
	"Byte/s":          "Bps",
	"Bytes/s":         "Bps",
	"kbit/s":          "KBs",
	"KB/s":            "KBs",

	"KB/op": "none", // TODO: maybe create a metric for this
	"ms/op": "none", // TODO: maybe create a metric for this

	"ms/count": "ms",
	"ms":       "ms",
	"μs":       "μs",
	"second":   "s",

	"GB":    "decgbytes",
	"byte":  "bytes",
	"count": "none", // none is just number
	"Count": "none", // none is just number
	"N/A":   "none", // none is just number
	"":      "none", // none is just number

}

func ConvertOtcMetricToGrafana(m string) string {
	grafanaMetric, ok := metricMapping[m]
	if !ok {
		fmt.Println("Could not find metric mapping for ", m)
	}
	return grafanaMetric
}

func OtcSouceDescToGraranaDashboardTitle(ds otcdoc.DocumentationSource) string {
	title := "OTC Prometheus Exporter - " + strings.ToUpper(ds.Namespace)
	if ds.SubComponent != "" {
		title += " (" + ds.SubComponent + ")"
	}
	return title
}

func OtcSourceDescToGrafanaUID(ds otcdoc.DocumentationSource) string {
	uid := "otc-" + ds.Namespace
	if ds.SubComponent != "" {
		uid += "-" + strings.ToLower(ds.SubComponent)
	}
	return uid
}

func OtcSourceDescToFilename(ds otcdoc.DocumentationSource) string {
	filename := ds.Namespace + ".json"
	if ds.SubComponent != "" {
		filename = ds.Namespace + "-" + strings.ToLower(ds.SubComponent) + ".json"
	}
	return filename
}

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
			UID:  "${datasource}",
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
					UID:  "${prometheus}",
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
		Tags:                 []string{"OTC"},
		Templating: Templating{
			List: []TemplatingVariable{
				{
					Current: Current{
						Selected: true,
						Text:     "Prometheus",
						Value:    "prometheus",
					},
					Description: "Datasource where the OTC Prometheus Exporter data is stored",
					Hide:        0,
					IncludeAll:  false,
					Label:       "Datasource",
					Multi:       false,
					Name:        "datasource",
					Options:     []any{},
					Query:       "prometheus",
					QueryValue:  "",
					Refresh:     1,
					Regex:       "",
					SkipURLSync: false,
					Type:        "datasource",
				},
			},
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
