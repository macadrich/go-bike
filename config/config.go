package config

import (
	"log"

	"github.com/spf13/viper"
)

type ThirdpartyAPI struct {
	APIKey     string
	BikeURL    string
	WeatherURL string
	City       string
}

type DBConfig struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	ThirdpartyAPI ThirdpartyAPI
}

func Config() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("./config/config.yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	return v
}

func LoadDBConfig() *DBConfig {
	v := Config()
	return &DBConfig{
		DBHost:     v.GetString("Database.Host"),
		DBPort:     v.GetString("Database.Port"),
		DBUser:     v.GetString("Database.User"),
		DBPassword: v.GetString("Database.Password"),
		DBName:     v.GetString("Database.Name"),
		ThirdpartyAPI: ThirdpartyAPI{
			APIKey:     v.GetString("ThirdpartAPI.APIKey"),
			BikeURL:    v.GetString("ThirdpartAPI.BicycleTrancit"),
			WeatherURL: v.GetString("ThirdpartAPI.CurrentWeather"),
			City:       v.GetString("ThirdpartAPI.City"),
		},
	}
}

func LoadAuthorization() string {
	v := Config()
	return v.GetString("Authorization.Token")
}
