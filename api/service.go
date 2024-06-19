package api

import (
	"context"
	"fmt"

	"github.com/macadrich/go-bike/client"
	"github.com/macadrich/go-bike/config"
	"github.com/macadrich/go-bike/database"
	"github.com/macadrich/go-bike/database/models"
)

type IService interface {
	InsertStation(ctx context.Context) error
	QueryAllStation(ctx context.Context, lastUpdate string) (*models.StationsResponse, error)
	QuerySpecificStation(kioskId int, lastUpdate string) (*models.Stations, error)
}

type service struct {
	db     database.Database
	client client.IClient
	cfg    *config.DBConfig
}

func NewService(db database.Database, client client.IClient) IService {
	return &service{db, client, config.LoadDBConfig()}
}

func (s *service) InsertStation(ctx context.Context) error {
	resp, err := s.client.GetData(ctx, s.cfg.ThirdpartyAPI.BikeURL)
	if err != nil || resp.HasInValidData() {
		return err
	}

	for _, v := range resp.Stations() {
		if err := s.db.InsertStation(resp.LastUpdated(), &v.Properties); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) QueryAllStation(ctx context.Context, lastUpdate string) (*models.StationsResponse, error) {
	listOfStations, err := s.db.QueryAllStation(lastUpdate)
	if err != nil {
		return nil, err
	}

	weatherURL := fmt.Sprintf("%s?q=%s&appid=%s&units=imperial", s.cfg.ThirdpartyAPI.WeatherURL, s.cfg.ThirdpartyAPI.City, s.cfg.ThirdpartyAPI.APIKey)
	result, err := s.client.GetData(ctx, weatherURL)
	if err != nil {
		return nil, err
	}

	response := models.StationsResponse{
		At:       lastUpdate,
		Stations: listOfStations,
		Weather:  *result.WeatherUpdate(),
	}

	return &response, nil
}

func (s *service) QuerySpecificStation(kioskId int, lastUpdate string) (*models.Stations, error) {
	station, err := s.db.QuerySpecificStation(kioskId, lastUpdate)
	if err != nil {
		return nil, err
	}

	return station, nil
}
