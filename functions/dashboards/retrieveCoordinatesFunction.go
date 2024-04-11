package dashboards

import (
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
)

// RetrieveCoordinates Retrieves the country coordinates for a dashboard
func RetrieveCoordinates(country string, id int) (resources.CoordinatesValues, error) {
	// Variable used in error message for HttpRequest function.
	fetching := "coordinates"

	// Construct URL
	url := fmt.Sprintf(resources.GEOCODING_METEO+"/search?name=%s&count=1&language=en&format=json", country)

	// Make HTTP request to specified URL
	response, err := HttpRequest(url, fetching, id)
	// Defer close of response body
	defer CloseResponseBody(response.Body, fetching, id)

	// Decode the JSON response
	var coordinatesResponse resources.CoordinatesResponse
	err = json.NewDecoder(response.Body).Decode(&coordinatesResponse)
	if err != nil {
		return resources.CoordinatesValues{}, fmt.Errorf("failed to decode JSON response: %s", err)
	}

	// Check if there are any results
	if len(coordinatesResponse.Results) == 0 {
		return resources.CoordinatesValues{}, fmt.Errorf("no coordinates found for dashboard: %d", id)
	}

	// Extract latitude and longitude from json response
	latitude := coordinatesResponse.Results[0].Latitude
	longitude := coordinatesResponse.Results[0].Longitude

	// Create and store coordinates in coordinates struct
	coordinates := resources.CoordinatesValues{
		Latitude:  latitude,
		Longitude: longitude,
	}

	// Log and make sure coordinates are retrieved from the response
	log.Printf("Retrieved coordinates: %+v", coordinates)

	// Return data
	return coordinates, nil
}