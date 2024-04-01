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

func GetDocumentID(ctx context.Context, client *firestore.Client,
	requestedId string) (string, error) {

	requestedIdInt, err := strconv.Atoi(requestedId)
	if err != nil {
		log.Println("Registration id " + requestedId + " could not be parsed: " + err.Error())
		return "", err
	}

	iter := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==",
		requestedIdInt).Documents(ctx)

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
		return document.Ref.ID, nil
	}

	if !found {
		log.Println("The document with ID" + requestedId + " was not found.")
		return "", errors.New("document not found")
	}

	return "", nil
}

func CreatePUTRequest(ctx context.Context, client *firestore.Client, w http.ResponseWriter,
	data resources.RegistrationsPOSTandPUT, documentID string) {
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
		"lastChange": time.Now().Format("20060102 15:04"),
	}

	// Update the document
	_, err3 := client.Collection(resources.REGISTRATIONS_COLLECTION).Doc(documentID).Set(ctx,
		putRegistration, firestore.MergeAll)

	if err3 != nil {
		log.Println("The lastChange field could not be changed: ", err3.Error())
		http.Error(w, "An error occurred when changing the last change"+
			" timestamp of the registration, Please try again. ", http.StatusInternalServerError)
	}
}
