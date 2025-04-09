package models

type TableQueryMessageRequest struct {
	RequesterId   string `json:"requesterId"`
	RoutingPolicy string `json:"routingPolicy"`
	Payload       []byte `json:"payload"`
}

type TableQueryMessageResponse struct {
	RequesterId   string `json:"requesterId"`
	RoutingPolicy string `json:"routingPolicy"`
	Payload       []byte `json:"payload"`
}
