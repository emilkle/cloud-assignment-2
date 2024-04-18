package registrations

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"strconv"
)

// CreateRegistrationsGET retrieves registration data based on the provided ID.
// It retrieves data from Firestore, constructs a RegistrationsGET struct, and returns it along with an error, if any.
func CreateRegistrationsGET(ctx context.Context, client *firestore.Client, idParam string) (resources.RegistrationsGET, error) {
	// Parse ID parameter to integer.
	idNumber, err1 := strconv.Atoi(idParam)
	if err1 != nil {
		log.Println("This id could not be parsed, try another id.", err1.Error())
		return resources.RegistrationsGET{}, err1
	}

	// Query Firestore for documents with matching ID.
	query := client.Collection(resources.REGISTRATIONS_COLLECTION).Where("id", "==", idNumber).Limit(1)
	documents, err2 := query.Documents(ctx).GetAll()
	if err2 != nil {
		log.Println("Failed to fetch documents:", err2)
		return resources.RegistrationsGET{}, err2
	}

	// Check if any documents were found.
	if len(documents) == 0 {
		err3 := fmt.Errorf("no document found with ID: %s", idParam)
		log.Println(err3)
		return resources.RegistrationsGET{}, err3
	}

	// Construct RegistrationsGET struct from retrieved data.
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

	// Print the error message to the server log if the retrieving of the document fails.
	log.Println("Document with ID", idParam, "was not found.")
	return resources.RegistrationsGET{}, nil
}

// GetTargetCurrencies retrieves target currencies from featuresData.
// It returns a slice of strings containing the target currencies.
func GetTargetCurrencies(featuresData map[string]interface{}) []string {
	targetCurrenciesInterface, ok := featuresData["targetCurrencies"].([]interface{})
	if !ok {
		// TargetCurrencies is not an array of interfaces.
		log.Println("The requested currency data was not found.")
	}

	// Converting []interface{} to []string.
	var targetCurrencies []string
	for _, currency := range targetCurrenciesInterface {
		if c, ok1 := currency.(string); ok1 {
			targetCurrencies = append(targetCurrencies, c)
		} else {
			// TargetCurrencies array is not a string.
			log.Println("The requested currency data is invalid.")
			return nil
		}
	}
	return targetCurrencies
}

// GetAllRegisteredDocuments retrieves all registration documents from Firestore.
// It constructs RegistrationsGET structs for each document and returns them along with an error, if any.
func GetAllRegisteredDocuments(ctx context.Context, client *firestore.Client) ([]resources.RegistrationsGET, error) {
	// Iterate over documents in ascending order of lastChange timestamp.
	iter := client.Collection(resources.REGISTRATIONS_COLLECTION).OrderBy("lastChange", firestore.Asc).Documents(ctx)
	var registrationsResponses []resources.RegistrationsGET
	idIndex := 1

	// Iterate through documents and construct RegistrationsGET structs.
	for {
		document, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			return nil, err1
		}
		data := document.Data()

		// Retrieve the lastChange timestamp from the document.
		lastChange, ok := data["lastChange"].(string)
		if !ok {
			return nil, fmt.Errorf("the timestamp of the last change %v could not be converted to string",
				data["lastChange"])
		}

		// Construct RegistrationsGET struct.
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

// CreateRegistrationsResponse constructs a RegistrationsGET struct from provided data.
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

// UpdateId updates the ID field of a document in Firestore.
func UpdateId(ctx context.Context, client *firestore.Client,
	documentID string, getResponse resources.RegistrationsGET) {

	// Update the document's id field.
	_, err := client.Collection(resources.REGISTRATIONS_COLLECTION).Doc(documentID).Set(ctx,
		map[string]interface{}{"id": getResponse.Id}, firestore.MergeAll)

	if err != nil {
		log.Println("The id field could not be set: ", err.Error())
	}
}
