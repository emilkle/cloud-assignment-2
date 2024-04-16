package firestoreEmulator

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

var ctx context.Context
var client *firestore.Client
var err error

// PopulateFirestoreData populates an emulated Firestore to be used for testing purposes
func PopulateFirestoreData() {
	err1 := os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8081")
	if err1 != nil {
		log.Printf("error setting environment variable: %v\n", err)
		return
	}
	log.Println("Emulator host has been set to: " + os.Getenv("FIRESTORE_EMULATOR_HOST"))

	ctx = context.Background()
	client, err = firestore.NewClient(ctx, "countries-dashboard-service",
		option.WithEndpoint(os.Getenv("FIRESTORE_EMULATOR_HOST")), option.WithoutAuthentication())
	//DEBUGGING
	log.Println("DEBUGGING client 2: ", client)
	if err != nil {
		log.Fatalf("Failed to connect to Firestore emulator: %v", err)
		return
	}

	// Define registration data to be inserted in the emulated Firestore
	registrations := []map[string]interface{}{
		{
			"id":      1,
			"country": "Norway",
			"isoCode": "NO",
			"features": map[string]interface{}{
				"temperature":      true,
				"precipitation":    true,
				"capital":          true,
				"coordinates":      true,
				"population":       true,
				"area":             false,
				"targetCurrencies": []string{"EUR", "USD", "SEK"},
			},
			"lastChange": "20240229 14:07",
		},
		{
			"id":      2,
			"country": "Sweden",
			"isoCode": "SE",
			"features": map[string]interface{}{
				"temperature":      true,
				"precipitation":    false,
				"capital":          true,
				"coordinates":      true,
				"population":       true,
				"area":             true,
				"targetCurrencies": []string{"EUR", "SEK"},
			},
			"lastChange": "20240301 15:10",
		},
		// more registrations can be added if needed
	}

	// Iterate over registrations and add them to Firestore
	for _, reg := range registrations {
		documentSnapshot, err2 := client.
			Collection(resources.REGISTRATIONS_COLLECTION).Doc(fmt.Sprintf("%d", reg["id"])).Get(ctx)
		if documentSnapshot.Exists() || err2 != nil {
			// Document with the same ID already exists
			log.Printf("Document with ID %d already exists, skipping.", reg["id"])
			continue
		}

		_, _, err3 := client.Collection(resources.REGISTRATIONS_COLLECTION).Add(ctx, reg)
		if err3 != nil {
			log.Printf("Failed to add registration: %v", err3)
		} else {
			log.Println("Registration added to the Firestore collection.")
		}
	}
}

// Server function to handle HTTP requests to populate Emulated Firestore
func StartServer() {
	http.HandleFunc("/populate", func(w http.ResponseWriter, r *http.Request) {
		PopulateFirestoreData()
		w.Write([]byte("Emulated Firestore populated successfully"))
	})

	log.Println("Server starting on :8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

// GetEmulatorClient() gets the firestore client
func GetEmulatorClient() *firestore.Client {
	return client
}

// GetEmulatorContext() gets the firestore context
func GetEmulatorContext() context.Context {
	return ctx
}
