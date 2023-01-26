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
