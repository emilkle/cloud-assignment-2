package main

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"log"
	"net/http"
	"os"
)

//Firebase context and client used by Firestore functions throughout the program.
//var ctx context.Context
//var client *firestore.Client

func main() {

	// Initialize firestore
	if err := database.InitializeFirestore(); err != nil {
		log.Fatalf("Failed to initalize firestore: %v", err)
	}

	//ctx = context.Background()

	/**
	// Connection to Firebase
	opt := option.WithCredentialsFile("./database/fb_key.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println("error initializing app: %v", err)
		return
	}

	//Initialize client
	client, err = app.Firestore(ctx)
	*/

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
	//http.HandleFunc(NOTIFICATIONS_PATH, handler.)
	http.HandleFunc(resources.STATUS_PATH, handlers.StatusHandler)
	//handler.StartTimer()

	log.Println("Service is listening on port " + port)
	http.ListenAndServe(":"+port, nil)
}
