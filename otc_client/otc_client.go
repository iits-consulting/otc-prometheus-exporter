package otc_client

import (
	"fmt"
)

const (
	EcsEndpointTemplate      = "https://ecs.eu-de.otc.t-systems.com/v2.1/%s/servers"
	RdsEndpointTemplate      = "https://rds.eu-de.otc.t-systems.com/v3/%s/instances"
	DmsEndpointTemplate      = "https://dms.eu-de.otc.t-systems.com/v2/%s/instances"
	NatEndpointTemplate      = "https://nat.eu-de.otc.t-systems.com/v2.0/nat_gateways?tenant_id=%s"
	ElbEndpointTemplate      = "https://elb.eu-de.otc.t-systems.com/v1.0/%s/elbaas/loadbalancers"
	MetricsEndpointTemplate  = "https://ces.eu-de.otc.t-systems.com/V1.0/%s/metrics"
	CloudEyeEndpointTemplate = "https://ces.eu-de.otc.t-systems.com/V1.0/%s/metric-data"
)

type OtcClient struct {
	secret           string
	ecsEndpoint      string
	rdsEndpoint      string
	natEndpoint      string
	dmsEndpoint      string
	elbEndpoint      string
	metricsEndpoint  string
	cloudEyeEndpoint string
}

func NewOtcClient(projectId, secret string) OtcClient {
	ecsEndpoint := fmt.Sprintf(EcsEndpointTemplate, projectId)
	rdsEndpoint := fmt.Sprintf(RdsEndpointTemplate, projectId)
	dmsEndpoint := fmt.Sprintf(DmsEndpointTemplate, projectId)
	natEndpoint := fmt.Sprintf(NatEndpointTemplate, projectId)
	elbEndpoint := fmt.Sprintf(ElbEndpointTemplate, projectId)
	metricsEndpoint := fmt.Sprintf(MetricsEndpointTemplate, projectId)
	cloudEyeEndpoint := fmt.Sprintf(CloudEyeEndpointTemplate, projectId)

	return OtcClient{
		secret:           secret,
		ecsEndpoint:      ecsEndpoint,
		rdsEndpoint:      rdsEndpoint,
		natEndpoint:      natEndpoint,
		elbEndpoint:      elbEndpoint,
		metricsEndpoint:  metricsEndpoint,
		cloudEyeEndpoint: cloudEyeEndpoint,
		dmsEndpoint:      dmsEndpoint,
	}
}
