package dashboards

import (
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
)

// RetrieveCapitalPopulationAndArea Retrieves the capital, population and area of a country to be inserted in a dashboard
func RetrieveCapitalPopulationAndArea(isoCode string, id int) (resources.CapitalPopulationArea, error) {
	// Variable used in error message for HttpRequest function.
	fetching := "capital, population and area"

	// Construct URL
	url := fmt.Sprintf(resources.REST_COUNTRIES_PATH+"/alpha/%s", isoCode)

	// Make HTTP request to specified URL
	response, err := HttpRequest(url, fetching, id)
	// Defer close of response body
	defer CloseResponseBody(response.Body, fetching, id)

	// Decode the JSON response
	var capPopArea []resources.CapitalPopulationArea
	err = json.NewDecoder(response.Body).Decode(&capPopArea)
	if err != nil {
		return resources.CapitalPopulationArea{}, fmt.Errorf("failed to decode JSON response: %s", err)
	}

	// Check if data has any results
	if len(capPopArea) == 0 {
		return resources.CapitalPopulationArea{}, fmt.Errorf("no data found for ISO code: %s", isoCode)
	}

	// Log and make sure data was returned
	log.Printf("Retrieved capital, population, and area data: %+v", capPopArea[0])

	return capPopArea[0], nil
}