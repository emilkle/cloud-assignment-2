package resources

// Registrations endpoint structs

// RegistrationsGET struct for the HTTP GET request's response body.
type RegistrationsGET struct {
	Id         int      `json:"id"`
	Country    string   `json:"country"`
	IsoCode    string   `json:"isoCode"`
	Features   Features `json:"features"`
	LastChange string   `json:"lastChange"`
}

// RegistrationsPOSTandPUT struct for HTTP POST and PUT requests.
type RegistrationsPOSTandPUT struct {
	Country  string   `json:"country"`
	IsoCode  string   `json:"isoCode"`
	Features Features `json:"features"`
}

// RegistrationsPOSTResponse struct for the HTTP POST request's response body.
type RegistrationsPOSTResponse struct {
	Id         int    `json:"id"`
	LastChange string `json:"lastChange"`
}

// Features struct for the features contained in each registration document.
type Features struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

// Status endpoint structs

// StatusResponse struct for the status response body
type StatusResponse struct {
	CountriesApi   int     `json:"countries_api"`
	MeteoApi       int     `json:"meteo_api"`
	CurrencyApi    int     `json:"currency_api"`
	NotificationDB int     `json:"notification_db"`
	Webhooks       int     `json:"webhooks"`
	Version        string  `json:"version"`
	Uptime         float64 `json:"uptime"`
}

//###############################################################################

// DashboardsGet struct for the dashboard response
type DashboardsGet struct {
	Country       string        `json:"country"`
	IsoCode       string        `json:"isoCode"`
	FeatureValues FeatureValues `json:"features"`
	LastRetrieval string        `json:"last_retrieval"`
}

// FeatureValues struct for features in a dashboard
type FeatureValues struct {
	Temperature      float64            `json:"temperature"`
	Precipitation    float64            `json:"precipitation"`
	Capital          string             `json:"capital"`
	Coordinates      CoordinatesValues  `json:"coordinates"`
	Population       int                `json:"population"`
	Area             float64            `json:"area"`
	TargetCurrencies map[string]float64 `json:"target_currencies"`
}

// ForecastResponse struct for the forecastResponse used to fetch temperature and precipitation for a dashboard
type ForecastResponse struct {
	Hourly HourlyData `json:"hourly"`
}

// HourlyData struct for forecasted time, temperature and precipitation
type HourlyData struct {
	Time          []string  `json:"time"`
	Temperature   []float64 `json:"temperature_2m"`
	Precipitation []float64 `json:"precipitation"`
}

// CapitalPopulationArea struct for holding capital, population and area for a dashboard
type CapitalPopulationArea struct {
	Capital    []string `json:"capital"`
	Population int      `json:"population"`
	Area       float64  `json:"area"`
}

// CoordinatesResponse Struct to store the latitude and longitude from rest api
type CoordinatesResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

// CoordinatesValues struct for latitude and longitude for a dashboard
type CoordinatesValues struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TargetCurrencyValues struct for holding currency values for a dashboard
type TargetCurrencyValues struct {
	TargetCurrencies map[string]float64 `json:"target_currencies"`
}

//###########################################################################

// Notification endpoint structs

// WebhookPOSTRequest struct for POST request
type WebhookPOSTRequest struct {
	URL string `json:"url"`
}

// WebhookPOST struct for POST request
type WebhookPOST struct {
	URL     string `json:"url"`
	Country string `json:"country"`
	Event   string `json:"event"`
}

// WebhookPOSTResponse struct for POST response
type WebhookPOSTResponse struct {
	ID string `json:"ID"`
}

// Remove this struct? Or implement it as response when deleting
// WebhookDELETEResponse struct for POST response
type WebhookDELETEResponse struct {
	ID    string `json:"ID"`
	URL   string `json:"URL"`
	Event string `json:"Event"`
}

// WebhookGET struct for view specific webhook GET response
type WebhookGET struct {
	ID      string `json:"ID"`
	URL     string `json:"URL"`
	Country string `json:"Country"`
	Event   string `json:"Event"`
}

// WebhookInvocation struct for webhook invocation
type WebhookInvocation struct {
	ID      string `json:"ID"`
	Country string `json:"Country"`
	Event   string `json:"Event"`
	Time    string `json:"Time"`
}
