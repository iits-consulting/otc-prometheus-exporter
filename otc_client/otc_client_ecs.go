package otc_client

import (
	"encoding/json"
	"io"
	"net/http"
)

type EcsResponse struct {
	Servers []Server `json:"servers"`
}

type Links struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Server struct {
	Name string `json:"name"`
	Link []Links
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Id   string `json:"id"`
}

func (o OtcClient) GetEcsData() (*EcsResponse, error) {
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

	var response EcsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}
