package registrations

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"strconv"
)

func DeleteDocumentWithRequestedId(ctx context.Context, client *firestore.Client, w http.ResponseWriter,
	requestedId int) {

	iter := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==",
		requestedId).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error iterating over query results: %v", err)
		}

		// Retrieve the DocumentRef for the document
		documentToDelete := doc.Ref

		_, err2 := client.Collection(resources.REGISTRATIONS_COLLECTION).
			Doc(documentToDelete.ID).Delete(ctx)
		if err2 != nil {
			log.Println("Error when deleting the given document. ", err2)
			http.Error(w, "Error when deleting the given document. ", http.StatusForbidden)
		}

	}
	log.Println("The document(s) with ID(s) " + strconv.Itoa(requestedId) + " was successfully deleted.")
	http.Error(w, "The document(s) with ID(s) "+strconv.Itoa(requestedId)+
		" was successfully deleted.", http.StatusNoContent)
}
