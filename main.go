package main

import (
	"errors"
	"fmt"
	"os"

	cmdPackage "github.com/iits-consulting/otc-prometheus-exporter/cmd"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"time"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

func collectMetricsInBackground() {
	go func() {
		client, err := internal.NewOtcClientFromConfig(internal.Config)
		if err != nil {
			panic(err)
		}

		var resourceIdToName map[string]string

		if internal.Config.ResourceIdNameMappingFlag {
			resourceIdToName, err = FetchResourceIdToNameMapping(client, internal.Config.Namespaces)
			if err != nil {
				panic(err)
			}
		}

		metrics, err := client.GetMetrics()
		if err != nil {
			panic(err)
		}

		filteredMetrics := internal.FilterByNamespaces(metrics, internal.Config.Namespaces)

		internal.PrometheusMetrics = internal.SetupPrometheusMetricsFromOtcMetrics(filteredMetrics)

		for {
			batchedMetricsResponse, err := client.GetMetricDataBatched(filteredMetrics)
			if err != nil {
				panic(err)
			}

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
			}

			time.Sleep(internal.Config.WaitDuration)
		}
	}()
}

func FetchResourceIdToNameMapping(client *internal.OtcWrapper, namespaces []string) (map[string]string, error) {
	resourceIdToName := make(map[string]string)

	if slices.Contains(namespaces, internal.EcsNamespace) {
		result, err := client.GetEcsIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.RdsNamespace) {
		result, err := client.GetRdsIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.DmsNamespace) {
		result, err := client.GetDmsIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.NatNamespace) {
		result, err := client.GetNatIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}

	if slices.Contains(namespaces, internal.ElbNamespace) {
		result, err := client.GetElbIdNameMapping()
		if err != nil {
			return map[string]string{}, err
		}
		maps.Copy(resourceIdToName, result)
	}
	fmt.Printf("Collected %d resources\n", len(resourceIdToName))
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
		waitDuration          time.Duration
	)

	var rootCmd = &cobra.Command{
		Use: "otc-prometheus-exporter",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(region)
			return cmdPackage.InitializeConfig(cmd, map[string]string{
				"port":                  "PORT",
				"region":                "REGION",
				"namespaces":            "NAMESPACES",
				"os_Username":           "Os_USERNAME",
				"os_Password":           "Os_PASSWORD",
				"os_access_key":         "OS_ACCESS_KEY",
				"os_secret_key":         "OS_SECRET_KEY",
				"projectId":             "OS_PROJECT_ID",
				"osDomainName":          "OS_DOMAIN_NAME",
				"waitDuration":          "WAITDURATION",
				"fetchResourceIdToname": "FETCH_RESOURCE_ID_TO_NAME",
			})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			isAkSkAuthentication := false
			switch {
			case username != "" && password != "":
				isAkSkAuthentication = false
			case accessKey != "" && secretKey != "":
				isAkSkAuthentication = true
			default:
				return errors.New("no valid authentication data provided. please provide either username and password or accessKey and secretKey")
			}
			fmt.Println(isAkSkAuthentication)
			otcRegion, err := internal.NewOtcRegionFromString(region)
			if err != nil {
				return err
			}
			fmt.Println(otcRegion)
			fmt.Println(">>>>", port)
			fmt.Println(region)
			fmt.Println(namespaces, username, password, accessKey, secretKey, waitDuration, fetchResourceIdToname)
			return nil
		},
	}
	rootCmd.Flags().Uint16VarP(&port, "port", "", 39100, "Port on which metrics are served")
	rootCmd.Flags().StringVarP(&region, "region", "r", "eu-de", "region where your project is located ")
	rootCmd.Flags().StringSliceVarP(&namespaces, "namespaces", "n", maps.Values(internal.OtcNamespacesMapping), "namespaces for instances you want to get the metrics from")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "user in the OTC with access to the API. Must be provided together with password and can't be provided with AK/SK")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "password for the user. Must be provided together with username and can't be provided with AK/SK")
	rootCmd.Flags().StringVarP(&accessKey, "accessKey", "", "", "you can instead of username/password also provide the users AK and SK")
	rootCmd.Flags().StringVarP(&secretKey, "secretKey", "", "", "you can instead of username/password also provide the users AK and SK")
	rootCmd.Flags().StringVarP(&projectId, "projectId", "", "", "project from which the metrics should be gathered")
	rootCmd.Flags().StringVarP(&osDomainName, "osDomainName", "", "", "Domainname/Tenant ID")
	rootCmd.Flags().BoolVarP(&fetchResourceIdToname, "fetchResourceIdToname", "", false, "turns the mapping of resource id to resource name on or off")
	rootCmd.Flags().DurationVarP(&waitDuration, "waitDuration", "", 60*time.Second, "time in seconds between two API call fetches")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// collectMetricsInBackground()
	//
	// http.Handle("/metrics", promhttp.Handler())
	// address := fmt.Sprintf(":%d", internal.Config.Port)
	// err := http.ListenAndServe(address, nil)
	// if err != nil {
	//	panic(err)
	// }
}
