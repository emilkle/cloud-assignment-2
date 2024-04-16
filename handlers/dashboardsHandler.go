package handlers

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/database"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var client *firestore.Client
var ctx context.Context

//var client = database.GetFirestoreClient()
//var ctx = database.GetFirestoreContext()

func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure only get method/request is allowed to the endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Checks if the handler makes request to production or emulated Firestore database
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "8081" {
		client = firestoreEmulator.GetEmulatorClient()
		ctx = firestoreEmulator.GetEmulatorContext()
	} else {
		client = database.GetFirestoreClient()
		ctx = database.GetFirestoreContext()
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
	dashboard, err := dashboards.RetrieveDashboardGet(client, ctx, IDs[0], false)
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
	WebhookTrigger(http.MethodGet, w, r)
}
