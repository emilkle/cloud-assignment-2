package resources

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Features struct {
	Temperature      float64     `json:"temperature"`
	Precipitation    bool        `json:"precipitation"`
	Capital          bool        `json:"capital"`
	Coordinates      Coordinates `json:"coordinates"`
	Population       bool        `json:"population"`
	Area             bool        `json:"area"`
	TargetCurrencies []string    `json:"targetCurrencies"`
}
