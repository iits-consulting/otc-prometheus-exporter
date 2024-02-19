package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/iits-consulting/otc-prometheus-exporter/grafana"
	"github.com/iits-consulting/otc-prometheus-exporter/otcdoc"
	"github.com/spf13/cobra"
)

func processDocumentationPages(outputPath string) {
	_, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, 0644)
		if err != nil {
			log.Fatalf("Could not create missing output directory %s because of error: %s", outputPath, err)
		}
	}

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

		outputFile := path.Join(outputPath, fmt.Sprintf("%s-metrics.json", ds.Description))
		os.WriteFile(outputFile, b, 0644)
	}
}

func main() {

	var outputPath string
	var rootCmd = &cobra.Command{
		Use:   "grafanadashboards",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			processDocumentationPages(outputPath)
		},
	}
	rootCmd.Flags().StringVar(&outputPath, "output-path", "", "Directory where all the dashboards will be written to.")
	rootCmd.MarkFlagRequired("output-path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
