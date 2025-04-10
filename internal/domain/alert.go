package domain

import "time"

type Alert struct {
	AppName   string    `json:"appName" bson:"appname"`
	CreatedAt time.Time `json:"createdAt" bson:"createdat"`
}

type AlertRequest struct {
	AppName string `json:"appName"`
}

type AlertResponse struct {
	AppName   string    `json:"appName"`
	CreatedAt time.Time `json:"createdAt"`
}
