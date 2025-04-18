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

type JobRouting struct {
	JobName string `json:"job_name" bson:"job_name"`
	// The ServiceIPPriority field is a list, whose entries depict the routing priority of the service instance
	// in comparison to other service instances of the same service balancing policy.
	// The higher the priority, the more likely the service instance will be selected
	// for routing.
	// The priority is a value between 0 and 1, where 0 is the lowest priority and 1 is the highest.
	// The default priority is 0.5.
	ServiceIPPriority []PriorityEntry `json:"service_ip_priority" bson:"service_ip_priority"`
}

type PriorityEntry struct {
	IpType   ServiceIpType `json:"IpType" bson:"IpType"`
	Priority float64       `json:"Priority" bson:"Priority"`
}
