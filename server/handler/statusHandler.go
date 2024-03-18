package handler

import (
	"countries-dashboard-service/server/endpoints"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var startTime = time.Now()

func CheckEndpointStatus(url string) int {
	statusResponse, err := http.Get(url)
	if err != nil {
		return http.StatusServiceUnavailable
	}
	defer func(Body io.ReadCloser) {
		if err != nil {
			err := Body.Close()
			if err != nil {
				fmt.Printf("failed to close response body from endpoint: %s, during status check. %v", url, err)
			}
		}
	}(statusResponse.Body)
	return statusResponse.StatusCode
}

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
	status["countries_api"] = CheckEndpointStatus(endpoints.Restcountries + "no/")
	status["currency_api"] = CheckEndpointStatus(endpoints.Currency + "NOK/")

	//Calculate time since server started
	uptime := time.Since(startTime).Seconds()

	//Struct for the response body
	type StatusResponse struct {
		CountriesApi int     `json:"countries_api"`
		CurrencyApi  int     `json:"currency_api"`
		Version      string  `json:"version"`
		Uptime       float64 `json:"uptime"`
	}

	//Make instance of the response struct
	statusResponse := StatusResponse{
		status["countries_api"],
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
