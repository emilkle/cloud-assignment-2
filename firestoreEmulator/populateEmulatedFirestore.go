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

var ctx context.Context
var client *firestore.Client
var err error

func InitializeFirestoreEmulator() {
	err1 := os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8081")
	if err1 != nil {
		log.Printf("error setting environment variable: %v\n", err)
		return
	}
	log.Println("Emulator host has been set to: " + os.Getenv("FIRESTORE_EMULATOR_HOST"))

	ctx = context.Background()
	client, err = firestore.NewClient(ctx, "countries-dashboard-service",
		option.WithEndpoint(os.Getenv("FIRESTORE_EMULATOR_HOST")), option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("Failed to connect to Firestore emulator: %v", err)
		return
	}
}

// PopulateFirestoreWithRegistrations populates an emulated Firestore to be used for testing purposes
func PopulateFirestoreWithRegistrations() {
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

	docs, err2 := client.Collection(resources.REGISTRATIONS_COLLECTION).Documents(ctx).GetAll()
	if err2 != nil {
		log.Println("Failed to retrieve documents: ", err2.Error())
		return
	}

	// Iterate over registrations and add them to Firestore
	if len(docs) < 2 {
		for _, reg := range registrations {
			_, _, err3 := client.Collection(resources.REGISTRATIONS_COLLECTION).Add(ctx, reg)
			if err3 != nil {
				log.Printf("Failed to add registration: %v", err3)
			} else {
				log.Println("Registration added to the Firestore collection.")
			}
		}
	} else {
		log.Println("There are already 2 documents in the collection. Skipping document addition.")
	}
}

// PopulateFirestoreWithRegistrations populates an emulated Firestore to be used for testing purposes
func PopulateFirestoreWithWebhooks() {
	webhooks := []map[string]interface{}{
		{
			"ID":      "1",
			"URL":     "URL1",
			"Country": "NO",
			"Event":   "POST",
		},
		{
			"ID":      "2",
			"URL":     "URL2",
			"Country": "NO",
			"Event":   "POST",
		},
	}

	//Iterate over webhooks and add them to firestore
	webhookDocs, err4 := client.Collection(resources.WEBHOOK_COLLECTION).Documents(ctx).GetAll()
	if err4 != nil {
		log.Println("Failed to retrieve documents: ", err4.Error())
		return
	}

	// Iterate over webhooks and add them to Firestore
	if len(webhookDocs) < 2 {
		for _, reg := range webhooks {
			_, _, err5 := client.Collection(resources.WEBHOOK_COLLECTION).Add(ctx, reg)
			if err5 != nil {
				log.Printf("Failed to add webhook: %v", err5)
			} else {
				log.Println("Webhook added to the Firestore collection.")
			}
		}
	} else {
		log.Println("There are already 2 documents in the collection. Skipping document addition.")
	}
}

// Server function to handle HTTP requests to populate Emulated Firestore
func StartServer() {
	http.HandleFunc("/populate", func(w http.ResponseWriter, r *http.Request) {
		PopulateFirestoreWithRegistrations()
		PopulateFirestoreWithWebhooks()
		w.Write([]byte("Emulated Firestore populated successfully"))
	})

	log.Println("Server starting on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// GetEmulatorClient() gets the firestore client
func GetEmulatorClient() *firestore.Client {
	return client
}

// GetEmulatorContext() gets the firestore context
func GetEmulatorContext() context.Context {
	return ctx
}
