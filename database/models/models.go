package models

type Snapshots struct {
	At       string     `json:"at"`
	Stations Stations   `json:"stations"`
	Weather  WeatherMap `json:"weather"`
}

type WeatherMap struct {
	Id         int64      `json:"id"`
	Dt         int        `json:"dt"`
	TimeZone   int        `json:"timezone"`
	Visibility int        `json:"visibility"`
	Base       string     `json:"base"`
	Rain       Rain       `json:"rain,omitempty"`
	Name       string     `json:"name"`
	COD        int        `json:"cod"`
	Clouds     Clouds     `json:"clouds"`
	Coordinate Coordinate `json:"coord"`
	Sys        Sys        `json:"sys"`
	Main       Main       `json:"main"`
	Wind       Wind       `json:"wind"`
	Weather    []Weather  `json:"weather"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type Rain struct {
	Hour float32
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	Sunrise int32  `json:"sunrise"`
	Sunset  int32  `json:"sunset"`
	Id      int64  `json:"id"`
	Country string `json:"country"`
}

type Wind struct {
	Deg   int     `json:"deg"`
	Speed float32 `json:"speed"`
}

type Main struct {
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
	Temp      float32 `json:"temp"`
	FeelsLike float32 `json:"feels_like"`
	TempMin   float32 `json:"temp_min"`
	TempMax   float32 `json:"temp_max"`
}

type CurrentWeather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Stations struct {
	Id                     int       `json:"id"`
	At                     string    `json:"at"`
	IsEventBased           bool      `json:"isEventBased"`
	IsVirtual              bool      `json:"isVirtual"`
	KioskId                int       `json:"kioskId"`
	TrikesAvailable        int       `json:"trikesAvailable"`
	TotalDocks             int       `json:"totalDocks"`
	DocksAvailable         int       `json:"docksAvailable"`
	BikesAvailable         int       `json:"bikesAvailable"`
	ClassicBikesAvailable  int       `json:"classicBikesAvailable"`
	SmartBikesAvailable    int       `json:"smartBikesAvailable"`
	ElectricBikesAvailable int       `json:"electricBikesAvailable"`
	RewardBikesAvailable   int       `json:"rewardBikesAvailable"`
	RewardDocksAvailable   int       `json:"rewardDocksAvailable"`
	KioskType              int       `json:"kioskType"`
	Latitude               float64   `json:"latitude"`
	Longitude              float64   `json:"longitude"`
	Name                   string    `json:"name"`
	KioskStatus            string    `json:"kiokStatus"`
	KioskPublicStatus      string    `json:"kioskPublicStatus"`
	KioskConnectionStatus  string    `json:"kioskConnectionStatus"`
	AddressStreet          string    `json:"addressStreet"`
	AddressCity            string    `json:"addressCity"`
	AddressState           string    `json:"addressState"`
	AddressZipCode         string    `json:"addressZipCode"`
	CloseTime              string    `json:"closeTime"`
	EventEnd               string    `json:"eventEnd"`
	EventStart             string    `json:"eventStart"`
	Notes                  string    `json:"notes"`
	OpenTime               string    `json:"openTime"`
	PublicText             string    `json:"publicText"`
	TimeZone               string    `json:"timeZone"`
	Coordinates            []float64 `json:"coordinates"`
	Bikes                  []Bike    `json:"bikes"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Bike struct {
	Id          int  `json:"id"`
	KioskId     int  `json:"kioskId"`
	DockNumber  int  `json:"dockNumber"`
	IsElectric  bool `json:"isElectric"`
	IsAvailable bool `json:"isAvailable"`
	Battery     int  `json:"battery"`
}

// weather coordinates
type Coordinate struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}

// responses
type StationsResponse struct {
	At       string     `json:"at"`
	Stations []Stations `json:"stations"`
	Weather  WeatherMap `json:"weather"`
}

type StationResponse struct {
	At      string     `json:"at"`
	Station Stations   `json:"stations"`
	Weather WeatherMap `json:"weather"`
}
