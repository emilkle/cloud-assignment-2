package dashboards

import (
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// TestUrlRetrieveCurrencyExchangeRates variable used when testing TestUrlRetrieveCurrencyExchangeRates function
var TestUrlRetrieveCurrencyExchangeRates string

// RetrieveCurrencyExchangeRates Fetches the exchange rates of currencies with NOK as base (NOK to currency)
func RetrieveCurrencyExchangeRates(id int, runTest bool) (resources.TargetCurrencyValues, error) {
	// Variable used in error message for HttpRequest function.
	fetching := "exchange rates"

	// Construct URL
	var urlPath = resources.CurrencyPath + "NOK"
	url := ConstructUrlForApiOrTest(urlPath, TestUrlRetrieveCurrencyExchangeRates, runTest)

	// Make HTTP request to specified URL
	response, err := HttpRequest(url, fetching, id)
	// Defer close of response body
	defer CloseResponseBody(response.Body, fetching, id)

	// Check status code of response
	if response.StatusCode != http.StatusOK {
		return resources.TargetCurrencyValues{}, fmt.Errorf("HTTP error: %s", response.Status)
	}

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
