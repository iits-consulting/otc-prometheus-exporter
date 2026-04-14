package otcdoc

import (
	"strings"
)

// ParseDocumentationPageFromRstBytes parses RST content, extracting the OTC namespace (from the
// "Namespace" section, e.g. "SYS.NAT") and all metrics from grid tables.
// Duplicate metric IDs across multiple tables in the same file are deduplicated, keeping the first.
func ParseDocumentationPageFromRstBytes(data []byte) (DocumentationPage, error) {
	lines := strings.Split(string(data), "\n")

	ns := parseNamespaceFromRst(lines)
	blocks := findRstGridTables(lines)

	seen := map[string]bool{}
	var metrics []MetricDocumentationEntry

	for _, block := range blocks {
		if len(block) < 3 {
			continue
		}
		colBounds := parseColBounds(block[0])
		if len(colBounds) < 3 {
			continue
		}

		headerLines, bodyRows := splitHeaderAndBody(block)
		headers := extractCells(headerLines, colBounds)

		idCol := findColIndex(headers, "metric id", "metrics", "metric", "id")
		nameCol := findColIndexExcluding(headers, idCol, "parameter", "name", "metric name", "metrics name", "metric", "metrics")
		descCol := findColIndex(headers, "description", "meaning")
		unitCol := findColIndex(headers, "unit")

		if idCol < 0 || nameCol < 0 {
			continue // not a metrics table
		}

		for _, rowLines := range bodyRows {
			cells := extractCells(rowLines, colBounds)
			if len(cells) <= idCol || len(cells) <= nameCol {
				continue
			}
			id := strings.TrimSpace(cells[idCol])
			name := strings.TrimSpace(cells[nameCol])
			if id == "" {
				continue
			}
			if seen[id] {
				continue
			}
			seen[id] = true

			var unit string
			if unitCol >= 0 && unitCol < len(cells) {
				u := strings.TrimSpace(cells[unitCol])
				if u != "N/A" {
					unit = u
				}
			} else if descCol >= 0 && descCol < len(cells) {
				unit = extractUnitFromDesc(cells[descCol])
			}

			var desc string
			if descCol >= 0 && descCol < len(cells) {
				desc = strings.TrimSpace(cells[descCol])
			}

			metrics = append(metrics, MetricDocumentationEntry{
				MetricId:    id,
				MetricName:  name,
				Unit:        unit,
				Description: desc,
			})
		}
	}

	return DocumentationPage{Namespace: ns, Metrics: metrics}, nil
}

// parseNamespaceFromRst extracts the OTC namespace (e.g. "SYS.NAT") from the RST
// "Namespace" section heading.
func parseNamespaceFromRst(lines []string) string {
	for i, line := range lines {
		if strings.TrimSpace(line) != "Namespace" {
			continue
		}
		if i+1 >= len(lines) || !isDashLine(lines[i+1]) {
			continue
		}
		for j := i + 2; j < len(lines); j++ {
			t := strings.TrimSpace(lines[j])
			if t != "" {
				return t
			}
		}
	}
	return ""
}

func isDashLine(line string) bool {
	t := strings.TrimSpace(line)
	return len(t) > 0 && strings.Trim(t, "-") == ""
}

// findRstGridTables collects contiguous blocks of RST grid table lines (starting with + or |).
func findRstGridTables(lines []string) [][]string {
	var tables [][]string
	var cur []string

	for _, line := range lines {
		t := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(t, "+") || strings.HasPrefix(t, "|") {
			cur = append(cur, t)
		} else {
			if len(cur) >= 3 {
				tables = append(tables, cur)
			}
			cur = nil
		}
	}
	if len(cur) >= 3 {
		tables = append(tables, cur)
	}
	return tables
}

// parseColBounds returns the positions of '+' in an RST grid table separator line.
func parseColBounds(sepLine string) []int {
	var bounds []int
	for i, c := range sepLine {
		if c == '+' {
			bounds = append(bounds, i)
		}
	}
	return bounds
}

// splitHeaderAndBody splits a table block into header lines and groups of data row lines.
// The first separator line containing '=' marks the end of the header section.
func splitHeaderAndBody(tableLines []string) (headerLines []string, bodyRows [][]string) {
	headerDone := false
	var curRow []string

	for _, line := range tableLines {
		isSep := strings.HasPrefix(line, "+")
		isHeaderSep := isSep && strings.ContainsRune(line, '=')

		if !headerDone {
			if isHeaderSep {
				headerDone = true
			} else if !isSep {
				headerLines = append(headerLines, line)
			}
		} else {
			if isSep {
				if len(curRow) > 0 {
					bodyRows = append(bodyRows, curRow)
					curRow = nil
				}
			} else {
				curRow = append(curRow, line)
			}
		}
	}
	if len(curRow) > 0 {
		bodyRows = append(bodyRows, curRow)
	}
	return
}

// extractCells extracts and joins the text content of each column across content lines.
// Empty rows between content within a cell are treated as paragraph breaks (\n\n),
// preserving the structure of RST table cells that use blank rows to separate sections.
func extractCells(contentLines []string, colBounds []int) []string {
	if len(colBounds) < 2 {
		return nil
	}
	numCols := len(colBounds) - 1
	builders := make([]strings.Builder, numCols)
	lastEmpty := make([]bool, numCols)

	for _, line := range contentLines {
		for col := 0; col < numCols; col++ {
			start := colBounds[col] + 1
			end := colBounds[col+1]
			if start >= len(line) {
				lastEmpty[col] = true
				continue
			}
			if end > len(line) {
				end = len(line)
			}
			part := strings.TrimSpace(line[start:end])
			if part == "" {
				lastEmpty[col] = true
				continue
			}
			if builders[col].Len() > 0 {
				if lastEmpty[col] {
					builders[col].WriteString("\n\n")
				} else {
					builders[col].WriteByte(' ')
				}
			}
			builders[col].WriteString(part)
			lastEmpty[col] = false
		}
	}

	cells := make([]string, len(builders))
	for i, b := range builders {
		cells[i] = b.String()
	}
	return cells
}

// findColIndex returns the index of the first header matching any of the given names
// (case-insensitive). Returns -1 if none match.
func findColIndex(headers []string, names ...string) int {
	return findColIndexExcluding(headers, -1, names...)
}

// findColIndexExcluding is like findColIndex but skips the column at excludeIdx.
// Use this when searching for nameCol to avoid re-matching the idCol.
func findColIndexExcluding(headers []string, excludeIdx int, names ...string) int {
	for i, h := range headers {
		if i == excludeIdx {
			continue
		}
		hl := strings.ToLower(strings.TrimSpace(h))
		for _, name := range names {
			if hl == name {
				return i
			}
		}
	}
	return -1
}

// extractUnitFromDesc parses "Unit: X" from a description cell string.
func extractUnitFromDesc(desc string) string {
	idx := strings.Index(desc, "Unit:")
	if idx < 0 {
		return ""
	}
	rest := strings.TrimSpace(desc[idx+5:])
	fields := strings.Fields(rest)
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
}
