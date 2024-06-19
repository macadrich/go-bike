package client

import (
	"encoding/json"
	"log"

	"github.com/macadrich/go-bike/database/models"
)

const (
	Latest   = "last_updated"
	Stations = "features"
)

type ClientResponse struct {
	data map[string]interface{}
}

type stationsResponse struct {
	Geometry   models.Geometry `json:"geometry"`
	Properties models.Stations `json:"properties"`
	Type       string          `json:"type"`
}

func (resp *ClientResponse) HasInValidData() bool {
	if resp.data != nil {
		return false
	}
	return len(resp.LastUpdated()) == 0 || len(resp.Stations()) == 0
}

func (resp *ClientResponse) LastUpdated() string {
	if resp.data != nil {
		lastUpdated, ok := resp.data[Latest].(string)
		if !ok {
			log.Println("type assertion failed:", Latest)
			return ""
		}
		return lastUpdated
	}
	return ""
}

func (resp *ClientResponse) Stations() []stationsResponse {
	if resp.data != nil {
		features, ok := resp.data[Stations].([]any)
		if !ok {
			log.Println("Type assertion failed:", Stations)
			return nil
		}

		var listOfStations []stationsResponse
		jsonData, err := json.Marshal(features)
		if err != nil {
			log.Println("Error marshaling features:", err)
			return nil
		}

		err = json.Unmarshal(jsonData, &listOfStations)
		if err != nil {
			log.Println("Error unmarshaling features:", err)
			return nil
		}

		return listOfStations
	}

	return nil
}

func (resp *ClientResponse) WeatherUpdate() *models.WeatherMap {
	if resp.data != nil {
		weatherData, err := json.Marshal(resp.data)
		if err != nil {
			return nil
		}

		var latestUpdateWeather models.WeatherMap
		err = json.Unmarshal(weatherData, &latestUpdateWeather)
		if err != nil {
			return nil
		}

		return &latestUpdateWeather
	}

	return nil
}
