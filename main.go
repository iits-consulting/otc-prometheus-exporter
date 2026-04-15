package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	cmdPackage "github.com/iits-consulting/otc-prometheus-exporter/cmd"
	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/iits-consulting/otc-prometheus-exporter/otcclient"
	"github.com/iits-consulting/otc-prometheus-exporter/provider"
	"github.com/spf13/cobra"
)

// registerAllProviders registers all metric providers with the given registry.
func registerAllProviders(registry *provider.Registry) {
	registry.Register(&provider.ECSProvider{})
	registry.Register(&provider.BMSProvider{})
	registry.Register(&provider.EVSProvider{})
	registry.Register(&provider.CSSProvider{})
	registry.Register(&provider.RDSProvider{})
	registry.Register(&provider.ELBProvider{})
	registry.Register(&provider.DMSProvider{})
	registry.Register(&provider.NATProvider{})
	registry.Register(&provider.DCSProvider{})
	registry.Register(&provider.DDSProvider{})
	registry.Register(&provider.VPCProvider{})
	registry.Register(&provider.OBSProvider{})
	registry.Register(&provider.DWSProvider{})
	registry.Register(&provider.SFSProvider{})
	registry.Register(&provider.EFSProvider{})
	registry.Register(&provider.DCaaSProvider{})
	registry.Register(&provider.CBRProvider{})
	registry.Register(&provider.ASProvider{})
	registry.Register(&provider.AlarmProvider{})
}

func main() {
	var (
		port            uint16
		region          string
		username        string
		password        string
		accessKey       string
		secretKey       string
		projectID       string
		regionProjectID string
		osDomainName    string
		logLevel        string
		cesBatchSize    int
		aomBatchSize    int
		aomConcurrency  int
		requestTimeout  int
		idleConnTimeout int
		cesLookback     int
		collectTimeout  int
	)

	var rootCmd = &cobra.Command{
		Use:     "otc-prometheus-exporter",
		Short:   "Exports OTC Cloud Eye metrics as Prometheus metrics",
		Version: Version,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cmdPackage.InitializeConfig(cmd, map[string]string{
				"port":              "PORT",
				"region":            "REGION",
				"os-username":       "OS_USERNAME",
				"os-password":       "OS_PASSWORD",
				"access-key":        "OS_ACCESS_KEY",
				"secret-key":        "OS_SECRET_KEY",
				"project-id":        "OS_PROJECT_ID",
				"region-project-id": "OS_REGION_PROJECT_ID",
				"os-domain-name":    "OS_DOMAIN_NAME",
				"log-level":         "LOG_LEVEL",
				"ces-batch-size":    "CES_BATCH_SIZE",
				"aom-batch-size":    "AOM_BATCH_SIZE",
				"aom-concurrency":   "AOM_CONCURRENCY",
				"request-timeout":   "REQUEST_TIMEOUT",
				"idle-conn-timeout": "IDLE_CONN_TIMEOUT",
				"ces-lookback":      "CES_LOOKBACK",
				"collect-timeout":   "COLLECT_TIMEOUT",
			})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := internal.NewLogger(logLevel)

			//nolint:errcheck // best-effort flush
			defer logger.Sync()

			// Validate authentication credentials.
			hasUserPass := username != "" && password != ""
			hasAKSK := accessKey != "" && secretKey != ""
			if !hasUserPass && !hasAKSK {
				return errors.New("no valid authentication data provided; provide either username/password or access-key/secret-key")
			}

			client, err := otcclient.New(otcclient.Config{
				Username:        username,
				Password:        password,
				AccessKey:       accessKey,
				SecretKey:       secretKey,
				ProjectID:       projectID,
				DomainName:      osDomainName,
				Region:          region,
				RequestTimeout:  time.Duration(requestTimeout) * time.Second,
				IdleConnTimeout: time.Duration(idleConnTimeout) * time.Second,
			}, logger)
			if err != nil {
				return fmt.Errorf("creating OTC client: %w", err)
			}

			if regionProjectID != "" {
				if err := client.SetRegionProjectID(regionProjectID); err != nil {
					return fmt.Errorf("region project: %w", err)
				}
			} else {
				err := client.DiscoverRegionProjectID()
				if err != nil {
					return err
				}
			}

			provider.Config.CESBatchSize = cesBatchSize
			provider.Config.AOMBatchSize = aomBatchSize
			provider.Config.AOMConcurrency = aomConcurrency
			provider.Config.CESLookback = time.Duration(cesLookback) * time.Minute
			provider.Config.CollectTimeout = time.Duration(collectTimeout) * time.Second

			registry := provider.NewRegistry()
			registerAllProviders(registry)

			exporterReg := newExporterRegistry()

			mux := http.NewServeMux()
			mux.HandleFunc("/metrics", metricsHandler(registry, client, logger, exporterReg))
			mux.HandleFunc("/healthz", healthzHandler)

			address := fmt.Sprintf(":%d", port)
			logger.Info("starting server", "address", address)
			return http.ListenAndServe(address, mux)
		},
	}

	rootCmd.Flags().Uint16VarP(&port, "port", "", 39100, "Port on which metrics are served")
	rootCmd.Flags().StringVarP(&region, "region", "r", "eu-de", "OTC region")
	rootCmd.Flags().StringVarP(&username, "os-username", "u", "", "OTC username (user/password auth)")
	rootCmd.Flags().StringVarP(&password, "os-password", "p", "", "OTC password (user/password auth)")
	rootCmd.Flags().StringVarP(&accessKey, "access-key", "", "", "OTC access key (AK/SK auth)")
	rootCmd.Flags().StringVarP(&secretKey, "secret-key", "", "", "OTC secret key (AK/SK auth)")
	rootCmd.Flags().StringVarP(&projectID, "project-id", "", "", "OTC project ID")
	rootCmd.Flags().StringVarP(&regionProjectID, "region-project-id", "", "", "OTC region-level project ID for global services (OBS). If not set, auto-discovered.")
	rootCmd.Flags().StringVarP(&osDomainName, "os-domain-name", "", "", "OTC domain name / tenant ID")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	rootCmd.Flags().IntVarP(&cesBatchSize, "ces-batch-size", "", 500, "Max metrics per CES batch API request")
	rootCmd.Flags().IntVarP(&aomBatchSize, "aom-batch-size", "", 20, "Max metrics per AOM data API request")
	rootCmd.Flags().IntVarP(&aomConcurrency, "aom-concurrency", "", 5, "Max concurrent AOM API calls per scrape")
	rootCmd.Flags().IntVarP(&requestTimeout, "request-timeout", "", 10, "HTTP request timeout in seconds for OTC API calls")
	rootCmd.Flags().IntVarP(&idleConnTimeout, "idle-conn-timeout", "", 90, "How long idle HTTP connections stay in the pool in seconds (set higher than scrape interval)")
	rootCmd.Flags().IntVarP(&cesLookback, "ces-lookback", "", 10, "CES lookback window in minutes for metric data queries")
	rootCmd.Flags().IntVarP(&collectTimeout, "collect-timeout", "", 55, "Maximum collect time in seconds per scrape (should be less than Prometheus scrape_timeout)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
