package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
)

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var client *firestore.Client

// InitializeFirestore initializes the production Firestore database or an emulated Firestore
// for testing, depending on if the environment variable FIRESTORE_EMULATOR_HOST is set.
func InitializeFirestore() error {
	// Set the context
	ctx = context.Background()

	var app *firebase.App
	var err error

	// Check if the FIRESTORE_EMULATOR_HOST environment variable is set
	if emulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST"); emulatorHost != "" {
		//UNCOMMENT IF YOU WANT TO TEST USING THE EMULATED FIRESTORE
		//firestoreEmulator.InitializeFirestoreEmulator()
		// Point to the Firestore Emulator
		conf := &firebase.Config{ProjectID: "countries-dashboard-service"}
		opts := option.WithEndpoint("http://" + emulatorHost)
		app, err = firebase.NewApp(ctx, conf, opts, option.WithoutAuthentication())
		//START SERVER FOR POPULATING THE EMULATED FIREBASE DATABASE
		go firestoreEmulator.StartServer()
	} else {
		// Connection to Firebase with credentials for production
		opt := option.WithCredentialsFile("./database/fb_key.json")
		app, err = firebase.NewApp(ctx, nil, opt)
	}

	if err != nil {
		log.Printf("error initializing app: %v\n", err)
		return err
	}

	// Initialize Firestore client
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Printf("error initializing Firestore client: %v\n", err)
		return err
	}

	return nil
}

// GetFirestoreClient() returns the Firestore client
func GetFirestoreClient() *firestore.Client {
	return client
}

// GetFirestoreContext() returns the Firestore context
func GetFirestoreContext() context.Context {
	return ctx
}
