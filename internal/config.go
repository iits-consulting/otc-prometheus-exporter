package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfigStruct struct {
	Namespaces      []string
	OtcProjectId    string
	OtcProjectToken string
	Port            int
	WaitDuration    time.Duration
}

const defaultPort = 8000
const defaultWaitDuration = 60 * time.Second

var Config ConfigStruct

func init() {
	const path = "~/.otc-auth-config"
	config, err := LoadConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	projectName, ok := os.LookupEnv("PROJECT_NAME")
	if !ok {
		panic("PROJECT_NAME not set\n")
	}

	namespacesraw, ok := os.LookupEnv("NAMESPACES")
	if !ok {
		panic("NAMESPACES not set\n")
	}

	namespacesarray := strings.Split(namespacesraw, ",")

	namespaces := make([]string, len(namespacesarray))

	for i, namespace := range namespacesarray {
		namespaces[i] = WithPrefixIfNotPresent(namespace, "SYS.")
		namespaces = append(namespaces, namespaces[i])
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
		numseconds, err := strconv.Atoi(rawtime)
		if err != err {
			panic(err)
		}

		waitDuration = time.Duration(numseconds) * time.Second
	}

	project, err := GetProjectByName(*config, projectName)
	if err != nil {
		panic(err)
	}

	valid, _ := project.ScopedToken.IsValidNow()
	if err != nil {
		panic(err)
	}
	if !valid {
		panic("Projecttoken is not valid anymore")
	}

	Config = ConfigStruct{
		Namespaces:      namespaces,
		OtcProjectId:    project.Id,
		OtcProjectToken: project.ScopedToken.Secret,
		Port:            port,
		WaitDuration:    waitDuration,
	}
}
