package internal

import (
	"strings"

	"github.com/iits-consulting/otc-prometheus-exporter/otc_client"
)

func WithPrefixIfNotPresent(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s
	}
	return p + s
}

func GetEcsIdToNameMapping(response otc_client.EcsResponse) map[string]string {
	resourceIdToName := make(map[string]string)
	for _, s := range response.Servers {
		resourceIdToName[s.Id] = s.Name
	}
	return resourceIdToName
}

func GetRdsIdToNameMapping(response otc_client.RdsResponse) map[string]string {
	resourceIdToName := make(map[string]string)
	for _, s := range response.Instance {
		resourceIdToName[s.Id] = s.Name
	}
	return resourceIdToName
}

func GetDmsIdToNameMapping(response otc_client.DmsResponse) map[string]string {
	resourceIdToName := make(map[string]string)
	for _, s := range response.Instances {
		resourceIdToName[s.InstanceId] = s.Name
	}
	return resourceIdToName
}

func GetNatIdtoNameMapping(response otc_client.NatResponse) map[string]string {
	resourceIdToName := make(map[string]string)
	for _, s := range response.NatGateways {
		resourceIdToName[s.Id] = s.Name
	}
	return resourceIdToName
}

func GetElbIdToNameMapping(response otc_client.ElbResponse) map[string]string {
	resourceIdToName := make(map[string]string)
	for _, s := range response.Loadbalancers {
		resourceIdToName[s.Id] = s.Name
	}
	return resourceIdToName
}
