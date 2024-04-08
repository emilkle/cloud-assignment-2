package dashboards

import (
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
)

// RetrieveCurrencyExchangeRates Fetches the exchange rates of currencies with NOK as base (NOK to currency)
func RetrieveCurrencyExchangeRates(id int) (resources.TargetCurrencyValues, error) {
	// Variable used in error message for HttpRequest function.
	fetching := "exchange rates"

	// Construct URL
	url := resources.CURRENCY_PATH + "NOK"

	// Make HTTP request to specified URL
	response, err := HttpRequest(url, fetching, id)
	// Defer close of response body
	defer CloseResponseBody(response.Body, fetching, id)

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
