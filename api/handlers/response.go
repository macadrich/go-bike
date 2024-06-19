package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/macadrich/go-bike/database/models"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type StationsResponse struct {
	At       string            `json:"at"`
	Stations []models.Stations `json:"stations,omitempty"`
	Weather  models.WeatherMap `json:"weather,omitempty"`
}

func sendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
