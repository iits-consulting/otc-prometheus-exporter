package otc_client

import (
	"encoding/json"
	"io"
	"net/http"
)

type RdsResponse struct {
	Instance   []Instances `json:"instances"`
	TotalCount int         `json:"total_count"`
}

type Datastores struct {
	Type            string `json:"type"`
	Version         string `json:"version"`
	CompleteVersion string `json:"complete_version"`
}

type Volume struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

type BackUpStrategy struct {
	StartTime string `json:"start_time"`
	KeepDays  int    `json:"keep_days"`
}

type Nodes struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Role             string `json:"role"`
	Status           string `json:"status"`
	AvailabilityZone string `json:"availability_zone"`
}

type RelatedInstances struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Tags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ChargeInfo struct {
	ChargeMode string `json:"charge_mode"`
}

type Instances struct {
	Id                  string             `json:"id"`
	Name                string             `json:"name"`
	Status              string             `json:"status"`
	Alias               string             `json:"alias"`
	PrivateIps          []string           `json:"private_ips"`
	PublicIps           []string           `json:"public_ips"`
	Port                int                `json:"port"`
	Type                string             `json:"type"`
	Region              string             `json:"region"`
	Datastore           Datastores         `json:"datastore"`
	Created             string             `json:"created"`
	Updated             string             `json:"updated"`
	DbUserName          string             `json:"db_user_name"`
	VpcId               string             `json:"vpc_id"`
	SubnetId            string             `json:"subnet_id"`
	SecurityGroupId     string             `json:"security_group_id"`
	Cpu                 string             `json:"cpu"`
	Mem                 string             `json:"mem"`
	FlavorRef           string             `json:"flavor_ref"`
	Volume              Volume             `json:"volume"`
	SwitchStrategy      string             `json:"switch_strategy"`
	BackUpStrategy      BackUpStrategy     `json:"backup_strategy"`
	MaintenanceWindow   string             `json:"maintenance_window"`
	Node                []Nodes            `json:"nodes"`
	RelatedInstance     []RelatedInstances `json:"related_instance"`
	DiskEncryptionId    string             `json:"disk_encryption_id"`
	EnterpriseProjectId string             `json:"enterprise_project_id"`
	TimeZone            string             `json:"time_zone"`
	ChargeInfo          ChargeInfo         `json:"charge_info"`
	Tag                 []Tags             `json:"tags"`
	AssociatedWithDdm   bool               `json:"associated_with_ddm"`
}

func (o OtcClient) GetRdsData() (*RdsResponse, error) {

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

	var response RdsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
