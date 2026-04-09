package otcdoc

import "strings"

// OTCNamespace resolves the OTC namespace string (e.g. "SYS.ECS") for a source.
// It prefers the namespace parsed from the RST file. When the RST has no "Namespace"
// section it falls back to a known-exceptions map, then to "SYS.<UPPER>".
func OTCNamespace(ds DocumentationSource, parsedNs string) string {
	if parsedNs != "" {
		return parsedNs
	}
	// RST files that don't declare a namespace section, or where the OTC namespace
	// differs from the simple "SYS.<upper(namespace)>" pattern.
	exceptions := map[string]string{
		"css": "SYS.ES",   // OTC publishes CSS metrics under SYS.ES
		"ddm": "SYS.DDMS", // RST and Huawei catalog both use SYS.DDMS
	}
	if ns, ok := exceptions[ds.Namespace]; ok {
		return ns
	}
	return "SYS." + strings.ToUpper(ds.Namespace)
}

// PrometheusPrefix converts an OTC namespace (e.g. "SYS.NAT", "SERVICE.BMS") to the
// Prometheus metric name prefix used by PrometheusMetricName (e.g. "nat", "bms").
func PrometheusPrefix(otcNs string) string {
	ns := strings.TrimPrefix(otcNs, "SYS.")
	ns = strings.TrimPrefix(ns, "SERVICE.")
	return strings.ToLower(ns)
}
