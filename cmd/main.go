package main

import (
	"fmt"
	"os"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
)

// 1. ENV Namespaces auslesen und parsen(Ziel: Array von strings)
// 2. filtern nach Namespaces
// 3. Anfrage
func main() {
	const path = "~/.otc-auth-config"
	config, err := internal.LoadConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	projectName, ok := os.LookupEnv("PROJECT_NAME")
	if !ok {
		panic("PROJECT_NAME not set\n")
	}

	namesspaces, ok := os.LookupEnv("NAMESPACES")
	if !ok {
		panic("NAMESPACES not set\n")
	}
	fmt.Println(namesspaces)

	project, err := internal.GetProjectByName(*config, projectName)
	if err != nil {
		panic(err)
	}

	//fmt.Println(project.ScopedToken.ExpiresAt)
	valid, _ := project.ScopedToken.IsValidNow()
	if err != nil {
		panic(err)
	}
	if !valid {
		panic("Projecttoken is not valid anymore")
	}

	client := internal.NewOtcClient(project.Id, project.ScopedToken.Secret)

	result, _ := client.GetEcsData()

	m := make(map[string]string)

	for _, s := range result.Servers {
		m[s.Id] = s.Name
	}

	metrics, _ := client.GetMetricTypes()
	filterdMetrics := metrics.FilterByNamespaces([]string{namesspaces})

	y, err := client.GetAllMetricData(filterdMetrics)
	if err != nil {
		panic(err)
	}
	fmt.Println(y)

	//MetricNameSpaces := internal.GetMetricNamespaces()
	//fmt.Println(MetricNameSpaces)

	//fmt.Println(metrics)
	//fmt.Println(m["98a3ec65-4da2-437a-b83c-0045186d74ec"])

}
