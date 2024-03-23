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

type Features struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}
