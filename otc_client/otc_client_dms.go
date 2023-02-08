package otc_client

import (
	"encoding/json"
	"io"
	"net/http"
)

type DmsResponse struct {
	Instances []Instance `json:"instances"`
}

type Instance struct {
	Name                       string   `json:"name"`
	Engine                     string   `json:"engine"`
	EngineVersion              string   `json:"egine_version"`
	Description                string   `json:"description"`
	Specification              string   `json:"specification"`
	StorageSpace               int      `json:"storage_space"`
	PartitionNum               string   `json:"partition_num"`
	UsedStorageSpace           int      `json:"used_storage_space"`
	ConnectAddress             string   `json:"connect_adress"`
	Port                       int      `json:"port"`
	Status                     string   `json:"status"`
	InstanceId                 string   `json:"instance_id"`
	ResourceSpecCode           string   `json:"resource_spec_code"`
	VpcId                      string   `json:"vpc_id"`
	VpcName                    string   `json:"vpc_name"`
	CreatedAt                  string   `json:"created_at"`
	SubnetName                 string   `json:"subnet_name"`
	SubnetCidr                 string   `json:"subnet_cidr"`
	UserId                     string   `json:"user_id"`
	UserName                   string   `json:"user_name"`
	AccesUser                  string   `json:"access_user"`
	MaintainBegin              string   `json:"maintain_begin"`
	MaintainEnd                string   `json:"maintain_end"`
	EnablePublicip             string   `json:"enable_publicip"`
	SslEnable                  string   `json:"ssl_enable"`
	SslTwoWayEnable            string   `json:"sslTwoWayEnable"`
	EnableAutoTopic            string   `json:"enable_auto_topic"`
	Type                       string   `json:"type"`
	ProductId                  string   `json:"product_id"`
	SecurityGroupId            string   `json:"security_group_id"`
	SecurityGroupName          string   `json:"security_group_Name"`
	SubnetId                   string   `json:"subnet_id"`
	AvailableZones             []string `json:"available_zones"`
	TotalStorageSpace          int      `json:"total_storage_space"`
	PublicConnectAdress        string   `json:"public_connect_address"`
	StorageResourceId          string   `json:"storage_resource_id"`
	StorageSpecCode            string   `json:"storage_spec_code"`
	ServiceType                string   `json:"service_type"`
	StorageType                string   `json:"storage_type"`
	RetentionPolicy            string   `json:"retention_policy"`
	KafkaPublicStatus          string   `json:"kafka_public_status"`
	PublicBandwidth            int      `json:"public_bandwidth"`
	CrossVpcInfo               string   `json:"cross_vpc_info"`
	RestConnectAddress         string   `json:"rest_connect_address"`
	PublicBoundwidth           int      `json:"public_boundwidth"`
	MessageQueryInstEnable     bool     `json:"message_query_inst_enable"`
	VpcClientPlain             bool     `json:"vpc_client_plain"`
	SupportFeatures            string   `json:"support_features"`
	PodConnectAddress          string   `json:"pod_connect_address"`
	DiskEncrypted              bool     `json:"disk_encrypted"`
	DiskEncryptedKey           string   `json:"disk_encrypted_key"`
	KafkaPrivateConnectAddress string   `json:"kafka_private_connect_address"`
	PublicAccessEnabled        string   `json:"public_access_enabled"`
	NodeNum                    int      `json:"node_num"`
	EnableAcl                  string   `json:"enable_acl"`
	BrokerNum                  int      `json:"broker_num"`
	Tags                       []Tag    `json:"tag"`
	DrEnable                   bool     `json:"dr_enable"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (o OtcClient) GetDmsData() (*DmsResponse, error) {

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

	var response DmsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
