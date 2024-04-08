package functions

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

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

	// Variables for data in the dashboards
	var tempAndPrecip resources.HourlyData
	var coordinates resources.CoordinatesValues
	var capitalPopArea resources.CapitalPopulationArea
	var capital string
	var population int
	var area float64
	var meanTemperature float64
	var meanPrecipitation float64
	var exchangeRates resources.TargetCurrencyValues
	var selectedExchangeRates resources.TargetCurrencyValues

	// Helper variables
	var latitude float64
	var longitude float64
	var coordinateData resources.CoordinatesValues
	coordinateData, err = RetrieveCoordinates(data["country"].(string), idNumber)
	// Lat and Long used in the RetrieveTempAndPrecipitation function.
	// Because a dashboard configuration might not have the coordinates set to true
	latitude = coordinateData.Latitude
	longitude = coordinateData.Longitude

	// Checks if coordinates belong in this dashboard configuration
	if featuresData["coordinates"].(bool) {
		coordinates, err = RetrieveCoordinates(data["country"].(string), idNumber)
	}

	// Retrieve capital, population and area
	capitalPopArea, err = RetrieveCapitalPopulationAndArea(data["isoCode"].(string), idNumber)

	// Check if dashboard configuration supports capital, population or area
	if featuresData["capital"].(bool) {
		capital = capitalPopArea.Capital[0]
	}
	if featuresData["population"].(bool) {
		population = capitalPopArea.Population
	}
	if featuresData["area"].(bool) {
		area = capitalPopArea.Area
	}

	// Retrieve temperature and precipitation data
	tempAndPrecip, err = RetrieveTempAndPrecipitation(latitude, longitude, idNumber)

	//DEBUGGING
	// Print temperature and precipitation data for debugging
	fmt.Println("Temperature data:", tempAndPrecip.Temperature)
	fmt.Println("Precipitation data:", tempAndPrecip.Precipitation)

	//check if temperature is part of the dashboard config and calculate the mean
	if featuresData["temperature"].(bool) {
		sumTemperature := 0.0
		for _, temp := range tempAndPrecip.Temperature {
			if temp != 0.0 {
				sumTemperature += temp
			}
		}
		meanTemperature = sumTemperature / float64(len(tempAndPrecip.Time))
		meanTemperature = math.Round(meanTemperature*10) / 10
	}
	//check if Precipitation is part of the dashboard config and calculate the mean
	if featuresData["precipitation"].(bool) {
		sumPrecipitation := 0.0
		for _, prec := range tempAndPrecip.Precipitation {
			if prec != 0.0 {
				sumPrecipitation += prec
			}
		}
		meanPrecipitation = sumPrecipitation / float64(len(tempAndPrecip.Time))
		meanPrecipitation = math.Round(meanPrecipitation*100) / 100

	}

	//DEBUGGING
	// Print mean temperature and precipitation for debugging
	fmt.Println("Mean temperature:", meanTemperature)
	fmt.Println("Mean precipitation:", meanPrecipitation)

	//Exchange rates are always shown in a dashboard
	// Retrieve exchange rates
	exchangeRates, err = RetrieveCurrencyExchangeRates(idNumber)

	// Retrieve the target currencies as interface slice
	targetCurrenciesInterface := featuresData["targetCurrencies"].([]interface{})

	// Initialize a slice to store the string values
	targetCurrencies := make([]string, len(targetCurrenciesInterface))

	// Iterate over the interface slice and convert each element to a string
	for i, currency := range targetCurrenciesInterface {
		targetCurrencies[i] = currency.(string)
	}

	//DEBUGGING
	fmt.Println("exchangeRates:", exchangeRates)

	// Initializing the TargetCurrencies map within TargetCurrencyValues before using it
	selectedExchangeRates = resources.TargetCurrencyValues{
		TargetCurrencies: make(map[string]float64),
	}

	// Iterate over targetCurrencies slice and retrieve corresponding rates
	for _, currency := range targetCurrencies {
		rate, ok := exchangeRates.TargetCurrencies[currency]
		if ok {
			selectedExchangeRates.TargetCurrencies[currency] = rate
		}
	}
	//DEBUGGING
	fmt.Println("selectedExchangeRates:", selectedExchangeRates)

	// Returns dashboard populated with values depending on the configuration
	return resources.DashboardsGetTest{
		Country: data["country"].(string),
		IsoCode: data["isoCode"].(string),
		FeatureValues: resources.FeatureValues{
			Temperature:      meanTemperature,
			Precipitation:    meanPrecipitation,
			Capital:          capital,
			Coordinates:      coordinates,
			Population:       population,
			Area:             area,
			TargetCurrencies: selectedExchangeRates.TargetCurrencies,
		},
		LastRetrieval: lastRetrieved,
	}, nil
}

// RetrieveTempAndPrecipitation Retrieves 24 hour temperature and precipitation values at specified coordinates
func RetrieveTempAndPrecipitation(latitude, longitude float64, id int) (resources.HourlyData, error) {
	// Construct URL
	url := fmt.Sprintf(resources.METEO_TEMP_PERCIP+"/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation&forecast_days=1", latitude, longitude)

	//DEBUGGING
	log.Printf("API URL: %s", url)

	// Make HTTP request
	response, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch temp and precipitation data for dashboard with id: %d. Error: %s", id, err)
		return resources.HourlyData{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body when fetching temp and precipitation data for dashboard with id: %d. Error: %s", id, err)
		}
	}(response.Body)

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

	// Log and check if any temp and precipitation data was retrieved from the response
	log.Printf("Retrieved temp and precipitation: %+v", tempAndPrecipitationData)

	return tempAndPrecipitationData, nil
}

// RetrieveCoordinates Retrieves the country coordinates for a dashboard
func RetrieveCoordinates(country string, id int) (resources.CoordinatesValues, error) {
	// Construct URL
	url := fmt.Sprintf(resources.GEOCODING_METEO+"/search?name=%s&count=1&language=en&format=json", country)

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

// RetrieveCapitalPopulationAndArea Retrieves the capital, population and area
// of a country to be inserted in a dashboard
func RetrieveCapitalPopulationAndArea(isoCode string, id int) (resources.CapitalPopulationArea, error) {
	// Construct URL
	url := fmt.Sprintf(resources.REST_COUNTRIES_PATH+"/alpha/%s", isoCode)

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

	// Decode the JSON response
	var data []resources.CapitalPopulationArea
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return resources.CapitalPopulationArea{}, fmt.Errorf("failed to decode JSON response: %s", err)
	}

	// Check if data has any results
	if len(data) == 0 {
		return resources.CapitalPopulationArea{}, fmt.Errorf("no data found for ISO code: %s", isoCode)
	}

	// Log and make sure data was returned
	log.Printf("Retrieved capital, population, and area data: %+v", data[0])

	return data[0], nil
}

// RetrieveCurrencyExchangeRates Fetches the exchange rates of currencies with NOK as base (NOK to currency)
func RetrieveCurrencyExchangeRates(id int) (resources.TargetCurrencyValues, error) {
	// Construct URL
	url := resources.CURRENCY_PATH + "NOK"

	// Make HTTP request
	response, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch exchange rates for dashboard with id: %d. Error: %s", id, err)
		return resources.TargetCurrencyValues{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body when fetching exchange rates for dashboard with id: %d. Error: %s", id, err)
		}
	}(response.Body)

	// Decode the JSON response
	var responseData map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return resources.TargetCurrencyValues{}, fmt.Errorf("failed to decode JSON response: %s", err)
	}

	// Extract rates from the response data
	ratesData, ok := responseData["rates"].(map[string]interface{})
	if !ok {
		return resources.TargetCurrencyValues{}, fmt.Errorf("failed to extract rates from JSON response")
	}

	// Populate the TargetCurrencyValues struct
	targetCurrencies := make(map[string]float64)
	for currency, rate := range ratesData {
		rateValue, ok := rate.(float64)
		if ok {
			targetCurrencies[currency] = rateValue
		}
	}

	// Create the TargetCurrencyValues struct
	exchangeRatesResponse := resources.TargetCurrencyValues{
		TargetCurrencies: targetCurrencies,
	}

	// Log and make sure data was returned
	log.Printf("Retrieved exchange rates: %+v", exchangeRatesResponse)

	return exchangeRatesResponse, nil
}
