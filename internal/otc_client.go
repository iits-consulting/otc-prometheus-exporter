package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const EcsEndpointTemplate = "https://ecs.eu-de.otc.t-systems.com/v2.1/%s/servers"
const MetricsEndpointTemplate = "https://ces.eu-de.otc.t-systems.com/V1.0/%s/metrics"
const CloudEyeEndpointTemplate = "https://ces.eu-de.otc.t-systems.com/V1.0/%s/metric-data"

type OtcClient struct {
	secret           string
	ecsEndpoint      string
	metricsEndpoint  string
	cloudEyeEndpoint string
}

type EcsResponse struct {
	Servers []Server `json:"servers"`
}

type Links struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Server struct {
	Name string `json:"name"`
	Link []Links
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Id   string `json:"id"`
}

type MetricsResponse struct {
	Metrics []Metric `json:"metrics"`
}

type Dimension struct {
	Name  string `json:"Name"`
	Value string `json:"value"`
}

type Metric struct {
	Namespace  string      `json:"namespace"`
	Dimensions []Dimension `json:"dimensions"`
	MetricName string      `json:"metric_name"`
	Unit       string      `json:"unit"`
}

type Datapoint struct {
	Average   float64 `json:"average"`
	Timestamp int     `json:"timestamp"`
	Unit      string  `json:"unit"`
}

type CloudEyeResponse struct {
	DataPoints []Datapoint `json:"datapoints"`
	MetricName string      `json:"metric_name"`
}

func NewOtcClient(projectId, secret string) OtcClient {
	ecsEndpoint := fmt.Sprintf(EcsEndpointTemplate, projectId)
	metricsEndpoint := fmt.Sprintf(MetricsEndpointTemplate, projectId)
	cloudEyeEndpoint := fmt.Sprintf(CloudEyeEndpointTemplate, projectId)
	return OtcClient{
		secret:           secret,
		ecsEndpoint:      ecsEndpoint,
		metricsEndpoint:  metricsEndpoint,
		cloudEyeEndpoint: cloudEyeEndpoint,
	}
}

func (o OtcClient) GetAllMetricData(mr MetricsResponse) (map[string]CloudEyeResponse, error) {

	const SleepDurationSeconds = 1

	result := map[string]CloudEyeResponse{}
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Second)

	for i, m := range mr.Metrics {
		y, err := o.GetMetricData(m.Namespace, m.MetricName, m.Dimensions[0].Name, m.Dimensions[0].Value, startTime, endTime)
		if err != nil {
			return map[string]CloudEyeResponse{}, err
		}
		result[m.StandardPrometheusMetricName()] = *y
		fmt.Println(i, result[m.StandardPrometheusMetricName()])

		time.Sleep(time.Second * SleepDurationSeconds)
	}
	return result, nil

}

func (o OtcClient) GetEcsData() (*EcsResponse, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", o.ecsEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", o.secret)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response EcsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

func (o OtcClient) GetMetricTypes() (*MetricsResponse, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", o.metricsEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", o.secret)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response MetricsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (mR MetricsResponse) FilterByNamespaces(namespaces []string) MetricsResponse {
	filterdMetrics := []Metric{}
	for _, m := range mR.Metrics {
		if m.IsFromNamespace(namespaces) {
			filterdMetrics = append(filterdMetrics, m)
		}
	}
	return MetricsResponse{filterdMetrics}
}

func (m Metric) IsFromNamespace(namespaces []string) bool {
	for _, n := range namespaces {
		if n == m.Namespace {
			return true
		}
	}
	return false
}

func (m Metric) StandardPrometheusMetricName() string {
	return fmt.Sprintf(
		"%s_%s",
		strings.TrimPrefix(strings.ToLower(m.Namespace), "sys."),
		strings.ToLower(m.MetricName),
	)
}

func (o OtcClient) GetMetricData(namespace, metricname, dimesionkey, dimensionvalue string, startTime time.Time, endTime time.Time) (*CloudEyeResponse, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", o.cloudEyeEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", o.secret)
	q := req.URL.Query()
	q.Add("namespace", namespace)
	q.Add("metric_name", metricname)
	q.Add("from", strconv.Itoa(int(startTime.UnixMilli())))
	q.Add("to", strconv.Itoa(int(endTime.UnixMilli())))
	q.Add("period", "300")
	q.Add("filter", "average")
	q.Add("dim.0", fmt.Sprintf("%s,%s", dimesionkey, dimensionvalue))

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response CloudEyeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, err
}
