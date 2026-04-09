package otcdoc

import (
	"testing"
)

const huaweiMarkdown = `
| 维度 | 指标名 | 指标描述 | 指标单位 |
| :-- | :-- | :-- | :-- |
| server | cpu_util | CPU使用率 | % |
| server | mem_util | 内存使用率 | 毫秒 |
| server | net_in | 网络流入 | 次/秒 |
| server | unknown_unit | 未知指标 | 未知单位 |
| server | cpu_util | CPU使用率 | % |
`

func TestParseHuaweiMarkdown(t *testing.T) {
	entries := ParseHuaweiMarkdownMetrics([]byte(huaweiMarkdown))
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries (deduped), got %d", len(entries))
	}
}

func TestHuaweiMarkdownSkipsHeader(t *testing.T) {
	entries := ParseHuaweiMarkdownMetrics([]byte(huaweiMarkdown))
	for _, e := range entries {
		if e.MetricId == "指标名" || e.MetricId == "Metric ID" {
			t.Errorf("header row was not skipped, got entry with ID %q", e.MetricId)
		}
	}
}

func TestHuaweiMarkdownDeduplicates(t *testing.T) {
	entries := ParseHuaweiMarkdownMetrics([]byte(huaweiMarkdown))
	seen := map[string]int{}
	for _, e := range entries {
		seen[e.MetricId]++
	}
	if seen["cpu_util"] != 1 {
		t.Errorf("expected cpu_util to appear once after dedup, got %d", seen["cpu_util"])
	}
}

func TestHuaweiMarkdownTranslatesChineseUnits(t *testing.T) {
	entries := ParseHuaweiMarkdownMetrics([]byte(huaweiMarkdown))
	byID := map[string]string{}
	for _, e := range entries {
		byID[e.MetricId] = e.Unit
	}
	if byID["mem_util"] != "ms" {
		t.Errorf("expected 毫秒 to translate to ms, got %q", byID["mem_util"])
	}
	if byID["net_in"] != "count/s" {
		t.Errorf("expected 次/秒 to translate to count/s, got %q", byID["net_in"])
	}
}

func TestHuaweiMarkdownDiscardsUnknownNonASCIIUnit(t *testing.T) {
	entries := ParseHuaweiMarkdownMetrics([]byte(huaweiMarkdown))
	byID := map[string]string{}
	for _, e := range entries {
		byID[e.MetricId] = e.Unit
	}
	if byID["unknown_unit"] != "" {
		t.Errorf("expected unknown non-ASCII unit to be discarded, got %q", byID["unknown_unit"])
	}
}

func TestHuaweiMarkdownKeepsASCIIUnit(t *testing.T) {
	entries := ParseHuaweiMarkdownMetrics([]byte(huaweiMarkdown))
	byID := map[string]string{}
	for _, e := range entries {
		byID[e.MetricId] = e.Unit
	}
	if byID["cpu_util"] != "%" {
		t.Errorf("expected ASCII unit %% to be kept as-is, got %q", byID["cpu_util"])
	}
}

func TestHuaweiMarkdownSkipsTooFewCells(t *testing.T) {
	// Parser requires at least 5 pipe-split segments (4 content cells).
	// "| only | two |" splits into 4 segments — below threshold, should be skipped.
	input := "| only | two |\n"
	entries := ParseHuaweiMarkdownMetrics([]byte(input))
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for row with too few cells, got %d", len(entries))
	}
}

func TestTranslateUnit(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"ms", "ms"},
		{"%", "%"},
		{"Bytes/s", "Bytes/s"},
		{"毫秒", "ms"},
		{"次/秒", "count/s"},
		{"次/分钟", "count/min"},
		{"未知单位", ""},
	}
	for _, c := range cases {
		got := translateUnit(c.input)
		if got != c.want {
			t.Errorf("translateUnit(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
