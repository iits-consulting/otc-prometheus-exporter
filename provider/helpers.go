package provider

import (
	"context"
	"sort"
	"strings"

	dto "github.com/prometheus/client_model/go"
)

type contextKey string

const enrichKey contextKey = "enrich"

// WithEnrich returns a context that carries the enrich flag.
func WithEnrich(ctx context.Context, enrich bool) context.Context {
	return context.WithValue(ctx, enrichKey, enrich)
}

// ShouldEnrich returns whether providers should call service-specific APIs.
// Defaults to true if not set in context. Controlled by the ?enrich= query param.
//
// Note: API-only providers (AS, CBR, ALARM) that have no CES metrics return
// nil when this is false, which effectively disables them entirely.
func ShouldEnrich(ctx context.Context) bool {
	v, ok := ctx.Value(enrichKey).(bool)
	if !ok {
		return true
	}
	return v
}

// NewGaugeMetricFamily creates a MetricFamily of type GAUGE with the given name and metrics.
func NewGaugeMetricFamily(name string, metrics []*dto.Metric) *dto.MetricFamily {
	return &dto.MetricFamily{
		Name:   &name,
		Type:   dto.MetricType_GAUGE.Enum(),
		Metric: metrics,
	}
}

// NewGaugeMetric creates a Metric with a gauge value and sorted label pairs.
// Label pairs are sorted alphabetically by key name as required by Prometheus.
func NewGaugeMetric(value float64, labels map[string]string) *dto.Metric {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]*dto.LabelPair, 0, len(labels))
	for _, k := range keys {
		name, labelValue := k, labels[k]
		pairs = append(pairs, &dto.LabelPair{
			Name:  &name,
			Value: &labelValue,
		})
	}

	return &dto.Metric{
		Label: pairs,
		Gauge: &dto.Gauge{Value: &value},
	}
}

// NewGaugeMetricWithTimestamp creates a Metric with a gauge value, sorted label
// pairs, and an explicit timestamp in milliseconds. This allows Prometheus to
// record the metric at the original measurement time rather than the scrape time.
func NewGaugeMetricWithTimestamp(value float64, labels map[string]string, timestampMs int64) *dto.Metric {
	m := NewGaugeMetric(value, labels)
	m.TimestampMs = &timestampMs
	return m
}

// EnrichWithNames iterates all metrics in all families, finds the "resource_id" label value,
// looks it up in nameMap, and sets the "resource_name" label value.
func EnrichWithNames(families []*dto.MetricFamily, nameMap map[string]string) {
	if nameMap == nil {
		return
	}
	for _, fam := range families {
		for _, m := range fam.Metric {
			var resourceID string
			var nameLabel *dto.LabelPair
			for _, lp := range m.Label {
				if lp.GetName() == "resource_id" {
					resourceID = lp.GetValue()
				}
				if lp.GetName() == "resource_name" {
					nameLabel = lp
				}
			}
			if nameLabel != nil && resourceID != "" {
				if name, ok := nameMap[resourceID]; ok {
					nameLabel.Value = &name
				}
			}
		}
	}
}

// FillResourceNameFromLabel copies the value of sourceLabel into resource_name
// for any metric where resource_name is empty. Useful for services like OBS
// where CES provides a descriptive dimension (e.g. bucket_name) but no
// service-API enrichment is available.
func FillResourceNameFromLabel(families []*dto.MetricFamily, sourceLabel string) {
	for _, fam := range families {
		for _, m := range fam.Metric {
			var nameLabel *dto.LabelPair
			var sourceValue string
			for _, lp := range m.Label {
				if lp.GetName() == "resource_name" {
					nameLabel = lp
				}
				if lp.GetName() == sourceLabel {
					sourceValue = lp.GetValue()
				}
			}
			if nameLabel != nil && nameLabel.GetValue() == "" && sourceValue != "" {
				nameLabel.Value = &sourceValue
			}
		}
	}
}

// EnrichWithHelp fills the Help field on any MetricFamily whose Help is currently
// empty, using the generated MetricHelpStrings map. Called once per scrape after
// Collect returns, so no individual provider needs to know about the help strings.
func EnrichWithHelp(families []*dto.MetricFamily) {
	for _, fam := range families {
		if fam.Help != nil && *fam.Help != "" {
			continue
		}
		if help, ok := MetricHelpStrings[fam.GetName()]; ok && help != "" {
			fam.Help = &help
		}
	}
}

// PrometheusMetricName converts an OTC CES namespace and metric name into a
// Prometheus-style metric name. It strips "SYS." or "SERVICE." prefixes from the
// namespace, lowercases everything, and joins with an underscore.
// Example: ("SYS.ECS", "cpu_util") -> "ecs_cpu_util"
func PrometheusMetricName(namespace, metricName string) string {
	// Strip known prefixes.
	ns := namespace
	ns = strings.TrimPrefix(ns, "SYS.")
	ns = strings.TrimPrefix(ns, "SERVICE.")

	return strings.ToLower(ns + "_" + metricName)
}
