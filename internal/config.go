package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfigStruct struct {
	Namespaces          []string
	OtcUsername         string
	OtcPassword         string
	OtcDomainName       string
	OtcProjectId        string
	OtcIdentityEndpoint string
	Port                int
	WaitDuration        time.Duration
}

const (
	defaultPort             = 8000
	defaultWaitDuration     = 60 * time.Second
	defaultIdentityEndpoint = "https://iam.eu-de.otc.t-systems.com:443/v3"
)

var Config ConfigStruct

func init() {
	LoadConfig()
}

func LoadConfig() {
	var err error
	namespacesRaw, ok := os.LookupEnv("NAMESPACES")
	if !ok || namespacesRaw == "" {
		panic("NAMESPACES not set or empty\n")
	}

	namespaces := strings.Split(namespacesRaw, ",")
	namespacesProcessed := make([]string, len(namespaces))

	for _, namespace := range namespaces {
		namespacesProcessed = append(namespacesProcessed, WithPrefixIfNotPresent(namespace, "SYS."))
	}

	port := defaultPort
	rawport, ok := os.LookupEnv("PORT")
	if ok {
		port, err = strconv.Atoi(rawport)
		if err != nil {
			panic(fmt.Errorf("it looks like the input for the port '%s' was not a number", rawport))
		}
	}

	waitDuration := defaultWaitDuration
	rawtime, ok := os.LookupEnv("WAITDURATION")
	if ok {
		numSeconds, err := strconv.Atoi(rawtime)
		if err != nil {
			panic(err)
		}

		waitDuration = time.Duration(numSeconds) * time.Second
	}

	otcUsername, ok := os.LookupEnv("OTC_USERNAME")
	if !ok {
		panic("OTC_USERNAME environment variable is not set")
	}
	otcPassword, ok := os.LookupEnv("OTC_PASSWORD")
	if !ok {
		panic("OTC_PASSWORD environment variable is not set")
	}
	otcProjectId, ok := os.LookupEnv("OTC_PROJECT_ID")
	if !ok {
		panic("OTC_PROJECT_ID environment variable is not set")
	}
	otcDomainName, ok := os.LookupEnv("OTC_DOMAIN_NAME")
	if !ok {
		panic("OTC_DOMAIN_NAME environment variable is not set")
	}

	Config = ConfigStruct{
		OtcUsername:         otcUsername,
		OtcPassword:         otcPassword,
		OtcProjectId:        otcProjectId,
		OtcDomainName:       otcDomainName,
		OtcIdentityEndpoint: defaultIdentityEndpoint,
		Namespaces:          namespacesProcessed,
		Port:                port,
		WaitDuration:        waitDuration,
	}
}
