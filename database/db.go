package database

import "github.com/macadrich/go-bike/database/models"

type Database interface {
	InsertStation(lastUpdated string, station *models.Stations) error
	QueryAllStation(lastUpdate string) ([]models.Stations, error)
	QuerySpecificStation(kioskId int, lastUpdate string) (*models.Stations, error)
}
