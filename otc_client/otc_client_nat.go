package otc_client

import (
	"encoding/json"
	"io"
	"net/http"
)

type NatResponse struct {
	NatGateways []NatGateway `json:"nat_gateways"`
}

type NatGateway struct {
	Id                string `json:"id"`
	TenantId          string `json:"tenant_id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Spec              string `json:"spec"`
	RouterId          string `json:"router_id"`
	InternalNetworkId string `json:"internal_network_id"`
	Status            string `json:"status"`
	AdminStateUp      bool   `json:"admin_state_up"`
	CreatedAt         string `json:"created_at"`
}

func (o OtcClient) GetNatData() (*NatResponse, error) {

	client := http.Client{}
	req, err := http.NewRequest("GET", o.rdsEndpoint, nil)
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

	var response NatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
