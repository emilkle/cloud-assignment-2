package functionTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/resources"
	"google.golang.org/api/iterator"
	"log"
)

var emulatorClient *firestore.Client
var emulatorCtx context.Context

func GetEmulatorClient() *firestore.Client {
	return emulatorClient
}

func GetEmulatorCtx() context.Context {
	return emulatorCtx
}

// SetupFirestoreDatabase resets the firestore emulator before each test.
func SetupFirestoreDatabase(collectionPath string) {
	firestoreEmulator.InitializeFirestoreEmulator()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	iter := emulatorClient.Collection(collectionPath).
		OrderBy("lastChange", firestore.Asc).Documents(emulatorCtx)

	for {
		doc, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Failed to iterate over documents: %v", err1)
			return
		}
		_, err1 = doc.Ref.Delete(emulatorCtx)
		if err1 != nil {
			log.Printf("Failed to delete document: %v", err1)
		}
	}

	if collectionPath == resources.REGISTRATIONS_PATH {
		firestoreEmulator.PopulateFirestoreWithRegistrations()
	}
	if collectionPath == resources.WEBHOOK_COLLECTION {
		firestoreEmulator.PopulateFirestoreWithWebhooks()
	} else {
		log.Println("Wrong path has been passed as parameter")
	}

}
