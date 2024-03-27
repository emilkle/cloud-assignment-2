package handlers

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DashboardsGet Struct to display a dashboard and the last time it was retrieved
type DashboardsGet struct {
	Country       string             `json:"country"`
	IsoCode       string             `json:"isoCode"`
	Features      resources.Features `json:"features"`
	LastRetrieval string             `json:"last_retrieval"`
}

func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure only get method/request is allowed to the endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
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
	dashboard, err := RetrieveDashboardGet(IDs[0])
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
}

// RetrieveDashboardGet returns a single/specific dashboard based on the dashboard ID.
func RetrieveDashboardGet(dashboardId string) (DashboardsGet, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Convert/parse string to integer
	idNumber, err := strconv.Atoi(dashboardId)
	if err != nil {
		log.Printf("Failed to parse ID: %s. Error: %s", dashboardId, err)
		return DashboardsGet{}, err
	}

	// Make query to the database to return all documents based on the specified ID
	query := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==", idNumber).Limit(1)
	documents, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to fetch documents. Error: %s", err)
		return DashboardsGet{}, err
	}

	// Check if any document with the specified ID were found
	if len(documents) == 0 {
		err := fmt.Errorf("no document found with ID: %s", dashboardId)
		log.Println(err)
		return DashboardsGet{}, err
	}

	// Create a timestamp for the last time this dashboard was retrieved
	var lastRetrieved = time.Now().Format("20060102 15:04")

	// Take only the first document returned by the query
	data := documents[0].Data()
	featuresData := data["features"].(map[string]interface{})

	// Returns a dashboard
	return DashboardsGet{
		Country: data["country"].(string),
		IsoCode: data["isoCode"].(string),
		Features: resources.Features{
			Temperature:      featuresData["temperature"].(bool),
			Precipitation:    featuresData["precipitation"].(bool),
			Capital:          featuresData["capital"].(bool),
			Coordinates:      featuresData["coordinates"].(bool),
			Population:       featuresData["population"].(bool),
			Area:             featuresData["area"].(bool),
			TargetCurrencies: functions.GetTargetCurrencies(featuresData),
		},
		LastRetrieval: lastRetrieved,
	}, nil
}
