package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/iits-consulting/otc-prometheus-exporter/internal"
)

func main() {
	const path = "~/.otc-auth-config"
	config, err := internal.LoadConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	const endpointTemplate = "https://ecs.eu-de.otc.t-systems.com/v2.1/%s/servers"

	project, err := internal.GetProjectByName(*config, "eu-de_iits-central")
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("https://ecs.eu-de.otc.t-systems.com/v2.1/%s/servers", project.Id)

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Auth-Token", project.ScopedToken.Secret)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	stringbody := string(body)
	fmt.Println(stringbody)

}
