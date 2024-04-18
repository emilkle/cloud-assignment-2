package handlers

import (
	"countries-dashboard-service/functions"
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

var StartTime = time.Now()

// StatusHandler checks the HTTP status codes of endpoints
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure only get method/request is allowed to the endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, ctx = dashboards.RecognizeEnvironmentVariableForClientContext(client, ctx)
	//Make map to store status codes from rest API endpoints
	status := make(map[string]int)

	//Add endpoints to the status-map
	status["countries_api"] = functions.CheckEndpointStatusFunc(resources.RestCountriesPath + "/alpha/no/")
	status["currency_api"] = functions.CheckEndpointStatusFunc(resources.CurrencyPath + "NOK/")
	status["meteo_api"] = functions.CheckEndpointStatusFunc(resources.OpenMeteoPath)
	status["notification_db"] = functions.CheckFirestoreStatusFunc()
	status["webhooks"] = functions.NumberOfRegisteredWebhooksGetFunc(client, ctx)
	// TODO: Add number of webhooks

	//Calculate time since server started
	uptime := math.Round(time.Since(StartTime).Seconds())

	//Make instance of the response struct
	statusResponse := resources.StatusResponse{
		CountriesApi:   status["countries_api"],
		MeteoApi:       status["meteo_api"],
		CurrencyApi:    status["currency_api"],
		NotificationDB: status["notification_db"],
		Webhooks:       status["webhooks"],
		Version:        "V1",
		Uptime:         uptime,
	}

	//Marshal the response to JSON
	convertedResponse, err := json.Marshal(statusResponse)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(convertedResponse)
	if err != nil {
		fmt.Println("Error writing response: ", http.StatusInternalServerError)
	}
}
