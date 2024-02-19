package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/iits-consulting/otc-prometheus-exporter/grafana"
	"github.com/iits-consulting/otc-prometheus-exporter/otcdoc"
)

type GrafanaPanel struct {
	Title string
	Expr  string
}

func main() {

	for _, ds := range otcdoc.DocumentationSources {
		resp, err := http.Get(ds.Url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		htmlBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		docpage, err := otcdoc.ParseDocumentationPageFromHtmlBytes(htmlBytes, ds.Namespace)
		if err != nil {
			panic(err)
		}

		dashboadTitle := fmt.Sprintf("OTC Prometheus Exporter - %s Dashboard", strings.ToUpper(docpage.Namespace))
		dashboardUid := fmt.Sprintf("otc-prometheus-exporter-%s-dashboard", docpage.Namespace)
		board := grafana.NewDefaultDashboard(dashboadTitle, dashboardUid)
		numberColumns := 2
		for i, m := range docpage.Metrics {
			width := 12
			height := 15
			x := (i % numberColumns) * width
			y := (i / numberColumns) * width

			s := grafana.PanelSettings{
				Expr:   fmt.Sprintf("%s_%s", ds.Namespace, m.MetricId),
				Title:  m.MetricName,
				Id:     i,
				X:      x,
				Y:      y,
				Width:  width,
				Height: height,
				Unit:   grafana.ConvertOtcMetricToGrafana(m.Unit),
			}
			board.Panels = append(board.Panels, grafana.NewPanelWithSettings(s))

		}

		b, err := json.Marshal(board)
		if err != nil {
			panic(err)
		}

		filename := fmt.Sprintf("%s-metrics.json", ds.Description)
		os.WriteFile(filename, b, 0644)
	}
}
