package otcdoc

import (
	"strings"
)

// huaweiUnitMap translates Chinese unit strings found in Huawei catalog files to English equivalents.
// Units not present in this map that contain non-ASCII characters are discarded (replaced with "").
var huaweiUnitMap = map[string]string{
	"毫秒":   "ms",
	"次":    "count",
	"次/秒":  "count/s",
	"次/分钟": "count/min",
	"个":    "count",
	"个/秒":  "count/s",
	"包/秒":  "packets/s",
	"积分":   "credits",
}

// translateUnit maps known Chinese unit strings to English. Returns the original value
// unchanged for ASCII units (e.g. "ms", "%", "Bytes/s", "μs"). Returns "" for unknown
// non-ASCII strings.
func translateUnit(unit string) string {
	for _, r := range unit {
		if r > 127 {
			return huaweiUnitMap[unit] // "" if not in map
		}
	}
	return unit
}

// ParseHuaweiMarkdownMetrics parses a Huawei Cloud metrics catalog Markdown file.
// These files have a pipe-delimited table with columns:
//
//	维度 (dimension) | 指标名 (metric ID) | 指标描述 (description, Chinese) | 指标单位 (unit)
//
// Returns one entry per unique metric ID. Chinese descriptions are discarded.
// Chinese unit strings are translated to English via huaweiUnitMap; unknown non-ASCII units are discarded.
func ParseHuaweiMarkdownMetrics(data []byte) []MetricDocumentationEntry {
	seen := map[string]bool{}
	var result []MetricDocumentationEntry

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		// Skip separator rows (e.g. |:--|:--|:--|:--|)
		if strings.Contains(line, "---") || strings.Contains(line, ":--") {
			continue
		}
		cells := strings.Split(line, "|")
		// After split on "|": cells[0]="" cells[1]=dim cells[2]=metric_id cells[3]=desc cells[4]=unit cells[5]=""
		// Column positions are fixed based on the known Huawei catalog format. If the upstream
		// table gains or reorders columns this will silently read the wrong values.
		if len(cells) < 5 {
			continue
		}
		metricID := strings.TrimSpace(cells[2])
		unit := strings.TrimSpace(cells[4])

		// Skip header row
		if metricID == "指标名" || metricID == "Metric ID" || metricID == "" {
			continue
		}
		// Clean any embedded HTML (e.g. <br>)
		metricID = strings.ReplaceAll(metricID, "<br>", "")
		metricID = strings.TrimSpace(metricID)
		unit = strings.ReplaceAll(unit, "<br>", "")
		unit = strings.TrimSpace(unit)

		if metricID == "" || seen[metricID] {
			continue
		}
		seen[metricID] = true
		result = append(result, MetricDocumentationEntry{
			MetricId: metricID,
			Unit:     translateUnit(unit),
		})
	}
	return result
}
