package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	cmdPackage "github.com/iits-consulting/otc-prometheus-exporter/cmd"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"time"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func collectMetricsInBackground(config internal.ConfigStruct, logger internal.ILogger) {
	go func() {
		client, err := internal.NewOtcClientFromConfig(config, logger)
		if err != nil {
			logger.Panic("Error creating OTC client", "error", err)
		}

		logger.Info("New OTC Client created!")

		var resourceIdToName map[string]string

		if config.ResourceIdNameMappingFlag {
			resourceIdToName, err = FetchResourceIdToNameMapping(client)
			if err != nil {
				logger.Panic("Error mapping resource IDs to Names", "error", err)
			}
			logger.Info("Resource ID names mapped!", "idmaps", resourceIdToName)
		}

		metrics, err := client.GetMetrics()
		if err != nil {
			logger.Panic("Unable to get metrics", "error", err, "context", "Initial get metrics call to get metric baseline.")
		}

		logger.Info("Initial metrics retrieved!")

		filteredMetrics, removedMetrics := internal.FilterByNamespaces(metrics, client.Config.Namespaces)

		logger.Info("Metrics filtered.")

		logger.Debug("Metrics removed.", "removed_metrics", removedMetrics)

		internal.PrometheusMetrics = internal.SetupPrometheusMetricsFromOtcMetrics(filteredMetrics, client.Logger)

		logger.Info("Prometheus metrics initiated!")

		for {
			batchedMetricsResponse, getMetricErr := client.GetMetricDataBatched(filteredMetrics)
			if getMetricErr != nil {
				logger.Error("Unable to retrieve batched metric data", "error", getMetricErr)
				time.Sleep(config.WaitDuration)
				continue
			}

			logger.Info("Batch metrics retrieved.")

			for _, metric := range batchedMetricsResponse {
				if len(metric.Datapoints) == 0 {
					continue
				}
				internal.PrometheusMetrics[internal.StandardPrometheusBatchMetricName(metric)].With(
					prometheus.Labels{
						"unit":          metric.Unit,
						"resource_id":   metric.Dimensions[0].Value,
						"resource_name": resourceIdToName[metric.Dimensions[0].Value],
					}).Set(metric.Datapoints[len(metric.Datapoints)-1].Average)

				logger.Debug("Prometheus metric set", "metric", internal.StandardPrometheusBatchMetricName(metric))
			}

			time.Sleep(config.WaitDuration)
		}
	}()
}

func FetchResourceIdToNameMapping(client *internal.OtcWrapper) (map[string]string, error) {
	resourceIdToName := make(map[string]string)

	if slices.Contains(client.Config.Namespaces, internal.EcsNamespace) {
		result, err := client.GetEcsIdNameMapping()
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "ECS Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(client.Config.Namespaces, internal.RdsNamespace) {
		result, err := client.GetRdsIdNameMapping()
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "RDS Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(client.Config.Namespaces, internal.DmsNamespace) {
		result, err := client.GetDmsIdNameMapping()
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "DMS Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(client.Config.Namespaces, internal.NatNamespace) {
		result, err := client.GetNatIdNameMapping()
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "NAT Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(client.Config.Namespaces, internal.ElbNamespace) {
		result, err := client.GetElbIdNameMapping()
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "ELB Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(client.Config.Namespaces, internal.DdsNamespace) {
		result, err := client.GetDdsIdNameMapping()
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "DDS Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(client.Config.Namespaces, internal.DcsNamespace) {
		result, err := client.GetDcsIdNameMapping()
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "DCS Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(client.Config.Namespaces, internal.VpcNamespace) {
		result, err := client.GetVpcIdNameMapping(client.Config.AuthenticationData.ProjectId)
		if err != nil {
			client.Logger.Error("Unable to map namespace!", "context", "VPC Name mapping")
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	client.Logger.Info(fmt.Sprintf("Collected %d resources\n", len(resourceIdToName)))
	return resourceIdToName, nil
}

func main() {
	var (
		port                  uint16
		region                string
		namespaces            []string
		username              string
		password              string
		accessKey             string
		secretKey             string
		projectId             string
		osDomainName          string
		fetchResourceIdToname bool
		waitDuration          uint
		logLevel              string
	)

	var rootCmd = &cobra.Command{
		Use: "otc-prometheus-exporter",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cmdPackage.InitializeConfig(cmd, map[string]string{
				"port":                      "PORT",
				"region":                    "REGION",
				"namespaces":                "NAMESPACES",
				"os-username":               "OS_USERNAME",
				"os-password":               "OS_PASSWORD",
				"access-key":                "OS_ACCESS_KEY",
				"secret-key":                "OS_SECRET_KEY",
				"projectId":                 "OS_PROJECT_ID",
				"os-domain-name":            "OS_DOMAIN_NAME",
				"wait-duration":             "WAITDURATION",
				"fetch-resource-id-to-name": "FETCH_RESOURCE_ID_TO_NAME",
				"log-level":                 "LOG_LEVEL",
			})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Add main log system, passed as dependency
			logger := internal.NewLogger(logLevel)

			//nolint:errcheck // not relevant
			defer logger.Sync()

			logger.Info("Env settings",
				"port", port,
				"region", region,
				"namespaces", namespaces,
				"waitDuration", waitDuration,
				"fetchResourceIdToName", fetchResourceIdToname,
				"logging", logLevel)

			isAkSkAuthentication := false
			switch {
			case username != "" && password != "":
				isAkSkAuthentication = false
			case accessKey != "" && secretKey != "":
				isAkSkAuthentication = true
			default:
				return errors.New("no valid authentication data provided. please provide either username and password or accessKey and secretKey")
			}

			otcRegion, err := internal.NewOtcRegionFromString(region)
			if err != nil {
				return err
			}

			namespaces = internal.ResolveOtcShortHandNamespace(namespaces, logger)

			collectMetricsInBackground(internal.ConfigStruct{
				Port:                      int(port),
				Namespaces:                namespaces,
				WaitDuration:              time.Duration(waitDuration) * time.Second,
				ResourceIdNameMappingFlag: fetchResourceIdToname,
				AuthenticationData: internal.AuthenticationData{
					Username:             username,
					Password:             password,
					AccessKey:            accessKey,
					SecretKey:            secretKey,
					IsAkSkAuthentication: isAkSkAuthentication,
					ProjectId:            projectId,
					DomainName:           osDomainName,
					Region:               otcRegion,
				},
			}, logger)

			http.Handle("/metrics", promhttp.Handler())
			address := fmt.Sprintf(":%d", port)
			err = http.ListenAndServe(address, nil)
			return err
		},
	}
	rootCmd.Flags().Uint16VarP(&port, "port", "", 39100, "Port on which metrics are served")
	rootCmd.Flags().StringVarP(&region, "region", "r", "eu-de", "region where your project is located ")
	rootCmd.Flags().StringSliceVarP(&namespaces, "namespaces", "n", maps.Values(internal.OtcNamespacesMapping), "namespaces for instances you want to get the metrics from")
	rootCmd.Flags().StringVarP(&username, "os-username", "u", "", "user in the OTC with access to the API. Must be provided together with password and can't be provided with AK/SK")
	rootCmd.Flags().StringVarP(&password, "os-password", "p", "", "password for the user. Must be provided together with username and can't be provided with AK/SK")
	rootCmd.Flags().StringVarP(&accessKey, "access-key", "", "", "you can instead of username/password also provide the users AK and SK")
	rootCmd.Flags().StringVarP(&secretKey, "secret-key", "", "", "you can instead of username/password also provide the users AK and SK")
	rootCmd.Flags().StringVarP(&projectId, "projectId", "", "", "project from which the metrics should be gathered")
	rootCmd.Flags().StringVarP(&osDomainName, "os-domain-name", "", "", "Domainname/Tenant ID")
	rootCmd.Flags().BoolVarP(&fetchResourceIdToname, "fetch-resource-id-to-name", "", false, "turns the mapping of resource id to resource name on or off")
	rootCmd.Flags().UintVarP(&waitDuration, "wait-duration", "", 60, "time in seconds between two API call fetches")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "", "", "type of logging to use")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
