package grafana

type Dashboad struct {
	Inputs               []Input     `json:"__inputs"`
	Elements             Elements    `json:"__elements"`
	Requires             []Require   `json:"__requires"`
	Annotations          Annotations `json:"annotations"`
	Editable             bool        `json:"editable"`
	FiscalYearStartMonth int         `json:"fiscalYearStartMonth"`
	GraphTooltip         int         `json:"graphTooltip"`
	ID                   any         `json:"id"`
	Links                []any       `json:"links"`
	LiveNow              bool        `json:"liveNow"`
	Panels               []Panel     `json:"panels"`
	Refresh              string      `json:"refresh"`
	SchemaVersion        int         `json:"schemaVersion"`
	Tags                 []string    `json:"tags"`
	Templating           Templating  `json:"templating"`
	Time                 Time        `json:"time"`
	Timepicker           Timepicker  `json:"timepicker"`
	Timezone             string      `json:"timezone"`
	Title                string      `json:"title"`
	UID                  string      `json:"uid"`
	Version              int         `json:"version"`
	WeekStart            string      `json:"weekStart"`
}
type Input struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Type        string `json:"type"`
	PluginID    string `json:"pluginId"`
	PluginName  string `json:"pluginName"`
}
type Elements struct {
}
type Require struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}
type Datasource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}
type AnnotationList struct {
	BuiltIn    int        `json:"builtIn"`
	Datasource Datasource `json:"datasource"`
	Enable     bool       `json:"enable"`
	Hide       bool       `json:"hide"`
	IconColor  string     `json:"iconColor"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
}
type Annotations struct {
	List []AnnotationList `json:"list"`
}
type Color struct {
	Mode string `json:"mode"`
}
type HideFrom struct {
	Legend  bool `json:"legend"`
	Tooltip bool `json:"tooltip"`
	Viz     bool `json:"viz"`
}
type ScaleDistribution struct {
	Type string `json:"type"`
}
type Stacking struct {
	Group string `json:"group"`
	Mode  string `json:"mode"`
}
type ThresholdsStyle struct {
	Mode string `json:"mode"`
}
type Custom struct {
	AxisBorderShow    bool              `json:"axisBorderShow"`
	AxisCenteredZero  bool              `json:"axisCenteredZero"`
	AxisColorMode     string            `json:"axisColorMode"`
	AxisLabel         string            `json:"axisLabel"`
	AxisPlacement     string            `json:"axisPlacement"`
	BarAlignment      int               `json:"barAlignment"`
	DrawStyle         string            `json:"drawStyle"`
	FillOpacity       int               `json:"fillOpacity"`
	GradientMode      string            `json:"gradientMode"`
	HideFrom          HideFrom          `json:"hideFrom"`
	InsertNulls       bool              `json:"insertNulls"`
	LineInterpolation string            `json:"lineInterpolation"`
	LineWidth         int               `json:"lineWidth"`
	PointSize         int               `json:"pointSize"`
	ScaleDistribution ScaleDistribution `json:"scaleDistribution"`
	ShowPoints        string            `json:"showPoints"`
	SpanNulls         bool              `json:"spanNulls"`
	Stacking          Stacking          `json:"stacking"`
	ThresholdsStyle   ThresholdsStyle   `json:"thresholdsStyle"`
}
type Steps struct {
	Color string `json:"color"`
	Value any    `json:"value"`
}
type Thresholds struct {
	Mode  string  `json:"mode"`
	Steps []Steps `json:"steps"`
}
type Defaults struct {
	Color      Color      `json:"color"`
	Custom     Custom     `json:"custom"`
	Mappings   []any      `json:"mappings"`
	Thresholds Thresholds `json:"thresholds"`
	Unit       string     `json:"unit"`
	UnitScale  bool       `json:"unitScale"`
}
type FieldConfig struct {
	Defaults  Defaults       `json:"defaults"`
	Overrides map[string]any `json:"overrides"`
}
type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}
type Legend struct {
	Calcs       []any  `json:"calcs"`
	DisplayMode string `json:"displayMode"`
	Placement   string `json:"placement"`
	ShowLegend  bool   `json:"showLegend"`
}
type Tooltip struct {
	Mode string `json:"mode"`
	Sort string `json:"sort"`
}
type Options struct {
	Legend  Legend  `json:"legend"`
	Tooltip Tooltip `json:"tooltip"`
}
type Target struct {
	Datasource          Datasource `json:"datasource"`
	DisableTextWrap     bool       `json:"disableTextWrap"`
	EditorMode          string     `json:"editorMode"`
	Expr                string     `json:"expr"`
	FullMetaSearch      bool       `json:"fullMetaSearch"`
	IncludeNullMetadata bool       `json:"includeNullMetadata"`
	Instant             bool       `json:"instant"`
	LegendFormat        string     `json:"legendFormat"`
	Range               bool       `json:"range"`
	RefID               string     `json:"refId"`
	UseBackend          bool       `json:"useBackend"`
}
type Panel struct {
	Datasource  Datasource  `json:"datasource"`
	FieldConfig FieldConfig `json:"fieldConfig"`
	GridPos     GridPos     `json:"gridPos"`
	ID          int         `json:"id"`
	Options     Options     `json:"options"`
	Targets     []Target    `json:"targets"`
	Title       string      `json:"title"`
	Type        string      `json:"type"`
}
type Templating struct {
	List []any `json:"list"`
}
type Time struct {
	From string `json:"from"`
	To   string `json:"to"`
}
type Timepicker struct {
}
