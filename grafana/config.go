package grafana

// DashboardConfig describes a Grafana dashboard to generate.
type DashboardConfig struct {
	Title    string
	UID      string
	Sections []PanelSection
}

// PanelSection groups related panels under a collapsible row.
type PanelSection struct {
	Title  string
	Panels []PanelConfig
}

// PanelConfig describes a single Grafana panel.
type PanelConfig struct {
	Metric      string
	Title       string
	Description string
	Unit        string
	Type        PanelType
	Thresholds  []Threshold
	Mappings    []any  // Grafana value mappings (fieldConfig.defaults.mappings)
	Legend      string // Override legend format (default: "{{resource_name}}")
	Expr        string // Override PromQL expression (default: Metric{resource_name=~"$resource_name"})
}

// PanelType selects the Grafana visualization type.
type PanelType int

const (
	TimeSeries PanelType = iota
	Stat
	Gauge
	Table
)

// Threshold defines a color step for Grafana panel thresholds.
type Threshold struct {
	Value float64
	Color string
}
