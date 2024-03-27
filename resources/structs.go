package resources

// Regitrations structs

type RegistrationsGET struct {
	Id         int      `json:"id"`
	Country    string   `json:"country"`
	IsoCode    string   `json:"isoCode"`
	Features   Features `json:"features"`
	LastChange string   `json:"lastChange"`
}

type RegistrationsPOSTandPUT struct {
	Country  string   `json:"country"`
	IsoCode  string   `json:"isoCode"`
	Features Features `json:"features"`
}

type RegistrationsPOSTResponse struct {
	Id         int    `json:"id"`
	LastChange string `json:"lastChange"`
}

type Features struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

// StatusResponse struct for the status response body
type StatusResponse struct {
	CountriesApi int `json:"countries_api"`
	MeteoApi     int `json:"meteo_api"`
	CurrencyApi  int `json:"currency_api"`
	//add notification_db and webhook
	Version string  `json:"version"`
	Uptime  float64 `json:"uptime"`
}

// DashboardsGet Struct to display a dashboard and the last time it was retrieved
type DashboardsGet struct {
	Country       string   `json:"country"`
	IsoCode       string   `json:"isoCode"`
	Features      Features `json:"features"`
	LastRetrieval string   `json:"last_retrieval"`
}
