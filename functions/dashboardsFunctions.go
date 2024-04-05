package functions

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
// RetrieveDashboardGet returns a single/specific dashboard based on the dashboard ID.
func RetrieveDashboardGet(dashboardId string) (resources.DashboardsGet, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Convert/parse string to integer
	idNumber, err := strconv.Atoi(dashboardId)
	if err != nil {
		log.Printf("Failed to parse ID: %s. Error: %s", dashboardId, err)
		return resources.DashboardsGet{}, err
	}

	// Make query to the database to return all documents based on the specified ID
	query := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==", idNumber).Limit(1)
	documents, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to fetch documents. Error: %s", err)
		return resources.DashboardsGet{}, err
	}

	// Check if any document with the specified ID were found
	if len(documents) == 0 {
		err := fmt.Errorf("no document found with ID: %s", dashboardId)
		log.Println(err)
		return resources.DashboardsGet{}, err
	}

	// Create a timestamp for the last time this dashboard was retrieved
	var lastRetrieved = time.Now().Format("20060102 15:04")

	// Take only the first document returned by the query
	data := documents[0].Data()
	featuresData := data["features"].(map[string]interface{})

	// Returns a dashboard
	return resources.DashboardsGet{
		Country: data["country"].(string),
		IsoCode: data["isoCode"].(string),
		Features: resources.Features{
			Temperature:      featuresData["temperature"].(bool),
			Precipitation:    featuresData["precipitation"].(bool),
			Capital:          featuresData["capital"].(bool),
			Coordinates:      featuresData["coordinates"].(bool),
			Population:       featuresData["population"].(bool),
			Area:             featuresData["area"].(bool),
			TargetCurrencies: registrations.GetTargetCurrencies(featuresData),
		},
		LastRetrieval: lastRetrieved,
	}, nil
}
*/

// RetrieveDashboardGet returns a single/specific dashboard based on the dashboard ID.
func RetrieveDashboardGet(dashboardId string) (resources.DashboardsGetTest, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Convert/parse string to integer
	idNumber, err := strconv.Atoi(dashboardId)
	if err != nil {
		log.Printf("Failed to parse ID: %s. Error: %s", dashboardId, err)
		return resources.DashboardsGetTest{}, err
	}

	// Make query to the database to return all documents based on the specified ID
	query := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==", idNumber).Limit(1)
	documents, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to fetch documents. Error: %s", err)
		return resources.DashboardsGetTest{}, err
	}

	// Check if any document with the specified ID were found
	if len(documents) == 0 {
		err := fmt.Errorf("no document found with ID: %s", dashboardId)
		log.Println(err)
		return resources.DashboardsGetTest{}, err
	}

	// Create a timestamp for the last time this dashboard was retrieved
	var lastRetrieved = time.Now().Format("20060102 15:04")

	// Take only the first document returned by the query
	data := documents[0].Data()
	featuresData := data["features"].(map[string]interface{})

	// Create variable for coordinates
	var coordinates resources.CoordinatesValues
	var capitalPopArea resources.CapitalPopulationArea

	// Checks if coordinates belong in this dashboard configuration
	if featuresData["coordinates"].(bool) {
		coordinates, err = RetrieveCoordinates(data["country"].(string), idNumber)
	}

	// Retrieve capital, population and area
	capitalPopArea, err = RetrieveCapitalPopulationAndArea(data["isoCode"].(string), idNumber)

	var capital string
	var population int
	var area float64

	// Check if dashboard configuration wants capital, population or area
	if featuresData["capital"].(bool) {
		capital = capitalPopArea.Capital[0]
	}
	if featuresData["population"].(bool) {
		population = capitalPopArea.Population
	}
	if featuresData["area"].(bool) {
		area = capitalPopArea.Area
	}

	// Returns dashboard populated with values depending on the configuration
	return resources.DashboardsGetTest{
		Country: data["country"].(string),
		IsoCode: data["isoCode"].(string),
		FeatureValues: resources.FeatureValues{
			//Temperature:      featuresData["temperature"].(bool),
			//Precipitation:    featuresData["precipitation"].(bool),
			Capital:     capital,
			Coordinates: coordinates,
			Population:  population,
			Area:        area,
			//TargetCurrencies: registrations.GetTargetCurrencies(featuresData),
		},
		LastRetrieval: lastRetrieved,
	}, nil
}

func RetrieveMeanTemperature() {

}

func RetrieveMeanPercipitation() {

}

func RetrieveCoordinates(country string, id int) (resources.CoordinatesValues, error) {
	// Construct URL
	url := fmt.Sprintf(resources.GEOCODING_METEO+"/search?name=%s&count=1&language=en&format=json", country)
	//TEST
	//url := "https://geocoding-api.open-meteo.com/v1/search?name=norway&count=1&language=en&format=json"

	// Make HTTP request
	response, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch coordinates for dashboard with id: %d. Error: %s", id, err)
		return resources.CoordinatesValues{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body when fetching coordinates for dashboard with id: %d. Error: %s", id, err)
		}
	}(response.Body)

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
	log.Printf("Latitude: %f, Longitude: %f", latitude, longitude)

	// Create coordinate struct
	coordinates := resources.CoordinatesValues{
		Latitude:  latitude,
		Longitude: longitude,
	}

	// Log and make sure coordinates are retrieved from the response
	log.Printf("Retrieved coordinates: %+v", coordinates)

	// Return data
	return coordinates, nil
}

// RetrieveCapitalPopulationAndArea Retrieves the capital, population and area
// of a country to be inserted in a dashboard
func RetrieveCapitalPopulationAndArea(isoCode string, id int) (resources.CapitalPopulationArea, error) {
	// Construct URL
	url := fmt.Sprintf(resources.REST_COUNTRIES_PATH+"/alpha/%s", isoCode)
	//url := "http://129.241.150.113:8080/v3.1/alpha" + isoCode

	// Make HTTP request
	response, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch capital, population and area for dashboard with id: %d. Error: %s", id, err)
		return resources.CapitalPopulationArea{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body when fetching capital, population and area for dashboard with id: %d. Error: %s", id, err)
		}
	}(response.Body)

	var data []resources.CapitalPopulationArea
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return resources.CapitalPopulationArea{}, fmt.Errorf("failed to decode JSON response: %s", err)
	}

	// Check if data is not empty
	if len(data) == 0 {
		return resources.CapitalPopulationArea{}, fmt.Errorf("no data found for ISO code: %s", isoCode)
	}

	//TEST
	log.Printf("Retrieved capital, population, and area data: %+v", data[0])

	return data[0], nil
}

func RetrieveCurrencies() {

}
