package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func LaunchServer() {
	// Handle port assignment
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Set up handler endpoints
	//http.HandleFunc()

	fmt.Println("starting server on port " + port + ".")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
