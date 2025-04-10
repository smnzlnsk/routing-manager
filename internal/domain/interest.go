package domain

import "time"

type Interest struct {
	AppName   string    `json:"appname" bson:"appname"`
	ServiceIp string    `json:"serviceIp" bson:"serviceIp"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type InterestRequest struct {
	AppName   string `json:"appname"`
	ServiceIp string `json:"serviceIp"`
}

type InterestResponse struct {
	AppName   string `json:"appname"`
	ServiceIp string `json:"serviceIp"`
	Status    string `json:"status"`
}
