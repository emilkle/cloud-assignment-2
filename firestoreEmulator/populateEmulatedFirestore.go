package firestoreEmulator

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

// PopulateFirestoreData populates an emulated Firestore to be used for testing purposes
func PopulateFirestoreData() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "countries-dashboard-service", option.WithEndpoint(os.Getenv("FIRESTORE_EMULATOR_HOST")), option.WithoutAuthentication())
	//DEBUGGING
	log.Println("DEBUGGING client 2: ", client)
	if err != nil {
		log.Fatalf("Failed to connect to Firestore emulator: %v", err)
		return
	}
	defer client.Close()

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
		_, _, err := client.Collection(resources.REGISTRATIONS_COLLECTION).Add(ctx, reg)
		if err != nil {
			log.Printf("Failed to add registration: %v", err)
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
