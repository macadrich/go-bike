package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/macadrich/go-bike/config"
	"github.com/macadrich/go-bike/database"
	"github.com/macadrich/go-bike/database/models"
)

type postgresDB struct {
	db *sql.DB
}

func NewDB(cfg *config.DBConfig) (database.Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	return initDB("postgres", dsn, 25, 5)
}

func initDB(dialect, dsn string, idleConn, maxConn int) (database.Database, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return &postgresDB{db}, nil
}

func (p *postgresDB) InsertStation(lastUpdated string, station *models.Stations) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx, lastUpdated, station.Name, station.KioskId, station.TotalDocks,
		station.IsEventBased, station.IsVirtual, station.TrikesAvailable,
		station.DocksAvailable, station.BikesAvailable, station.ClassicBikesAvailable,
		station.SmartBikesAvailable, station.ElectricBikesAvailable, station.RewardBikesAvailable,
		station.RewardDocksAvailable, station.KioskType, station.Latitude, station.Longitude,
		station.KioskStatus, station.KioskPublicStatus, station.KioskConnectionStatus, station.AddressStreet,
		station.AddressCity, station.AddressState, station.AddressZipCode, station.CloseTime,
		station.EventEnd, station.EventStart, station.Notes, station.OpenTime,
		station.PublicText, station.TimeZone,
	)

	if err != nil {
		return err
	}

	rowAffected, rowErr := result.RowsAffected()
	if rowErr != nil {
		return fmt.Errorf("error on row affected: %w", rowErr)
	}

	if rowAffected > 0 {
		if len(station.Bikes) > 0 {
			err := p.insertBikes(ctx, station.KioskId, lastUpdated, station.Bikes)
			if err != nil {
				return err
			}
		}

		if len(station.Coordinates) > 0 {
			err := p.insertCoordinates(ctx, station.KioskId, station.Coordinates)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *postgresDB) QueryAllStation(lastUpdate string) ([]models.Stations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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

	rows, err := p.db.QueryContext(ctx, query, lastUpdate)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var stations []models.Stations
	for rows.Next() {
		var station models.Stations

		err := rows.Scan(
			&station.At, &station.Name, &station.KioskId, &station.TotalDocks, &station.IsEventBased,
			&station.IsVirtual, &station.TrikesAvailable, &station.DocksAvailable,
			&station.BikesAvailable, &station.ClassicBikesAvailable, &station.SmartBikesAvailable,
			&station.ElectricBikesAvailable, &station.RewardBikesAvailable, &station.RewardDocksAvailable,
			&station.KioskType, &station.Latitude, &station.Longitude, &station.KioskStatus,
			&station.KioskPublicStatus, &station.KioskConnectionStatus, &station.AddressStreet, &station.AddressCity,
			&station.AddressState, &station.AddressZipCode, &station.CloseTime, &station.EventEnd,
			&station.EventStart, &station.Notes, &station.OpenTime, &station.PublicText, &station.TimeZone,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		stations = append(stations, station)
	}

	for i, s := range stations {
		bikes, err := p.fetchBikes(ctx, s.KioskId)
		if err != nil {
			return nil, fmt.Errorf("fetch error: %w", err)
		}
		stations[i].Bikes = bikes
	}

	for i, s := range stations {
		coordinates, err := p.fetchCoordinates(ctx, s.KioskId)
		if err != nil {
			return nil, fmt.Errorf("fetch error: %w", err)
		}
		stations[i].Coordinates = coordinates
	}

	return stations, nil
}

func (p *postgresDB) QuerySpecificStation(kioskId int, lastUpdate string) (*models.Stations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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

	rows, err := p.db.QueryContext(ctx, query, kioskId, lastUpdate)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var station models.Stations
	for rows.Next() {
		err := rows.Scan(
			&station.At, &station.Name, &station.KioskId, &station.TotalDocks, &station.IsEventBased,
			&station.IsVirtual, &station.TrikesAvailable, &station.DocksAvailable,
			&station.BikesAvailable, &station.ClassicBikesAvailable, &station.SmartBikesAvailable,
			&station.ElectricBikesAvailable, &station.RewardBikesAvailable, &station.RewardDocksAvailable,
			&station.KioskType, &station.Latitude, &station.Longitude, &station.KioskStatus,
			&station.KioskPublicStatus, &station.KioskConnectionStatus, &station.AddressStreet, &station.AddressCity,
			&station.AddressState, &station.AddressZipCode, &station.CloseTime, &station.EventEnd,
			&station.EventStart, &station.Notes, &station.OpenTime, &station.PublicText, &station.TimeZone,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
	}

	bikes, err := p.fetchBikes(ctx, station.KioskId)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}
	station.Bikes = bikes

	coordinates, err := p.fetchCoordinates(ctx, station.KioskId)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}
	station.Coordinates = coordinates

	return &station, nil
}

func (p *postgresDB) insertCoordinates(ctx context.Context, kioskId int, coordinates []float64) error {
	query := "INSERT INTO coordinates (kiosk_id, longitude, latitude) VALUES($1,$2,$3)"

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, kioskId, coordinates[0], coordinates[1])
	if err != nil {
		return fmt.Errorf("error inserting data into database: %w", err)
	}

	return nil
}

func (p *postgresDB) fetchCoordinates(ctx context.Context, kioskId int) ([]float64, error) {
	query := "SELECT longitude, latitude FROM coordinates WHERE kiosk_id = $1"
	rows, err := p.db.QueryContext(ctx, query, kioskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coordinates []float64
	for rows.Next() {
		var longitude float64
		var latitude float64
		err := rows.Scan(&longitude, &latitude)
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, []float64{longitude, latitude}...)
	}

	return coordinates, nil
}

func (p *postgresDB) insertBikes(ctx context.Context, kioskId int, lastUpdated string, bikes []models.Bike) error {
	query := `
		INSERT INTO bikes 
		(at, kiosk_id, dock_number, is_electric, is_available, battery)
		VALUES($1,$2,$3,$4,$5,$6)
	`

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, bike := range bikes {
		_, err := stmt.ExecContext(ctx, lastUpdated, kioskId, bike.DockNumber, bike.IsElectric, bike.IsAvailable, bike.Battery)
		if err != nil {
			return fmt.Errorf("error inserting data into database: %w", err)
		}
	}

	return nil
}

func (p *postgresDB) fetchBikes(ctx context.Context, kioskId int) ([]models.Bike, error) {
	query := "SELECT id, kiosk_id, dock_number, is_electric, is_available, battery FROM bikes WHERE kiosk_id = $1"
	rows, err := p.db.QueryContext(ctx, query, kioskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bikes []models.Bike
	for rows.Next() {
		var bike models.Bike
		err := rows.Scan(&bike.Id, &bike.KioskId, &bike.DockNumber, &bike.IsElectric, &bike.IsAvailable, &bike.Battery)
		if err != nil {
			return nil, err
		}
		bikes = append(bikes, bike)
	}

	return bikes, nil
}
