package domain

import "time"

type Interest struct {
	AppName   string    `json:"appname"`
	ServiceIp string    `json:"serviceIp"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
