package handlers

import (
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

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
	dashboard, err := dashboards.RetrieveDashboardGet(IDs[0])
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
