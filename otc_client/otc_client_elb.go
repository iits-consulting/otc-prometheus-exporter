package otc_client

import (
	"encoding/json"
	"io"
	"net/http"
)

type ElbResponse struct {
	Loadbalancers []Loadbalancer `json:"load_balancers"`
}

type Loadbalancer struct {
	VipAddress      string `json:"vip_address"`
	UpdateTime      string `json:"update_time"`
	CreateTime      string `json:"create_time"`
	Id              string `json:"id"`
	Status          string `json:"status"`
	Bandwidth       int    `json:"bandwidth"`
	VpcId           string `json:"vpc_id"`
	AdminStateUp    int    `json:"admin_state_up"`
	VipSubnetId     string `json:"vip_subnet_id"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SecurityGroupId string `json:"security_group_id"`
}

func (o OtcClient) GetElbData() (*ElbResponse, error) {

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

	var response ElbResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
