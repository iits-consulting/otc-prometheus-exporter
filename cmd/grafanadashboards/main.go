package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
		defer func() {
			err := resp.Body.Close()
			log.Fatalf("Could not close the response body of %s because of %s\n", ds.Url, err)
		}()

		htmlBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Could not read the OTC documentation page for %s because of %s\n", ds.Namespace, err)
		}

		docpage, err := otcdoc.ParseDocumentationPageFromHtmlBytes(htmlBytes, ds.Namespace)
		if err != nil {
			log.Fatalf("Could not parse the HTML from the OTC documentation page for %s because of %s\n", ds.Namespace, err)
		}

		dashboadTitle := grafana.OtcSouceDescToGraranaDashboardTitle(ds)
		dashboardUid := grafana.OtcSourceDescToGrafanaUID(ds)
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

		outputFile := path.Join(outputPath, grafana.OtcSourceDescToFilename(ds))
		fmt.Println(outputFile)
		err = os.WriteFile(outputFile, b, 0644)
		if err != nil {
			log.Fatalf("Could not write the generated dashboard to %s because of %s", ds.Namespace, err)
		}
	}
}
func modifyDashboardFiles(directoryPath string) error {
	// List all files in the directory
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		// Check if the file is a JSON file
		if filepath.Ext(file.Name()) == ".json" {
			// Construct the full path to the JSON file
			filePath := filepath.Join(directoryPath, file.Name())

			// Read the content of the JSON file
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", filePath, err)
			}

			// Modify the JSON content to set legendFormat to {{ resource_id }}
			modifiedData := strings.ReplaceAll(string(data), `"legendFormat":"__auto"`, `"legendFormat":"{{ resource_id }}"`)

			// Write the modified content back to the JSON file
			err = ioutil.WriteFile(filePath, []byte(modifiedData), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filePath, err)
			}
		}
	}

	return nil
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

	directoryPath := "otc-prometheus-exporter/charts/otc-prometheus-exporter/dashboards" // set generic path or path in project, correct path?
	err := modifyDashboardFiles(directoryPath)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Dashboard files modified successfully.")
	}

}
