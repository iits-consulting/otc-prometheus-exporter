package internal

import (
	"fmt"
	"os"
	"strconv"
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

	namespaces, ok := os.LookupEnv("NAMESPACES")
	if !ok {
		panic("NAMESPACES not set\n")
	}
	fmt.Println(namespaces)

	port := defaultPort

	rawport, ok := os.LookupEnv("PORT")
	if ok {
		port, err = strconv.Atoi(rawport)
		if err != nil {
			panic(err)
		}
		fmt.Print(port)

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
		Namespaces:      []string{namespaces},
		OtcProjectId:    project.Id,
		OtcProjectToken: project.ScopedToken.Secret,
		Port:            port,
		WaitDuration:    waitDuration,
	}
}
