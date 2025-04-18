package domain

type RoutingChange struct {
	AppName              string                  `json:"appName" bson:"appName"`
	ServiceIP            string                  `json:"serviceIp" bson:"serviceIp"`
	InstancePriorityList []InstancePriorityEntry `json:"instancePriorityList" bson:"instancePriorityList"`
}

type InstancePriorityEntry struct {
	InstanceID string  `json:"instanceId" bson:"instanceId"`
	Priority   float64 `json:"priority" bson:"priority"`
}
