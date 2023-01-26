package internal

import (
	"github.com/iits-consulting/otc-prometheus-exporter/otc_client"
	"strings"
)

func WithPrefixIfNotPresent(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s
	}
	return p + s
}

func GetEcsIdToNameMapping(response otc_client.EcsResponse) map[string]string {
	// TODO
	return map[string]string{}
}

func GetRdsIdToNameMapping(response otc_client.RdsResponse) map[string]string {
	// TODO
	return map[string]string{}
}
