package otc_client

import (
	"fmt"
)

const (
	EcsEndpointTemplate      = "https://ecs.eu-de.otc.t-systems.com/v2.1/%s/servers"
	RdsEndpointTemplate      = "https://rds.eu-de.otc.t-systems.com/v3/%s/instances"
	MetricsEndpointTemplate  = "https://ces.eu-de.otc.t-systems.com/V1.0/%s/metrics"
	CloudEyeEndpointTemplate = "https://ces.eu-de.otc.t-systems.com/V1.0/%s/metric-data"
)

type OtcClient struct {
	secret           string
	ecsEndpoint      string
	rdsEndpoint      string
	metricsEndpoint  string
	cloudEyeEndpoint string
}

func NewOtcClient(projectId, secret string) OtcClient {
	ecsEndpoint := fmt.Sprintf(EcsEndpointTemplate, projectId)
	rdsEndpint := fmt.Sprintf(RdsEndpointTemplate, projectId)
	metricsEndpoint := fmt.Sprintf(MetricsEndpointTemplate, projectId)
	cloudEyeEndpoint := fmt.Sprintf(CloudEyeEndpointTemplate, projectId)
	return OtcClient{
		secret:           secret,
		ecsEndpoint:      ecsEndpoint,
		rdsEndpoint:      rdsEndpint,
		metricsEndpoint:  metricsEndpoint,
		cloudEyeEndpoint: cloudEyeEndpoint,
	}
}
