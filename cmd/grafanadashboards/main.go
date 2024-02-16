package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iits-consulting/otc-prometheus-exporter/grafana"
)

type MetricDocumentationEntry struct {
	MetricId        string
	Parameter       string
	Description     string
	Unit            string
	Formula         string
	ValueRange      string
	MonitoredObject string
	MonitorinPeriod string
}

type GrafanaPanel struct {
	Title string
	Expr  string
}

var urls = map[string]string{
	"ecs": "https://docs.otc.t-systems.com/elastic-cloud-server/umn/monitoring/basic_ecs_metrics.html",
	"dcs": "https://docs.otc.t-systems.com/distributed-cache-service/umn/monitoring/dcs_metrics.html",
}

var metricMapping = map[string]string{
	"Percent":         "percent",
	"byte/s":          "Bps",
	"request/s":       "reqps",
	"%":               "percent",
	"kbit/s":          "KBs",
	"KB/s":            "KBs",
	"packages/second": "pkgps", // TODO: verify on correctness
}

func parseUnitFromDescription(sDescription *goquery.Selection) string {
	unit := ""
	sDescription.Find("p").Each(func(i int, s *goquery.Selection) {
		if strings.HasPrefix(s.Text(), "Unit:") {
			unit, _ = strings.CutPrefix(s.Text(), "Unit:")
			unit = strings.TrimSpace(unit)
		}
	})

	if unit == "" {
		fmt.Println(
			"could not find the unit in the description", sDescription.Text(),
		)
	}
	return unit
}

func processTableRow(sTableRow *goquery.Selection) (MetricDocumentationEntry, error) {
	m := MetricDocumentationEntry{}
	var err error = nil
	sTableRow.Find("td").Each(func(i int, sTableData *goquery.Selection) {
		switch i {
		case 0:
			m.MetricId = sTableData.Text()
		case 1:
			m.Parameter = sTableData.Text()
		case 2:
			m.Description = sTableData.Text()
			m.Unit = parseUnitFromDescription(sTableData)
		case 3:
			m.ValueRange = sTableData.Text()
		case 4:
			m.MonitoredObject = sTableData.Text()
		case 5:
			m.MonitorinPeriod = sTableData.Text()
		default:
			err = errors.New("invalid amout of columns, the input data or the parsing do not work as expected")
		}
	})
	return m, err
}

func checkIfTableHasIncorrectStructure(sTable *goquery.Selection) bool {
	columnCounter := 0
	sTable.Find("thead tr").First().Find("th").Each(func(i int, sTableHeader *goquery.Selection) {
		fmt.Println(sTableHeader.Text())
		columnCounter++
	},
	)

	return columnCounter != 6
}

func parseMetricEntriesFromDocumentationPage(htmlBytes []byte) ([]MetricDocumentationEntry, error) {
	metricRows := make([]MetricDocumentationEntry, 0)

	reader := bytes.NewReader(htmlBytes)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}

	doc.Find("table").Each(func(i int, sTable *goquery.Selection) {

		if checkIfTableHasIncorrectStructure(sTable) {
			return
		}

		sTable.Find("tbody tr").Each(
			func(i int, sTableRow *goquery.Selection) {
				m, err := processTableRow(sTableRow)
				if err != nil {
					panic(err)
				}
				metricRows = append(metricRows, m)
			},
		)
	})

	return metricRows, nil
}

func main() {

	for namespace, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		htmlBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		metricRows, err := parseMetricEntriesFromDocumentationPage(htmlBytes)
		if err != nil {
			panic(err)
		}

		dashboadTitle := fmt.Sprintf("OTC Prometheus Exporter - %s Dashboard", strings.ToUpper(namespace))
		dashboardUid := fmt.Sprintf("otc-prometheus-exporter-%s-dashboard", namespace)
		board := grafana.NewDefaultDashboard(dashboadTitle, dashboardUid)
		numberColumns := 2
		for i, m := range metricRows {
			width := 12
			height := 15
			x := (i % numberColumns) * width
			y := (i / numberColumns) * width
			mappedMetric := metricMapping[m.Unit]
			if mappedMetric == "" {
				fmt.Printf("Need mapping for \"%s\"\n", m.Unit)
			}

			s := grafana.PanelSettings{
				Expr:   fmt.Sprintf("%s_%s", namespace, m.MetricId),
				Title:  m.Parameter,
				Id:     i,
				X:      x,
				Y:      y,
				Width:  width,
				Height: height,
				Unit:   metricMapping[m.Unit],
			}
			board.Panels = append(board.Panels, grafana.NewPanelWithSettings(s))

		}

		b, err := json.Marshal(board)
		if err != nil {
			panic(err)
		}

		filename := fmt.Sprintf("%s-metrics.json", namespace)
		os.WriteFile(filename, b, 0644)
	}
}
