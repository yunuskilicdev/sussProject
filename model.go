package main

import "gorm.io/gorm"

type EventLog struct {
	gorm.Model
	UserName      string
	Ip            uint32
	IpRaw         string
	Uuid          string
	UnixTimestamp int64
	Latitude      float64
	Longitude     float64
	Radius        int64
}

type CityBlockResponse struct {
	Latitude       float64
	Longitude      float64
	AccuracyRadius int64
}

type SussRequest struct {
	Username      string `json:"username"`
	UnixTimestamp int64  `json:"unix_timestamp"`
	EventUuid     string `json:"event_uuid"`
	IpAddress     string `json:"ip_address"`
}

type SussResponse struct {
	CurrentGeo                     CurrentGeo `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool       `json:"travelToCurrentGeoSuspicious"`
	TravelFromCurrentGeoSuspicious bool       `json:"travelFromCurrentGeoSuspicious"`
	PrecedingIPAccess              IPAccess   `json:"precedingIpAccess"`
	SubsequentIPAccess             IPAccess   `json:"subsequentIpAccess"`
}

type CurrentGeo struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius int64   `json:"radius"`
}

type IPAccess struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Radius    int64   `json:"radius"`
	Speed     float64 `json:"speed"`
	IP        string  `json:"ip"`
	Timestamp int64   `json:"timestamp"`
}
