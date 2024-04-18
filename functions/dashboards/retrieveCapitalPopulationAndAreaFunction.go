package dashboards

import (
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// TestUrlRetrieveCapitalPopulationAndArea variable used when testing the RetrieveCapitalPopulationAndArea function
var TestUrlRetrieveCapitalPopulationAndArea string

// RetrieveCapitalPopulationAndArea Retrieves the capital, population and area of a country to be inserted in a dashboard
func RetrieveCapitalPopulationAndArea(isoCode string, id int, runTest bool) (resources.CapitalPopulationArea, error) {
	// Variable used in error message for HttpRequest function.
	fetching := "capital, population and area"

	// Construct URL
	var urlPath = fmt.Sprintf(resources.RestCountriesPath+"/alpha/%s", isoCode)
	url := ConstructUrlForApiOrTest(urlPath, TestUrlRetrieveCapitalPopulationAndArea, runTest)

	// Make HTTP request to specified URL
	response, err := HttpRequest(url, fetching, id)
	// Defer close of response body
	defer CloseResponseBody(response.Body, fetching, id)

	// Check status code of response
	if response.StatusCode != http.StatusOK {
		return resources.CapitalPopulationArea{}, fmt.Errorf("HTTP error: %s", response.Status)
	}

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
