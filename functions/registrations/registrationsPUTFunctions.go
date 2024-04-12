package registrations

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"errors"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetDocumentID retrieves the document ID based on the requested ID.
// It searches for the document with the provided ID in Firestore and returns its ID along with an error, if any.
func GetDocumentID(ctx context.Context, client *firestore.Client,
	requestedId string) (string, error) {
	// Convert requested ID to integer.
	requestedIdInt, err := strconv.Atoi(requestedId)
	if err != nil {
		log.Println("Registration id " + requestedId + " could not be parsed: " + err.Error())
		return "", err
	}

	// Query Firestore for documents with matching ID.
	iter := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==",
		requestedIdInt).Documents(ctx)

	// Iterate through query results.
	var found bool
	for {
		document, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Error iterating over query results: %v", err1)
			return "", err1
		}
		found = true
		return document.Ref.ID, nil // Return document ID if found.
	}

	// If no document found with requested ID, return an error.
	if !found {
		log.Println("The document with ID" + requestedId + " was not found.")
		return "", errors.New("document not found")
	}

	return "", nil
}

// CreatePUTRequest updates the registration document in Firestore.
// It takes the provided data, constructs an update request, and updates the document in the database.
func CreatePUTRequest(ctx context.Context, client *firestore.Client, w http.ResponseWriter,
	data resources.RegistrationsPOSTandPUT, documentID string) {

	err := ValidateDataTypes(data, w)
	if err != nil {
		log.Println("The document has incorrect datatypes:", err.Error())
		http.Error(w, "The input datatypes or document structure is incorrect. Please use the following"+
			"format to update a document: "+resources.JSON_STRUCT_POST_AND_PUT, http.StatusInternalServerError)
	} else {
		putRegistration := map[string]interface{}{
			"country": data.Country,
			"isoCode": data.IsoCode,
			"features": map[string]interface{}{
				"temperature":      data.Features.Temperature,
				"precipitation":    data.Features.Precipitation,
				"capital":          data.Features.Capital,
				"coordinates":      data.Features.Coordinates,
				"population":       data.Features.Population,
				"area":             data.Features.Area,
				"targetCurrencies": data.Features.TargetCurrencies,
			},
			"lastChange": time.Now().Format("20060102 15:04"), // Update lastChange timestamp.
		}

		// Update the document in Firestore.
		_, err3 := client.Collection(resources.REGISTRATIONS_COLLECTION).Doc(documentID).Set(ctx,
			putRegistration, firestore.MergeAll)

		if err3 != nil {
			log.Println("The lastChange field could not be changed: ", err3.Error())
			http.Error(w, "An error occurred when changing the last change"+
				" timestamp of the registration, Please try again. ", http.StatusInternalServerError)
		}
	}
}
