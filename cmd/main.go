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
	projectToken, err := internal.GetScopedToken(*config, "eu-de_iits-infra")
	if err != nil {
		panic(err)
	}
	fmt.Println(projectToken)
}
