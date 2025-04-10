package domain

type ServiceIpType string

const (
	ServiceIpTypeRoundRobin ServiceIpType = "RR"
)

type ServiceIpListEntry struct {
	Address   string        `json:"Address" bson:"Address"`
	Addressv6 string        `json:"Address_v6" bson:"Address_v6"`
	IpType    ServiceIpType `json:"IpType" bson:"IpType"`
}

type ServiceInstanceListEntry struct {
	InstanceNumber int    `json:"instance_number" bson:"instance_number"`
	InstanceIP     string `json:"instance_ip" bson:"instance_ip"`
	InstanceIPv6   string `json:"instance_ip_v6" bson:"instance_ip_v6"`
}

type Job struct {
	JobName             string                     `json:"job_name" bson:"job_name"`
	ServiceIpList       []ServiceIpListEntry       `json:"service_ip_list" bson:"service_ip_list"`
	ServiceInstanceList []ServiceInstanceListEntry `json:"instance_list" bson:"instance_list"`
}
