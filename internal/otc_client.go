package internal

import (
	"fmt"
	"io"
	"net/http"
)

const endpointTemplate = "https://ecs.eu-de.otc.t-systems.com/v2.1/%s/servers"

type OtcClient struct {
	secret      string
	ecsEndpoint string
}

func NewOtcClient(projectId, secret string) OtcClient {
	url := fmt.Sprintf(endpointTemplate, projectId)
	return OtcClient{
		secret:      secret,
		ecsEndpoint: url,
	}
}

func (o OtcClient) GetEcsData() (*string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", o.ecsEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", o.secret)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stringbody := string(body)
	return &stringbody, nil

}
