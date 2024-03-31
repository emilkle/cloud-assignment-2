package registrations

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"google.golang.org/api/iterator"
	"log"
	"strconv"
)

func DeleteDocumentWithRequestedId(ctx context.Context, client *firestore.Client,
	requestedIds []string) []string {

	var notFoundIds []string
	for _, requestedId := range requestedIds {
		requestedIdInt, err := strconv.Atoi(requestedId)
		if err != nil {
			log.Println("Registration id " + requestedId + " could not be parsed: " + err.Error())
			continue
		}

		found := false
		iter := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==",
			requestedIdInt).Documents(ctx)

		for {
			document, err1 := iter.Next()
			if err1 == iterator.Done {
				break
			}
			if err1 != nil {
				log.Fatalf("Error iterating over query results: %v", err1)
				continue
			}

			found = true
			documentToDelete := document.Ref.ID

			_, err2 := client.Collection(resources.REGISTRATIONS_COLLECTION).Doc(documentToDelete).Delete(ctx)
			if err2 != nil {
				log.Println("An error occurred when deleting the given document. ", err2)
				continue
			}
		}

		if !found {
			notFoundIds = append(notFoundIds, requestedId)
		}
	}
	return notFoundIds
}
