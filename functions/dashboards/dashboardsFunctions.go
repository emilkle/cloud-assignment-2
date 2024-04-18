package dashboards

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

// RetrieveDashboardData fetches the dashboard data from Firestore DB
func RetrieveDashboardData(client *firestore.Client, ctx context.Context, dashboardId string) ([]*firestore.DocumentSnapshot, int, error) {
	idNumber, err := strconv.Atoi(dashboardId)
	if err != nil {
		log.Printf("Failed to parse ID: %s. Error: %s", dashboardId, err)
		return nil, 0, err
	}
	query := client.Collection(resources.RegistrationsCollection).Where("id", "==", idNumber).Limit(1)
	documents, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to fetch documents. Error: %s", err)
		return nil, 0, err
	}
	return documents, idNumber, nil
}

// RetrieveDashboardGet returns a single/specific dashboard based on the dashboard ID.
func RetrieveDashboardGet(client *firestore.Client, ctx context.Context, dashboardId string, runTest bool) (resources.DashboardsGet, error) {
	documents, idNumber, err := RetrieveDashboardData(client, ctx, dashboardId)
	if err != nil {
		return resources.DashboardsGet{}, err
	}

	// Check if documents contain documents
	if len(documents) == 0 {
		err = fmt.Errorf("no document found with ID: %s", dashboardId)
		log.Println(err)
		return resources.DashboardsGet{}, err
	}

	var lastRetrieved = time.Now().Format("20060102 15:04")
	data := documents[0].Data()
	featuresData, ok := data["features"].(map[string]interface{})
	if !ok {
		log.Println("Error: 'features' is not a valid map[string]interface{}")
		return resources.DashboardsGet{}, fmt.Errorf("'features' field is missing or not the expected type")
	}

	// Variables for data in the dashboards
	var tempAndPrecip resources.HourlyData
	var coordinates resources.CoordinatesValues
	var capitalPopArea resources.CapitalPopulationArea
	var capital string
	var population int
	var area float64
	var meanTemperature float64
	var meanPrecipitation float64
	var selectedExchangeRates resources.TargetCurrencyValues

	// Helper variables
	var latitude float64
	var longitude float64
	var coordinateData resources.CoordinatesValues
	coordinateData, err = RetrieveCoordinates(data["country"].(string), idNumber, runTest)
	// Lat and Long used in the RetrieveTempAndPrecipitation function.
	// Because a dashboard configuration might not have the coordinates set to true
	latitude = coordinateData.Latitude
	longitude = coordinateData.Longitude

	// Checks if dashboard configuration supports coordinates
	if featuresData["coordinates"].(bool) {
		coordinates.Longitude = longitude
		coordinates.Latitude = latitude
	}

	// Retrieve capital, population and area
	if featuresData["capital"].(bool) || featuresData["population"].(bool) || featuresData["area"].(bool) {
		capitalPopArea, err = RetrieveCapitalPopulationAndArea(data["isoCode"].(string), idNumber, runTest)

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
	}

	// Retrieve temperature and precipitation data
	if featuresData["temperature"].(bool) || featuresData["precipitation"].(bool) {
		tempAndPrecip, err = RetrieveTempAndPrecipitation(latitude, longitude, idNumber, runTest)
		temperature, precipitation := CalculateMeanTemperatureAndPrecipitation(tempAndPrecip)

		//Check if dashboard configuration support temperature and precipitation
		if featuresData["temperature"].(bool) {
			meanTemperature = temperature
		}
		if featuresData["precipitation"].(bool) {
			meanPrecipitation = precipitation
		}
	}

	// Exchange rates are always shown in a dashboard
	selectedExchangeRates, err = RetrieveTargetCurrenciesAndExchangeRates(featuresData, idNumber, runTest)
	if err != nil {
		return resources.DashboardsGet{}, err
	}

	// Returns dashboard populated with values depending on the configuration
	return resources.DashboardsGet{
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

// CalculateMeanTemperatureAndPrecipitation calculates the mean temperature and precipitation
func CalculateMeanTemperatureAndPrecipitation(tempAndPrecip resources.HourlyData) (float64, float64) {
	var meanTemperature, meanPrecipitation float64

	// Calculate mean temperature
	sumTemperature := 0.0
	for _, temp := range tempAndPrecip.Temperature {
		if temp != 0.0 {
			sumTemperature += temp
		}
	}
	meanTemperature = sumTemperature / float64(len(tempAndPrecip.Time))
	// Round meanTemperature to have 1 decimal
	meanTemperature = math.Round(meanTemperature*10) / 10

	// Calculate mean precipitation
	sumPrecipitation := 0.0
	for _, prec := range tempAndPrecip.Precipitation {
		if prec != 0.0 {
			sumPrecipitation += prec
		}
	}
	meanPrecipitation = sumPrecipitation / float64(len(tempAndPrecip.Time))
	// Round meanPrecipitation to have two decimals
	meanPrecipitation = math.Round(meanPrecipitation*100) / 100

	return meanTemperature, meanPrecipitation
}

// RetrieveTargetCurrenciesAndExchangeRates retrieves the currency exchange rates displayed in a dashboard configuration
func RetrieveTargetCurrenciesAndExchangeRates(featuresData map[string]interface{}, id int, runTest bool) (resources.TargetCurrencyValues, error) {
	// Retrieve exchange rates
	exchangeRates, err := RetrieveCurrencyExchangeRates(id, runTest)
	if err != nil {
		return resources.TargetCurrencyValues{}, err
	}

	// Retrieve the target currencies as interface slice
	targetCurrenciesInterface := featuresData["targetCurrencies"].([]interface{})

	// Initialize a slice to store the string values
	targetCurrencies := make([]string, len(targetCurrenciesInterface))

	// Iterate over the interface slice and convert each element to a string
	for i, currency := range targetCurrenciesInterface {
		targetCurrencies[i] = currency.(string)
	}

	// Initialize the TargetCurrencies map within TargetCurrencyValues before using it
	selectedExchangeRates := resources.TargetCurrencyValues{
		TargetCurrencies: make(map[string]float64),
	}

	// Iterate over targetCurrencies slice and retrieve corresponding rates
	for _, currency := range targetCurrencies {
		rate, ok := exchangeRates.TargetCurrencies[currency]
		if ok {
			selectedExchangeRates.TargetCurrencies[currency] = rate
		}
	}

	// Return the currencies corresponding with the dashboard
	return selectedExchangeRates, nil
}
