package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/macadrich/go-bike/database/models"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
	mockData map[string]any
}

func NewMockDB() *MockDB {
	return &MockDB{
		mockData: make(map[string]any),
	}
}

func (m *MockDB) InsertStation(ctx context.Context) error {
	args := m.Called(m.mockData["lastUpdated"])
	return args.Error(0)
}

func (m *MockDB) QueryAllStation(ctx context.Context, lastUpdate string) (*models.StationsResponse, error) {

	return nil, nil
}

func (m *MockDB) QuerySpecificStation(kioskId int, lastUpdate string) (*models.Stations, error) {

	return nil, nil
}

func TestInsertStation(t *testing.T) {
	mockDB := NewMockDB()

	mockDB.mockData["lastUpdated"] = []models.Stations{
		{
			Id: 3005,
			At: "2024-05-14T06:48:19.588Z",
		},
	}

	mockDB.On("InsertStation", mockDB.mockData["lastUpdated"]).Return(nil)

	handlers := NewHandlers(mockDB)
	req, err := http.NewRequest("GET", "/api/v1/indego-data-fetch-and-store-it-db", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.InsertStation)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestQueryAllStations(t *testing.T) {
	mockDB := NewMockDB()
	mockDB.mockData["lastUpdated"] = "2024-05-14T06:48:19.588Z"
	mockDB.On("QueryAllStation", mock.Anything, mockDB.mockData["lastUpdated"]).Return(&models.StationsResponse{}, nil)
	handlers := NewHandlers(mockDB)
	req, err := http.NewRequest("GET", "/api/v1/indego-data-fetch-and-store-it-db", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.QueryAllStation)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
