package internal

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

type ConfigStruct struct {
	AuthenticationData        AuthenticationData
	Namespaces                []string
	Port                      int
	WaitDuration              time.Duration
	ResourceIdNameMappingFlag bool
}

type AuthenticationData struct {
	Username             string
	Password             string
	AccessKey            string
	SecretKey            string
	IsAkSkAuthentication bool
	ProjectId            string
	DomainName           string
	Region               OtcRegion
}

func (ad AuthenticationData) ToOtcGopherAuthOptionsProvider() golangsdk.AuthOptionsProvider {
	var opts golangsdk.AuthOptionsProvider
	if ad.IsAkSkAuthentication {
		opts = golangsdk.AKSKAuthOptions{
			IdentityEndpoint: ad.Region.IamEndpoint(),
			AccessKey:        ad.AccessKey,
			SecretKey:        ad.SecretKey,
			Domain:           ad.DomainName,
			ProjectId:        ad.ProjectId,
		}
	} else {
		opts = golangsdk.AuthOptions{
			IdentityEndpoint: ad.Region.IamEndpoint(),
			Username:         ad.Username,
			Password:         ad.Password,
			DomainName:       ad.DomainName,
			TenantID:         ad.ProjectId,
			AllowReauth:      true,
		}
	}
	return opts
}

type OtcRegion string

const (
	otcRegionEuDe OtcRegion = "eu-de"
	otcRegionEuNl OtcRegion = "eu-nl"
)

func NewOtcRegionFromString(region string) (OtcRegion, error) {
	otcRegion := OtcRegion(region)
	if slices.Contains([]OtcRegion{otcRegionEuNl, otcRegionEuDe}, otcRegion) {
		return otcRegion, nil
	}

	return "", fmt.Errorf("invalid argument %s does not represent a valid region", region)
}

func (r OtcRegion) IamEndpoint() string {
	return fmt.Sprintf("https://iam.%s.otc.t-systems.com:443/v3", r)
}

// ResolveOtcShortHandNamespace maps the short code for the namespaces to the actual namespace name.
func ResolveOtcShortHandNamespace(namespaces []string) []string {
	fullNamespaces := make([]string, len(namespaces))
	for i, v := range namespaces {
		correctNamespaceName, ok := OtcNamespacesMapping[v]
		fullNamespaces[i] = v
		if ok {
			fullNamespaces[i] = correctNamespaceName
		}
	}

	return fullNamespaces
}
