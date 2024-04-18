package handlers

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

var client *firestore.Client
var ctx context.Context

// SkipRealCallOfRetrieveDashboardGet is a flag variable used when testing the DashboardsHandler
var SkipRealCallOfRetrieveDashboardGet bool

// DashboardsHandler handles HTTP requests related to dashboards used in the countries dashboard service.
// It supports HTTP GET method.
func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure only get method/request is allowed to the endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Checks if the handler is supposed request to production or emulated Firestore database
	// and set the client and context accordingly
	if !SkipRealCallOfRetrieveDashboardGet {
		client, ctx = dashboards.RecognizeEnvironmentVariableForClientContext(client, ctx)
	}

	// Extract id from url
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[len(urlParts)-1]

	IDs := strings.Split(id, ",")
	// Checks if more than one id was specified
	if len(IDs) != 1 {
		http.Error(w, "Cannot retrieve more than one dashboard, too many IDs specified.", http.StatusBadRequest)
		return
	}

	// Retrieve dashboard
	var dashboard resources.DashboardsGet
	var err error

	dashboard, err = helperFunctionForTesting(IDs)
	if err != nil {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	// Encode the dashboard response to JSON
	jsonDashboardResponse, err := json.Marshal(dashboard)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(resources.ENCODING_ERROR, err)
		return
	}

	// Set HTTP header
	w.Header().Set("Content-Type", "application/json")

	// Write response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonDashboardResponse)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		log.Println("Error writing response", err)
	}
	if !SkipRealCallOfRetrieveDashboardGet {
		WebhookTrigger(http.MethodGet, w, r)
	}
}

// helperFunctionForTesting aids the testing of the DashboardsHandler, so that no external services are
// interacted with during testing.
func helperFunctionForTesting(IDs []string) (resources.DashboardsGet, error) {
	var dashboard resources.DashboardsGet
	var err error
	if !SkipRealCallOfRetrieveDashboardGet {
		dashboard, err = dashboards.RetrieveDashboardGet(client, ctx, IDs[0], false)
	} else if IDs[0] == "1" {
		dashboard = resources.DashboardsGet{
			Country: "Norway",
			IsoCode: "NO",
			FeatureValues: resources.FeatureValues{
				Temperature:   2.0,
				Precipitation: 1.0,
				Capital:       "Oslo",
				Coordinates: resources.CoordinatesValues{
					Latitude:  62.0,
					Longitude: 10.0,
				},
				Population: 5379475,
				Area:       385180.0,
				TargetCurrencies: map[string]float64{
					"EUR": 0.086312,
					"USD": 0.998935,
					"SEK": 0.091928,
				},
			},
			LastRetrieval: "20240229 14:07",
		}
	} else {
		err = errors.New("dashboard not found")
	}
	return dashboard, err
}
