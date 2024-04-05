package registrations

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/database"
	"countries-dashboard-service/resources"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"strconv"
)

func CreateRegistrationsGET(idParam string) (resources.RegistrationsGET, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	idNumber, err1 := strconv.Atoi(idParam)
	if err1 != nil {
		log.Println("This id could not be parsed, try another id.", err1.Error())
		return resources.RegistrationsGET{}, err1
	}
	query := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==", idNumber).Limit(1)
	documents, err2 := query.Documents(ctx).GetAll()
	if err2 != nil {
		log.Println("Failed to fetch documents:", err2)
		return resources.RegistrationsGET{}, err2
	}

	// Check if any documents were found
	if len(documents) == 0 {
		err3 := fmt.Errorf("No document found with ID: %s", idParam)
		log.Println(err3)
		return resources.RegistrationsGET{}, err3
	}

	for _, document := range documents {
		data := document.Data()
		featuresData := data["features"].(map[string]interface{})

		return resources.RegistrationsGET{
			Id:      idNumber,
			Country: data["country"].(string),
			IsoCode: data["isoCode"].(string),
			Features: resources.Features{
				Temperature:      featuresData["temperature"].(bool),
				Precipitation:    featuresData["precipitation"].(bool),
				Capital:          featuresData["capital"].(bool),
				Coordinates:      featuresData["coordinates"].(bool),
				Population:       featuresData["population"].(bool),
				Area:             featuresData["area"].(bool),
				TargetCurrencies: GetTargetCurrencies(featuresData),
			},
			LastChange: data["lastChange"].(string),
		}, nil
	}
	log.Println("Document with ID", idParam, "was not found.")
	return resources.RegistrationsGET{}, nil
}

func GetTargetCurrencies(featuresData map[string]interface{}) []string {
	targetCurrenciesInterface, ok := featuresData["targetCurrencies"].([]interface{})
	if !ok {
		// TargetCurrencies is not an array of interfaces
		log.Println("The requested currency data was not found.")
	}

	// Converting []interface{} to []string
	var targetCurrencies []string
	for _, currency := range targetCurrenciesInterface {
		if c, ok1 := currency.(string); ok1 {
			targetCurrencies = append(targetCurrencies, c)
		} else {
			// TargetCurrencies is not a string
			log.Println("The requested currency data is invalid.")
			return nil
		}
	}
	return targetCurrencies
}

func GetAllRegisteredDocuments() ([]resources.RegistrationsGET, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	iter := client.Collection(resources.REGISTRATIONS_COLLECTION).OrderBy("lastChange", firestore.Asc).Documents(ctx)

	var registrationsResponses []resources.RegistrationsGET

	idIndex := 1
	for {
		document, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			return nil, err1
		}
		data := document.Data()

		// Retrieve the lastChange timestamp from the document
		lastChange, ok := data["lastChange"].(string)
		if !ok {
			log.Printf("The timestamp of the last change"+
				" %v could not be converted to string.", data["lastChange"])
			continue
		}

		registrationsResponse := CreateRegistrationsResponse(data, lastChange, idIndex)

		registrationID := document.Ref.ID

		// Update all the id fields in for the Firestore documents after deleting a document in the middle of the
		// ascending order, to ensure that all registration documents will be found.
		UpdateId(ctx, client, registrationID, registrationsResponse)

		registrationsResponses = append(registrationsResponses, registrationsResponse)

		idIndex++
	}

	return registrationsResponses, nil
}

func CreateRegistrationsResponse(data map[string]interface{}, lastChange string, idIndex int) resources.RegistrationsGET {
	featuresData := data["features"].(map[string]interface{})

	return resources.RegistrationsGET{
		Id:      idIndex,
		Country: data["country"].(string),
		IsoCode: data["isoCode"].(string),
		Features: resources.Features{
			Temperature:      featuresData["temperature"].(bool),
			Precipitation:    featuresData["precipitation"].(bool),
			Capital:          featuresData["capital"].(bool),
			Coordinates:      featuresData["coordinates"].(bool),
			Population:       featuresData["population"].(bool),
			Area:             featuresData["area"].(bool),
			TargetCurrencies: GetTargetCurrencies(featuresData),
		},
		LastChange: lastChange,
	}
}

func UpdateId(ctx context.Context, client *firestore.Client,
	documentID string, getResponse resources.RegistrationsGET) {

	// Update the document's id field.
	_, err := client.Collection(resources.REGISTRATIONS_COLLECTION).Doc(documentID).Set(ctx,
		map[string]interface{}{"id": getResponse.Id}, firestore.MergeAll)

	if err != nil {
		log.Println("The id field could not be set: ", err.Error())
	}
}
