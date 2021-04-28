package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
	"net"
	"net/http"
)

const maxSpeed = 400

func main() {
	InitialMigration()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5000", nil)
}

func InitialMigration() {
	dbConn, err := gorm.Open(sqlite.Open("suss.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	dbConn.AutoMigrate(&EventLog{})
}

func getDbConnection() *gorm.DB {
	dbConn, err := gorm.Open(sqlite.Open("suss.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return dbConn
}

func handler(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var requestBody SussRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestBody)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	response, err := handle(requestBody)
	if err != nil {
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handle(request SussRequest) (*SussResponse, error) {
	response := SussResponse{}
	ipLong, err := Ip2long(request.IpAddress)
	if err != nil {
		return nil, err
	}

	dbConn := getDbConnection()
	var cityBlock CityBlockResponse
	dbConn.Raw("select latitude, longitude, accuracy_radius from city_blocks where network_start_integer <= ? and network_last_integer >= ?;", ipLong, ipLong).Scan(&cityBlock)

	currentGeo := CurrentGeo{
		Lat:    cityBlock.Latitude,
		Lon:    cityBlock.Longitude,
		Radius: cityBlock.AccuracyRadius,
	}
	response.CurrentGeo = currentGeo

	var previousEvent EventLog
	preEventSubQuery := dbConn.Select("MAX(unix_timestamp)").Where("user_name = ? and unix_timestamp <= ? and ip != ?", request.Username, request.UnixTimestamp, ipLong).Table("event_logs")
	dbConn.Where("user_name = ? AND ip != ? AND unix_timestamp = (?)", request.Username, ipLong, preEventSubQuery).Find(&previousEvent)

	if previousEvent.Ip > 0 {
		distance := Distance(previousEvent.Latitude, previousEvent.Longitude, cityBlock.Latitude, cityBlock.Longitude, previousEvent.Radius, cityBlock.AccuracyRadius)

		speed := calculateSpeed(request, previousEvent, distance)

		preAccess := IPAccess{
			Lat:       previousEvent.Latitude,
			Lon:       previousEvent.Longitude,
			Radius:    previousEvent.Radius,
			Speed:     speed,
			IP:        previousEvent.IpRaw,
			Timestamp: previousEvent.UnixTimestamp,
		}
		response.PrecedingIPAccess = preAccess
		if speed > maxSpeed {
			response.TravelToCurrentGeoSuspicious = true
		}
	}

	var postEvent EventLog
	postEventSubQuery := dbConn.Select("MIN(unix_timestamp)").Where("user_name = ? and unix_timestamp > ? and ip != ?", request.Username, request.UnixTimestamp, ipLong).Table("event_logs")
	dbConn.Where("user_name = ? AND ip != ? AND unix_timestamp = (?)", request.Username, ipLong, postEventSubQuery).Find(&postEvent)

	if postEvent.Ip > 0 {
		distance := Distance(cityBlock.Latitude, cityBlock.Longitude, postEvent.Latitude, postEvent.Longitude, cityBlock.AccuracyRadius, postEvent.Radius)

		diff := postEvent.UnixTimestamp - request.UnixTimestamp
		speed := 0.0
		if diff != 0 {
			speed = math.Abs(distance * 3600 / float64(diff))
		} else {
			speed = math.Abs(distance)
		}

		postAccess := IPAccess{
			Lat:       postEvent.Latitude,
			Lon:       postEvent.Longitude,
			Radius:    postEvent.Radius,
			Speed:     speed,
			IP:        postEvent.IpRaw,
			Timestamp: postEvent.UnixTimestamp,
		}
		response.SubsequentIPAccess = postAccess
		if speed > maxSpeed {
			response.TravelFromCurrentGeoSuspicious = true
		}
	}

	dbConn.Create(&EventLog{
		UserName:      request.Username,
		Ip:            ipLong,
		IpRaw:         request.IpAddress,
		Uuid:          request.EventUuid,
		UnixTimestamp: request.UnixTimestamp,
		Latitude:      cityBlock.Latitude,
		Longitude:     cityBlock.Longitude,
		Radius:        cityBlock.AccuracyRadius,
	})

	return &response, nil
}

func calculateSpeed(request SussRequest, previousEvent EventLog, distance float64) float64 {
	diff := request.UnixTimestamp - previousEvent.UnixTimestamp
	speed := 0.0
	if diff != 0 {
		speed = math.Abs(distance * 3600 / float64(diff))
	} else { // If we have 2 different in same second, I assume that speed is equal for distance.
		speed = math.Abs(distance)
	}
	return speed
}

func Ip2long(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
