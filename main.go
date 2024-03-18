package countries_dashboard_service

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// This is needed to make render use the correct port set by their environment variables.
	port := os.Getenv("PORT")
	if port == "" {
		// If the port is not specified it is set to port 8080.
		log.Println("No port has been specified. Port has been set to default: " + DEFAULT_PORT)
		port = DEFAULT_PORT
	}

	// Initializes the handlers for the different endpoints that the API uses
	//http.HandleFunc(REGISTRATIONS_PATH, handler.)
	//http.HandleFunc(DASHBOARDS_PATH, handler.)
	//http.HandleFunc(NOTIFICATIONS_PATH, handler.)
	//http.HandleFunc(STATUS_PATH, handler.StatusHandler)
	//handler.StartTimer()

	log.Println("Service is listening on port " + port)
	http.ListenAndServe(":"+port, nil)
}
