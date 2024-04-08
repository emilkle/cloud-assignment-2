package main

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialize firestore database.
	if err := database.InitializeFirestore(); err != nil {
		log.Fatalf("Failed to initalize firestore: %v", err)
	}

	// This is needed to make render use the correct port set by their environment variables.
	port := os.Getenv("PORT")
	if port == "" {
		// If the port is not specified it is set to port 8080.
		log.Println("No port has been specified. Port has been set to default: " + resources.DEFAULT_PORT)
		port = resources.DEFAULT_PORT
	}

	// Initializes the handlers for the different endpoints that the API uses
	http.HandleFunc(resources.REGISTRATIONS_PATH, handlers.RegistrationsHandler)
	http.HandleFunc(resources.DASHBOARDS_PATH, handlers.DashboardsHandler)
	http.HandleFunc(resources.NOTIFICATIONS_PATH, handlers.WebhookHandler)
	http.HandleFunc(resources.TEMP_WEBHOOK_INV, handlers.ServiceHandler)
	http.HandleFunc(resources.STATUS_PATH, handlers.StatusHandler)

	log.Println("Service is listening on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("Could not listen to port "+port, err)
		return
	}
}
