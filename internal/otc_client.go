package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const EcsEndpointTemplate = "https://ecs.eu-de.otc.t-systems.com/v2.1/%s/servers"
const MetricsEndpointTemplate = "https://ces.eu-de.otc.t-systems.com/V1.0/%s/metrics"

type OtcClient struct {
	secret          string
	ecsEndpoint     string
	metricsEndpoint string
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
	return OtcClient{
		secret:          secret,
		ecsEndpoint:     ecsEndpoint,
		metricsEndpoint: metricsEndpoint,
	}
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
	json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

func (o OtcClient) GetMetricsData() (*MetricsResponse, error) {
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
