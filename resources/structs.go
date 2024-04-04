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

// Notification endpoint structs

// WebhookPOST struct for POST request
type WebhookPOST struct {
	URL     string `json:"url"`
	IsoCode string `json:"country"`
	Event   string `json:"event"`
}

// WebhookPOSTResponse struct for POST response
type WebhookPOSTResponse struct {
	ID         string `json:"ID"`
	LastChange string `json:"LastChange"`
}

// ViewWebhook struct for view specific webhook GET response
type ViewWebhook struct {
	ID      string `json:"ID"`
	URL     string `json:"URL"`
	Country string `json:"Country"`
	Event   string `json:"Event"`
}

// EventData struct for webhook invocation
type EventData struct {
	ID      string `json:"ID"`
	Country string `json:"Country"`
	Event   string `json:"Event"`
	Time    string `json:"Time"`
}
