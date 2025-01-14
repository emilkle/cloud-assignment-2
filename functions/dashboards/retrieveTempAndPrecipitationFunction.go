package dashboards

import (
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// TestUrlRetrieveTempAndPrecipitation variable used when testing the RetrieveTempAndPrecipitation function
var TestUrlRetrieveTempAndPrecipitation string

// RetrieveTempAndPrecipitation Retrieves 24 hour temperature and precipitation values at specified coordinates
func RetrieveTempAndPrecipitation(latitude, longitude float64, id int, runTest bool) (resources.HourlyData, error) {
	fetching := "temp and precipitation"

	// Construct URL
	var urlPath = fmt.Sprintf(resources.MeteoTempPercip+"/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation&forecast_days=1", latitude, longitude)
	url := ConstructUrlForApiOrTest(urlPath, TestUrlRetrieveTempAndPrecipitation, runTest)

	// Make HTTP request to specified URL
	response, err := HttpRequest(url, fetching, id)
	// Defer close of response body
	defer CloseResponseBody(response.Body, fetching, id)

	// Check status code of response
	if response.StatusCode != http.StatusOK {
		return resources.HourlyData{}, fmt.Errorf("HTTP error: %s", response.Status)
	}

	// Decode JSON response
	var forecastResponse resources.ForecastResponse
	err = json.NewDecoder(response.Body).Decode(&forecastResponse)
	if err != nil {
		return resources.HourlyData{}, fmt.Errorf("failed to decode JSON response: %s", err)
	}

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
