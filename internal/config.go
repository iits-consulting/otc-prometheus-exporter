package internal

import (
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"strconv"
	"strings"
	"time"

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

const (
	defaultRegion       = otcRegionEuDe
	defaultPort         = 8000
	defaultWaitDuration = 60 * time.Second
)

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

	return "", fmt.Errorf("Invalid argument %s does not represent a valid region.", region)
}

func (r OtcRegion) IamEndpoint() string {
	return fmt.Sprintf("https://iam.%s.otc.t-systems.com:443/v3", r)
}

var Config ConfigStruct

func init() {
	var err error
	Config, err = LoadConfig()
	if err != nil {
		panic(err)

	}

}

func loadNamespacesFromEnv() ([]string, error) {
	namespacesRaw := os.Getenv("NAMESPACES")
	if namespacesRaw == "" {
		return maps.Keys(OtcNamespacesMapping), nil
	}

	namespaces := strings.Split(namespacesRaw, ",")
	namespacesProcessed := make([]string, len(namespaces))

	for i, namespace := range namespaces {
		namespacesProcessed[i] = namespace
		if fullNamespace, ok := OtcNamespacesMapping[namespace]; ok {
			namespacesProcessed[i] = fullNamespace
		}
	}

	return namespacesProcessed, nil
}

func loadPortFromEnv() (int, error) {
	port := defaultPort
	rawport, ok := os.LookupEnv("PORT")
	if !ok {
		return port, nil
	}
	port, err := strconv.Atoi(rawport)
	if err != nil {
		return 0, fmt.Errorf("input port is not a number. got '%s'", rawport)
	}
	return port, nil

}

func loadWaitDurationFromEnv() (time.Duration, error) {
	waitDuration := defaultWaitDuration
	rawtime, ok := os.LookupEnv("WAITDURATION")

	if !ok {
		return waitDuration, nil
	}

	numSeconds, err := strconv.Atoi(rawtime)
	if err != nil {
		return 0, fmt.Errorf("input duration is not a number. got '%s'", waitDuration)
	}

	waitDuration = time.Duration(numSeconds) * time.Second
	return waitDuration, nil
}

func loadResourceIdNameMappingFlagFromEnv() (bool, error) {
	fetchResourceEnabledRaw, ok := os.LookupEnv("FETCH_RESOURCE_ID_TO_NAME")
	if !ok {
		return false, nil
	}
	fetchResourceEnabled, err := strconv.ParseBool(fetchResourceEnabledRaw)
	if err != nil {
		return false, err
	}
	return fetchResourceEnabled, nil

}

func loadRegionFromEnv() (OtcRegion, error) {
	if region, ok := os.LookupEnv("REGION"); ok {
		return NewOtcRegionFromString(region)
	}

	return defaultRegion, nil
}

func loadAuthenticationDataFromEnv() (*AuthenticationData, error) {

	otcUsername := os.Getenv("OS_USERNAME")
	otcPassword := os.Getenv("OS_PASSWORD")
	otcAccessKey := os.Getenv("OS_ACCESS_KEY")
	otcSecretKey := os.Getenv("OS_SECRET_KEY")

	isAkSkAuthentication := false

	switch {
	case otcUsername != "" && otcPassword != "":
		isAkSkAuthentication = false
	case otcAccessKey != "" && otcSecretKey != "":
		isAkSkAuthentication = true
	default:
		return nil, errors.New("no valid authentication data provided. please provide either \"OS_USERNAME\" and \"OS_PASSWORD\" or \"OS_ACCESS_KEY\" and \"OS_SECRET_KEY\"")
	}

	otcProjectId, projectIdOk := os.LookupEnv("OS_PROJECT_ID")
	if !projectIdOk {
		return nil, errors.New("environment variable \"OS_PROJECT_ID\" is not set")
	}
	otcDomainName, domainNameOk := os.LookupEnv("OS_DOMAIN_NAME")
	if !domainNameOk {
		return nil, errors.New("environment variable \"OS_DOMAIN_NAME\" is not set")
	}

	region, err := loadRegionFromEnv()
	if err != nil {
		return nil, err
	}

	return &AuthenticationData{
		Username:             otcUsername,
		Password:             otcPassword,
		AccessKey:            otcAccessKey,
		SecretKey:            otcSecretKey,
		IsAkSkAuthentication: isAkSkAuthentication,
		ProjectId:            otcProjectId,
		DomainName:           otcDomainName,
		Region:               region,
	}, nil
}

func LoadConfig() (ConfigStruct, error) {
	value, err := loadResourceIdNameMappingFlagFromEnv()
	if err != nil {
		panic(err)
	}
	namespaces, err := loadNamespacesFromEnv()
	if err != nil {
		return ConfigStruct{}, err
	}
	port, err := loadPortFromEnv()
	if err != nil {
		return ConfigStruct{}, err
	}
	waitDuration, err := loadWaitDurationFromEnv()
	if err != nil {
		return ConfigStruct{}, err
	}
	authenticationData, err := loadAuthenticationDataFromEnv()
	if err != nil {
		return ConfigStruct{}, err
	}

	return ConfigStruct{
		AuthenticationData:        *authenticationData,
		Namespaces:                namespaces,
		Port:                      port,
		WaitDuration:              waitDuration,
		ResourceIdNameMappingFlag: value,
	}, nil
}
