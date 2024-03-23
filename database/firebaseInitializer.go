package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var client *firestore.Client

func InitializeFirestore() error {
	// Connection to Firebase
	opt := option.WithCredentialsFile("./database/fb_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("error initializing app: %v\n", err)
		return err
	}

	//Initialize client
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Printf("error initializing Firestore client: %v\n", err)
		return err
	}

	// Set the context
	ctx = context.Background()

	return nil
}

// GetFirestoreClient() gets the firestore client
func GetFirestoreClient() *firestore.Client {
	return client
}

// GetFirestoreContext() gets the firestore context
func GetFirestoreContext() context.Context {
	return ctx
}
