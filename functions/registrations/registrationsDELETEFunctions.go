package registrations

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"google.golang.org/api/iterator"
	"log"
	"strconv"
)

// DeleteDocumentWithRequestedId deletes documents with the requested IDs from Firestore.
// It returns a slice containing IDs of documents that were not found and couldn't be deleted.
func DeleteDocumentWithRequestedId(ctx context.Context, client *firestore.Client, requestedIds []string) []string {
	var notFoundIds []string

	// Iterate through requested IDs.
	for _, requestedId := range requestedIds {
		requestedIdInt, err := strconv.Atoi(requestedId)
		if err != nil {
			log.Println("Registration id " + requestedId + " could not be parsed: " + err.Error())
			continue
		}

		// Attempt to delete document with the requested ID.
		found := FindDocumentWithId(ctx, client, requestedIdInt)
		if !found {
			notFoundIds = append(notFoundIds, requestedId)
		}
	}
	return notFoundIds
}

// FindDocumentWithId finds a document with the ID field of the provided document in Firestore and
// performs the deletion process.
// It returns true if the document is found and deleted successfully, false otherwise.
func FindDocumentWithId(ctx context.Context, client *firestore.Client, documentId int) bool {
	// Query Firestore for documents with matching ids.
	iter := client.Collection(resources.RegistrationsCollection).Where("id", "==",
		documentId).Documents(ctx)

	for {
		document, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Error iterating over query results: %v", err1)
			return false
		}

		documentToDelete := document.Ref.ID
		_, err2 := client.Collection(resources.RegistrationsCollection).Doc(documentToDelete).Delete(ctx)
		if err2 != nil {
			log.Println("An error occurred when deleting the given document. ", err2)
			return false
		}
		return true // The document was found and is successfully deleted.
	}
	return false // The document could not be found.
}
