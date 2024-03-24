package handlers

import (
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var startTime = time.Now()

// StatusHandler checks the status of endpoints
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure only get method/request is allowed to the endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Make map to store status codes from rest API endpoints
	status := make(map[string]int)

	//Add endpoints to the status-map
	status["countries_api"] = functions.CheckEndpointStatus(resources.REST_COUNTRIES_PATH + "/alpha/no/")
	status["currency_api"] = functions.CheckEndpointStatus(resources.CURRENCY_PATH + "NOK/")
	status["meteo_api"] = functions.CheckEndpointStatus(resources.OPEN_METEO_PATH)

	//Calculate time since server started
	uptime := time.Since(startTime).Seconds()

	//Make instance of the response struct
	statusResponse := resources.StatusResponse{
		status["countries_api"],
		status["meteo_api"],
		status["currency_api"],
		"V1",
		uptime,
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
