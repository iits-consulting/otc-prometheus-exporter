package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	otcMetricData "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metricdata"
	otcMetrics "github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/metrics"
	otcCompute "github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	ddsInstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/instances"
	dmsInstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/instances"
	elbInstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	natgatewayInstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/natgateways"
	rdsInstances "github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"golang.org/x/exp/slices"
)

type OtcWrapper struct {
	providerClient *golangsdk.ProviderClient
}

func NewOtcClientFromConfig(config ConfigStruct) (*OtcWrapper, error) {
	var opts golangsdk.AuthOptionsProvider = config.AuthenticationData.ToOtcGopherAuthOptionsProvider()
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	return &OtcWrapper{providerClient: provider}, nil
}

func (c *OtcWrapper) GetMetrics() ([]otcMetrics.MetricInfoList, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	cesClient, err := openstack.NewCESClient(c.providerClient, opts)
	if err != nil {
		return []otcMetrics.MetricInfoList{}, err
	}
	metricsResponsePages, err := otcMetrics.ListMetrics(cesClient, otcMetrics.ListMetricsRequest{}).AllPages()
	if err != nil {
		return []otcMetrics.MetricInfoList{}, err
	}
	metricsResponse, err := otcMetrics.ExtractMetrics(metricsResponsePages)
	if err != nil {
		return []otcMetrics.MetricInfoList{}, err
	}

	return metricsResponse.Metrics, nil
}

func (c *OtcWrapper) GetEcsIdNameMapping() (map[string]string, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	computeClient, err := openstack.NewComputeV2(c.providerClient, opts)
	if err != nil {
		return nil, err
	}

	computeResponsePages, err := otcCompute.List(computeClient, otcCompute.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}

	computeResponse, err := otcCompute.ExtractServers(computeResponsePages)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, instance := range computeResponse {
		result[instance.ID] = instance.Name
	}

	return result, nil
}

func (c *OtcWrapper) GetRdsIdNameMapping() (map[string]string, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	rdsClient, err := openstack.NewRDSV3(c.providerClient, opts)
	if err != nil {
		return nil, err
	}
	rdsResponse, err := rdsInstances.List(rdsClient, rdsInstances.ListOpts{})
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, instance := range rdsResponse.Instances {
		result[instance.Id] = instance.Name
	}

	return result, nil
}

func (c *OtcWrapper) GetDmsIdNameMapping() (map[string]string, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	dmsClient, err := openstack.NewDMSServiceV1(c.providerClient, opts)
	if err != nil {
		return nil, err
	}
	dmsResponsePages, err := dmsInstances.List(dmsClient, dmsInstances.ListDmsInstanceOpts{}).AllPages()
	if err != nil {
		return nil, err
	}

	dmsResponse, err := dmsInstances.ExtractDmsInstances(dmsResponsePages)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, instance := range dmsResponse.Instances {
		result[instance.InstanceID] = instance.Name
	}

	return result, nil
}

func (c *OtcWrapper) GetNatIdNameMapping() (map[string]string, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	natClient, err := openstack.NewNatV2(c.providerClient, opts)
	if err != nil {
		return nil, err
	}

	natResponsePages, err := natgatewayInstances.List(natClient, natgatewayInstances.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}

	natResponse, err := natgatewayInstances.ExtractNatGateways(natResponsePages)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, instance := range natResponse {
		result[instance.ID] = instance.Name
	}

	return result, nil
}

func (c *OtcWrapper) GetElbIdNameMapping() (map[string]string, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	elbClient, err := openstack.NewELBV2(c.providerClient, opts)
	if err != nil {
		return nil, err
	}

	elbResponsePages, err := elbInstances.List(elbClient, elbInstances.ListOpts{}).AllPages()
	if err != nil {
		fmt.Println(string(err.(golangsdk.ErrDefault404).Body))
		return nil, err
	}

	elbResponses, err := elbInstances.ExtractLoadBalancers(elbResponsePages)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, instance := range elbResponses {
		result[instance.ID] = instance.Name
	}

	return result, nil
}

// Add in the DDS Service client here (openstack client.go line 719) - David
func (c *OtcWrapper) GetDdsIdNameMapping() (map[string]string, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	ddsClient, err := openstack.NewDDSServiceV3(c.providerClient, opts)

	if err != nil {
		return nil, err
	}

	ddsListResponse, err := ddsInstances.List(ddsClient, ddsInstances.ListInstanceOpts{})

	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, instance := range ddsListResponse.Instances {
		result[instance.Groups[0].Nodes[0].Id] = instance.Name
	}

	return result, nil
}

func (c *OtcWrapper) GetMetricData(metric otcMetrics.MetricInfoList) (*otcMetricData.MetricData, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	cesClient, err := openstack.NewCESClient(c.providerClient, opts)
	if err != nil {
		return nil, err
	}

	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Minute)

	dim0Formatted := fmt.Sprintf("%s,%s", metric.Dimensions[0].Name, metric.Dimensions[0].Value)
	metricData, err := otcMetricData.ShowMetricData(
		cesClient,
		otcMetricData.ShowMetricDataOpts{
			Namespace:  metric.Namespace,
			MetricName: metric.MetricName,
			Dim0:       dim0Formatted,
			From:       strconv.FormatInt(startTime.UnixMilli(), 10),
			To:         strconv.FormatInt(endTime.UnixMilli(), 10),
			Filter:     "average",
			Period:     1,
		},
	)

	if err != nil {
		return nil, err
	}
	return metricData, nil
}

func (c *OtcWrapper) GetMetricDataBatched(metrics []otcMetrics.MetricInfoList) ([]otcMetricData.BatchMetricData, error) {
	opts := golangsdk.EndpointOpts{Region: "eu-de"}
	cesClient, err := openstack.NewCESClient(c.providerClient, opts)
	if err != nil {
		fmt.Print("Error creating CES Client!\n")
		return []otcMetricData.BatchMetricData{}, err
	}

	// cesClient.ResourceBase = "http://localhost:3001/"
	// log.Printf(cesClient.Endpoint)

	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Minute)

	requestedMetrics := make([]otcMetricData.Metric, len(metrics))

	for i, m := range metrics {
		requestedMetrics[i] = otcMetricData.Metric{
			Namespace:  m.Namespace,
			MetricName: m.MetricName,
			Dimensions: []otcMetricData.MetricsDimension{
				{
					Name:  m.Dimensions[0].Name,
					Value: m.Dimensions[0].Value,
				},
			},
		}
	}

	metricData, err := otcMetricData.BatchListMetricData(
		cesClient,
		otcMetricData.BatchListMetricDataOpts{
			Metrics: requestedMetrics,
			From:    startTime.UnixMilli(),
			To:      endTime.UnixMilli(),
			Filter:  "average",
			Period:  "1",
		},
	)

	return metricData, err
}

func FilterByNamespaces(metrics []otcMetrics.MetricInfoList, namespaces []string) []otcMetrics.MetricInfoList {
	var filteredMetrics = []otcMetrics.MetricInfoList{}
	for _, m := range metrics {
		if IsFromNamespace(m, namespaces) {
			filteredMetrics = append(filteredMetrics, m)
		}
	}
	return filteredMetrics
}

func IsFromNamespace(metric otcMetrics.MetricInfoList, namespaces []string) bool {
	return slices.Contains(namespaces, metric.Namespace)
}

func StandardPrometheusMetricName(metric otcMetrics.MetricInfoList) string {
	return fmt.Sprintf(
		"%s_%s",
		strings.TrimPrefix(strings.ToLower(metric.Namespace), "sys."),
		strings.ToLower(metric.MetricName),
	)
}

func StandardPrometheusBatchMetricName(metric otcMetricData.BatchMetricData) string {
	return fmt.Sprintf(
		"%s_%s",
		strings.TrimPrefix(strings.ToLower(metric.Namespace), "sys."),
		strings.ToLower(metric.MetricName),
	)
}
