package main

import (
	"fmt"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
)

func main() {
	const path = "~/.otc-auth-config"
	config, err := internal.LoadConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	project, err := internal.GetProjectByName(*config, "eu-de_iits-central")
	if err != nil {
		panic(err)
	}

	valid, err := project.ScopedToken.IsValidNow()
	if err != nil {
		panic(err)
	}
	if !valid {
		panic("Projecttoken is not valid anymore")
	}

	client := internal.NewOtcClient(project.Id, project.ScopedToken.Secret)

	result, _ := client.GetEcsData()
	fmt.Println(*result)

}
