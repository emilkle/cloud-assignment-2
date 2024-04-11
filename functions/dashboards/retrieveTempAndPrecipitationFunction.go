package dashboards

import (
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
)

// RetrieveTempAndPrecipitation Retrieves 24 hour temperature and precipitation values at specified coordinates
func RetrieveTempAndPrecipitation(latitude, longitude float64, id int) (resources.HourlyData, error) {
	fetching := "temp and precipitation"

	// Construct URL
	url := fmt.Sprintf(resources.METEO_TEMP_PERCIP+"/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation&forecast_days=1", latitude, longitude)

	// Make HTTP request to specified URL
	response, err := HttpRequest(url, fetching, id)
	// Defer close of response body
	defer CloseResponseBody(response.Body, fetching, id)

	// Decode JSON response
	var forecastResponse resources.ForecastResponse
	err = json.NewDecoder(response.Body).Decode(&forecastResponse)
	if err != nil {
		return resources.HourlyData{}, fmt.Errorf("failed to decode JSON response: %s", err)
	}

	// DEBUGGING
	log.Printf("Decoded API Response: %+v", forecastResponse)

	// Check if any values were returned
	if len(forecastResponse.Hourly.Temperature) == 0 &&
		len(forecastResponse.Hourly.Precipitation) == 0 {
		return resources.HourlyData{}, fmt.Errorf("no temperature and precipitation data returned")
	}

	// Create and store temperature and precipitation data in struct
	tempAndPrecipitationData := forecastResponse.Hourly

	// Log and make sure if any temp and precipitation data was retrieved from the response
	log.Printf("Retrieved temp and precipitation: %+v", tempAndPrecipitationData)

	return tempAndPrecipitationData, nil
}