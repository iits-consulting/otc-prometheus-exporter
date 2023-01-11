package internal

import (
	"fmt"
	"os"
)

type ConfigStruct struct {
	Namespaces      []string
	OtcProjectId    string
	OtcProjectToken string
}

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
	}
}
