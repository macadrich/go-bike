package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/macadrich/go-bike/database/models"
	"github.com/stretchr/testify/assert"
)

type MockIndegoStations struct {
	Geometry   models.Geometry `json:"geometry"`
	Properties models.Stations `json:"properties"`
	Type       string          `json:"type"`
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("database connection: '%s'", err)
	}

	return db, mock
}

func DumpFiles() []MockIndegoStations {
	var data map[string]any

	fileData, _ := os.ReadFile("../../phl.json")
	json.NewDecoder(bytes.NewBuffer(fileData)).Decode(&data)

	station := data["features"].([]any)

	jsonData, err := json.Marshal(station)
	if err != nil {
		log.Fatal("Error marshaling features:", err)
	}

	var listOfStations []MockIndegoStations
	err = json.Unmarshal(jsonData, &listOfStations)
	if err != nil {
		log.Fatal("Error unmarshaling features:", err)
	}

	return listOfStations
}

func TestInsertStation(t *testing.T) {
	db, mock := NewMock()
	postgres := &postgresDB{db}
	defer func() {
		postgres.db.Close()
	}()

	//var station = &models.Stations{}
	listOfStations := DumpFiles()
	station := listOfStations[0].Properties

	query := `
		INSERT INTO stations 
		(
			at, name, kiosk_id, total_docks, is_event_based,
			is_virtual, trikes_available, docks_available,
			bikes_available, classic_bikes_available, smart_bikes_available,
			electric_bikes_available, reward_bikes_available, reward_docks_available,
			kiosk_type, latitude, longitude, kiosk_status,
			kiosk_public_status, kiosk_connection_status, address_street,
			address_city, address_state, address_zipcode, close_time,
			event_end, event_start, notes, open_time, public_text, timezone 
		) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31)
	`

	lastUpdated := "2024-05-14T06:48:19.588Z"
	prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
	prep.ExpectExec().WithArgs(lastUpdated, station.Name, station.KioskId, station.TotalDocks,
		station.IsEventBased, station.IsVirtual, station.TrikesAvailable,
		station.DocksAvailable, station.BikesAvailable, station.ClassicBikesAvailable,
		station.SmartBikesAvailable, station.ElectricBikesAvailable, station.RewardBikesAvailable,
		station.RewardDocksAvailable, station.KioskType, station.Latitude, station.Longitude,
		station.KioskStatus, station.KioskPublicStatus, station.KioskConnectionStatus, station.AddressStreet,
		station.AddressCity, station.AddressState, station.AddressZipCode, station.CloseTime,
		station.EventEnd, station.EventStart, station.Notes, station.OpenTime,
		station.PublicText, station.TimeZone,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	kioskId := 3005
	// Expect bikes
	bikeQuery := "INSERT INTO bikes (at, kiosk_id, dock_number, is_electric, is_available, battery) VALUES($1,$2,$3,$4,$5,$6)"

	station.Bikes = station.Bikes[:1]
	bikes := station.Bikes
	bike := bikes[0]

	prep = mock.ExpectPrepare(regexp.QuoteMeta(bikeQuery))
	prep.ExpectExec().WithArgs(&lastUpdated, &kioskId, &bike.DockNumber, &bike.IsElectric, &bike.IsAvailable, &bike.Battery).WillReturnResult(sqlmock.NewResult(0, 1))

	// Expect coordinates
	coordinatesQuery := "INSERT INTO coordinates (kiosk_id, longitude, latitude) VALUES($1,$2,$3)"
	coordinates := station.Coordinates

	prep = mock.ExpectPrepare(regexp.QuoteMeta(coordinatesQuery))
	prep.ExpectExec().WithArgs(&kioskId, &coordinates[0], &coordinates[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	err := postgres.InsertStation(lastUpdated, &station)
	assert.NoError(t, err)
}

func TestQueryAllStation(t *testing.T) {
	db, mock := NewMock()
	postgres := &postgresDB{db}
	defer func() {
		postgres.db.Close()
	}()

	query := `
		SELECT at, name, kiosk_id, total_docks, is_event_based,
		is_virtual, trikes_available, docks_available,
		bikes_available, classic_bikes_available, smart_bikes_available,
		electric_bikes_available, reward_bikes_available, reward_docks_available,
		kiosk_type, latitude, longitude, kiosk_status,
		kiosk_public_status, kiosk_connection_status, address_street,
		address_city, address_state, address_zipcode, close_time,
		event_end, event_start, notes, open_time, public_text, timezone
		FROM stations WHERE at >= $1 ORDER BY at ASC
	`

	listOfStations := DumpFiles()
	station := listOfStations[0].Properties

	rows := sqlmock.NewRows([]string{"at", "name", "kiosk_id", "total_docks", "is_event_based",
		"is_virtual", "trikes_available", "docks_available",
		"bikes_available", "classic_bikes_available", "smart_bikes_available",
		"electric_bikes_available", "reward_bikes_available", "reward_docks_available",
		"kiosk_type", "latitude", "longitude", "kiosk_status",
		"kiosk_public_status", "kiosk_connection_status", "address_street",
		"address_city", "address_state", "address_zipcode", "close_time",
		"event_end", "event_start", "notes", "open_time", "public_text", "timezone"}).AddRow(station.At, station.Name, station.KioskId, station.TotalDocks,
		station.IsEventBased, station.IsVirtual, station.TrikesAvailable,
		station.DocksAvailable, station.BikesAvailable, station.ClassicBikesAvailable,
		station.SmartBikesAvailable, station.ElectricBikesAvailable, station.RewardBikesAvailable,
		station.RewardDocksAvailable, station.KioskType, station.Latitude, station.Longitude,
		station.KioskStatus, station.KioskPublicStatus, station.KioskConnectionStatus, station.AddressStreet,
		station.AddressCity, station.AddressState, station.AddressZipCode, station.CloseTime,
		station.EventEnd, station.EventStart, station.Notes, station.OpenTime,
		station.PublicText, station.TimeZone)

	lastUpdated := "2024-05-14T06:48:19.588Z"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(lastUpdated).WillReturnRows(rows)

	// Add the expected sub query for bikes
	bike := station.Bikes[0]
	bikesQuery := "SELECT id, kiosk_id, dock_number, is_electric, is_available, battery FROM bikes WHERE kiosk_id = $1"
	bikeRows := sqlmock.NewRows([]string{
		"id", "kiosk_id", "dock_number", "is_electric", "is_available", "battery",
	}).AddRow(&bike.Id, &bike.KioskId, &bike.DockNumber, &bike.IsElectric, &bike.IsAvailable, &bike.Battery)
	kioskId := 3005
	mock.ExpectQuery(regexp.QuoteMeta(bikesQuery)).WithArgs(kioskId).WillReturnRows(bikeRows)

	// Add the expected sub query for coordinates
	coordinate := station.Coordinates
	coordinateQuery := "SELECT longitude, latitude FROM coordinates WHERE kiosk_id = $1"
	coordinateRows := sqlmock.NewRows([]string{"longitude", "latitude"}).AddRow(&coordinate[0], &coordinate[1])
	mock.ExpectQuery(regexp.QuoteMeta(coordinateQuery)).WithArgs(kioskId).WillReturnRows(coordinateRows)

	stations, err := postgres.QueryAllStation(lastUpdated)
	assert.NotEmpty(t, stations)
	assert.NoError(t, err)
	assert.Len(t, stations, 1)
}

func TestQuerySpecificStation(t *testing.T) {
	db, mock := NewMock()
	postgres := &postgresDB{db}
	defer func() {
		postgres.db.Close()
	}()

	query := `
		SELECT at, name, kiosk_id, total_docks, is_event_based,
		is_virtual, trikes_available, docks_available,
		bikes_available, classic_bikes_available, smart_bikes_available,
		electric_bikes_available, reward_bikes_available, reward_docks_available,
		kiosk_type, latitude, longitude, kiosk_status,
		kiosk_public_status, kiosk_connection_status, address_street,
		address_city, address_state, address_zipcode, close_time,
		event_end, event_start, notes, open_time, public_text, timezone
		FROM stations WHERE kiosk_id = $1 AND at >= $2 ORDER BY at ASC
	`

	listOfStations := DumpFiles()
	station := listOfStations[0].Properties

	rows := sqlmock.NewRows([]string{"at", "name", "kiosk_id", "total_docks", "is_event_based",
		"is_virtual", "trikes_available", "docks_available",
		"bikes_available", "classic_bikes_available", "smart_bikes_available",
		"electric_bikes_available", "reward_bikes_available", "reward_docks_available",
		"kiosk_type", "latitude", "longitude", "kiosk_status",
		"kiosk_public_status", "kiosk_connection_status", "address_street",
		"address_city", "address_state", "address_zipcode", "close_time",
		"event_end", "event_start", "notes", "open_time", "public_text", "timezone"}).AddRow(station.At, station.Name, station.KioskId, station.TotalDocks,
		station.IsEventBased, station.IsVirtual, station.TrikesAvailable,
		station.DocksAvailable, station.BikesAvailable, station.ClassicBikesAvailable,
		station.SmartBikesAvailable, station.ElectricBikesAvailable, station.RewardBikesAvailable,
		station.RewardDocksAvailable, station.KioskType, station.Latitude, station.Longitude,
		station.KioskStatus, station.KioskPublicStatus, station.KioskConnectionStatus, station.AddressStreet,
		station.AddressCity, station.AddressState, station.AddressZipCode, station.CloseTime,
		station.EventEnd, station.EventStart, station.Notes, station.OpenTime,
		station.PublicText, station.TimeZone)

	lastUpdated := "2024-05-14T06:48:19.588Z"
	kioskId := int(3005)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(kioskId, lastUpdated).WillReturnRows(rows)

	// Add the expected sub query for bikes
	bike := station.Bikes[0]

	bikesQuery := "SELECT id, kiosk_id, dock_number, is_electric, is_available, battery FROM bikes WHERE kiosk_id = $1"
	bikeColumn := []string{"id", "kiosk_id", "dock_number", "is_electric", "is_available", "battery"}

	bikeRows := sqlmock.NewRows(bikeColumn).AddRow(&bike.Id, &bike.KioskId, &bike.DockNumber, &bike.IsElectric, &bike.IsAvailable, &bike.Battery)
	mock.ExpectQuery(regexp.QuoteMeta(bikesQuery)).WithArgs(kioskId).WillReturnRows(bikeRows)

	// Add the expected sub query for coordinates
	coordinate := station.Coordinates
	coordinateQuery := "SELECT longitude, latitude FROM coordinates WHERE kiosk_id = $1"
	coordinateColumn := []string{"longitude", "latitude"}

	coordinateRows := sqlmock.NewRows(coordinateColumn).AddRow(&coordinate[0], &coordinate[1])
	mock.ExpectQuery(regexp.QuoteMeta(coordinateQuery)).WithArgs(kioskId).WillReturnRows(coordinateRows)

	stations, err := postgres.QuerySpecificStation(kioskId, lastUpdated)
	assert.NotEmpty(t, stations)
	assert.NoError(t, err)
}

func TestFetchCooordinates(t *testing.T) {
	db, mock := NewMock()
	postgres := &postgresDB{db}
	defer func() {
		postgres.db.Close()
	}()

	listOfStations := DumpFiles()
	station := listOfStations[0].Properties

	kioskId := 3005

	coordinate := station.Coordinates
	coordinateQuery := "SELECT longitude, latitude FROM coordinates WHERE kiosk_id = $1"
	coordinateColumn := []string{"longitude", "latitude"}

	coordinateRows := sqlmock.NewRows(coordinateColumn).AddRow(&coordinate[0], &coordinate[1])
	mock.ExpectQuery(regexp.QuoteMeta(coordinateQuery)).WithArgs(kioskId).WillReturnRows(coordinateRows)

	bikes, err := postgres.fetchCoordinates(context.TODO(), kioskId)
	assert.NotEmpty(t, bikes)
	assert.NoError(t, err)
}

func TestInsertCoordinates(t *testing.T) {
	db, mock := NewMock()
	postgres := &postgresDB{db}
	defer func() {
		postgres.db.Close()
	}()

	listOfStations := DumpFiles()
	station := listOfStations[0].Properties

	query := "INSERT INTO coordinates (kiosk_id, longitude, latitude) VALUES($1,$2,$3)"

	coordinates := station.Coordinates
	kioskId := 3005

	prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
	prep.ExpectExec().WithArgs(&kioskId, &coordinates[0], &coordinates[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	err := postgres.insertCoordinates(context.TODO(), kioskId, coordinates)
	assert.NoError(t, err)
}

func TestFetchBikes(t *testing.T) {
	db, mock := NewMock()
	postgres := &postgresDB{db}
	defer func() {
		postgres.db.Close()
	}()

	listOfStations := DumpFiles()
	station := listOfStations[0].Properties

	bike := station.Bikes[0]
	kioskId := 3005

	bikesQuery := "SELECT id, kiosk_id, dock_number, is_electric, is_available, battery FROM bikes WHERE kiosk_id = $1"
	bikeColumn := []string{"id", "kiosk_id", "dock_number", "is_electric", "is_available", "battery"}

	bikeRows := sqlmock.NewRows(bikeColumn).AddRow(&bike.Id, &bike.KioskId, &bike.DockNumber, &bike.IsElectric, &bike.IsAvailable, &bike.Battery)
	mock.ExpectQuery(regexp.QuoteMeta(bikesQuery)).WithArgs(kioskId).WillReturnRows(bikeRows)

	bikes, err := postgres.fetchBikes(context.TODO(), kioskId)
	assert.NotEmpty(t, bikes)
	assert.NoError(t, err)
}

func TestInsertBikes(t *testing.T) {
	db, mock := NewMock()
	postgres := &postgresDB{db}
	defer func() {
		postgres.db.Close()
	}()

	listOfStations := DumpFiles()
	station := listOfStations[0].Properties

	query := "INSERT INTO bikes (at, kiosk_id, dock_number, is_electric, is_available, battery) VALUES($1,$2,$3,$4,$5,$6)"

	bikes := station.Bikes[:1] // limit 1
	bike := bikes[0]
	lastUpdated := "2024-05-14T06:48:19.588Z"
	kioskId := 3005

	prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
	prep.ExpectExec().WithArgs(&lastUpdated, &kioskId, &bike.DockNumber, &bike.IsElectric, &bike.IsAvailable, &bike.Battery).WillReturnResult(sqlmock.NewResult(0, 1))

	err := postgres.insertBikes(context.TODO(), kioskId, lastUpdated, bikes)
	assert.NoError(t, err)
}
