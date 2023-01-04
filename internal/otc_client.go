package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Servers []Server
}

/*type Name struct {
	Name string `json:"name"`
}*/

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
	Metrics []Metric
}

type Dimension struct {
	Name  string `json:"Name"`
	Value string `json:"value"`
}

type Metric struct {
	Namespace  string `json:"namespace"`
	Dimensions []Dimension
	MetricName string `json:"metric_name"`
	Unit       string `json:"unit"`
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

//func

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
	json.Unmarshal(body, &response)
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
	json.Unmarshal(body, &response)
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

// Period = 300; from & to berechnen = akutelle Zeit in millisekunden - 1 Sekunde, aktuelle Zeit in millisekunden
// filter = average
func (o OtcClient) GetMetricData(namespace, metricname, dimesionkey, dimensionvalue string) (*map[string]interface{}, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", o.cloudEyeEndpoint, nil)
	req.Header.Add("X-Auth-Token", o.secret)
	q := req.URL.Query()
	q.Add("namespace", namespace)
	q.Add("metric_name", metricname)
	q.Add("from", "1672845283629")
	q.Add("to", "1672845284629")
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

	var response map[string]interface{}
	json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, err
}
