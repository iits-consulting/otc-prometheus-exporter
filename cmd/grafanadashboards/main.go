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
		err := os.MkdirAll(outputPath, 0700)
		if err != nil {
			log.Fatalf("Could not create missing output directory %s because of %s\n", outputPath, err)
		}
	}

	for _, ds := range otcdoc.DocumentationSources {
		resp, err := http.Get(ds.Url)
		if err != nil {
			log.Fatalf("Could not fetch the OTC documentation page for %s on %s because of %s\n", ds.Namespace, ds.Url, err)
		}
		defer resp.Body.Close()

		htmlBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Could not read the OTC documentation page for %s because of %s\n", ds.Namespace, err)
		}

		docpage, err := otcdoc.ParseDocumentationPageFromHtmlBytes(htmlBytes, ds.Namespace)
		if err != nil {
			log.Fatalf("Could not parse the HTML from the OTC documentation page for %s because of %s\n", ds.Namespace, err)
		}

		dashboadTitle := fmt.Sprintf("OTC Prometheus Exporter - %s Dashboard", ds.Description)
		dashboardUid := fmt.Sprintf("otc-prometheus-exporter-%s-dashboard", strings.ToLower(ds.Description))
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
			log.Fatalf("Could not save the generated dashboard for %s because of %s\n", ds.Namespace, err)
		}

		outputFile := path.Join(outputPath, fmt.Sprintf("%s-metrics.json", ds.Description))
		fmt.Println(outputFile)
		err = os.WriteFile(outputFile, b, 0644)
		if err != nil {
			log.Fatalf("Could not write the generated dashboard to %s because of %s", ds.Namespace, err)
		}
	}
}

func main() {

	var outputPath string
	var rootCmd = &cobra.Command{
		Use:   "grafanadashboards",
		Short: "Generates Grafana dashboards for the OTC prometheus exporter.",
		Run: func(cmd *cobra.Command, args []string) {
			processDocumentationPages(outputPath)
		},
	}
	rootCmd.Flags().StringVar(&outputPath, "output-path", "", "Directory where all the dashboards will be written to.")
	rootCmd.MarkFlagRequired("output-path") //nolint:errcheck // will not err because only this one flag is used

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}

}
