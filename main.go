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
		log.Println("No port has been specified. Port has been set to default: " + resources.DefaultPort)
		port = resources.DefaultPort
	}

	// Initializes the handlers for the different endpoints that the API uses
	http.HandleFunc(resources.RegistrationsPath, handlers.RegistrationsHandler)
	http.HandleFunc(resources.DashboardsPath, handlers.DashboardsHandler)
	http.HandleFunc(resources.NotificationsPath, handlers.WebhookHandler)
	http.HandleFunc(resources.StatusPath, handlers.StatusHandler)

	log.Println("Service is listening on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("Could not listen to port "+port, err)
		return
	}
}
