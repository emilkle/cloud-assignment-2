package registrations

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"time"
)

// ValidateDataTypes validates the data structure representing a POST-registration or a PUT-registration payload.
func ValidateDataTypes(data resources.RegistrationsPOSTandPUT, w http.ResponseWriter) error {
	// Check if the 'country' field is a string
	if reflect.TypeOf(data.Country) != reflect.TypeOf("") || data.Country == "" {
		err := errors.New("'country' field is not a string")
		log.Println(resources.StandardDatatypeError, err.Error())
		http.Error(w, "Please use a 'string' as datatype for the country field.", http.StatusForbidden)
		return err
	}

	// Check if the 'isoCode' field is a string
	if reflect.TypeOf(data.IsoCode) != reflect.TypeOf("") || data.IsoCode == "" {
		err := errors.New("'isoCode' field is not a string")
		log.Println(resources.StandardDatatypeError, err.Error())
		http.Error(w, "Please use a 'string' as datatype for the isoCode field.", http.StatusForbidden)
		return err
	}

	// Check if the 'targetCurrencies' field is a slice of strings
	targetCurrencies := data.Features.TargetCurrencies
	for _, tc := range targetCurrencies {
		if (reflect.TypeOf(tc) != reflect.TypeOf("")) || tc == "" || targetCurrencies == nil {
			err := errors.New("element:" + tc + "of 'targetCurrencies' field is not a string, " +
				"or the array is not a string array")
			log.Println(resources.StandardDatatypeError, err.Error())
			http.Error(w, "Please use a 'string' as datatype for the elements in the targetCurrencies array.",
				http.StatusForbidden)
			return err
		}
	}
	return nil
}

// CreatePOSTRequest creates a new registration document in Firestore.
// It takes the provided data, constructs a document,
// adds it to the database, and returns the document ID along with an error, if any.
func CreatePOSTRequest(ctx context.Context, client *firestore.Client, w http.ResponseWriter,
	data resources.RegistrationsPOSTandPUT) (string, error) {
	err := ValidateDataTypes(data, w)
	if err != nil {
		log.Println("The document has incorrect datatypes:", err.Error())
		http.Error(w, "The input datatypes or document structure is incorrect. Please use the following"+
			"format to add a new document: "+resources.JsonStructPostAndPut, http.StatusForbidden)
		return "", err
	} else {

		// Construct registration document data.
		postRegistration := map[string]interface{}{
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
			}}

		// Create a new document with an autogenerated ID.
		// If everything is successful it will log a message in the server log with the newly
		// created documents autogenerated ID.
		newDocumentRef, _, err1 := client.Collection(resources.RegistrationsCollection).Add(ctx, postRegistration)
		if err1 != nil {
			log.Println("An error occurred when creating a new document:", err1.Error())
			http.Error(w, "An error occurred when creating a new document.", http.StatusInternalServerError)
			return "", err1
		} else {
			log.Println("Document added to the registrations collection. " +
				"Identifier of the added document: " + newDocumentRef.ID)
		}
		return newDocumentRef.ID, nil
	}
}

// CreatePOSTResponse creates a response for a successful POST request.
// It retrieves all registered documents, calculates the ID for the next document,
// and returns a response along with an error, if any.
func CreatePOSTResponse(ctx context.Context, client *firestore.Client, w http.ResponseWriter) (resources.RegistrationsPOSTResponse, error) {
	allDocuments, err := GetAllRegisteredDocuments(ctx, client)
	if err != nil {
		log.Println("Could not fetch all documents: ", err.Error())
		http.Error(w, "Could not fetch all documents from the firestore database.",
			http.StatusInternalServerError)
		return resources.RegistrationsPOSTResponse{}, err
	}
	nextId := len(allDocuments) + 1

	return resources.RegistrationsPOSTResponse{
		Id:         nextId,
		LastChange: time.Now().Format("20060102 15:04"),
	}, nil
}

// UpdatePOSTRequest updates the newly created registration document with its ID and last change timestamp.
func UpdatePOSTRequest(ctx context.Context, client *firestore.Client, w http.ResponseWriter,
	documentID string, postResponse any) {
	postResponseMap := make(map[string]interface{})
	jsonBytes, err1 := json.Marshal(&postResponse)
	if err1 != nil {
		log.Println("Unable to marshal the POST response: ", err1.Error())
		http.Error(w, resources.EncodingError+"of the POST response data.", http.StatusInternalServerError)
		return
	}

	err2 := json.Unmarshal(jsonBytes, &postResponseMap)
	if err2 != nil {
		log.Println("Unable to unmarshal the POST response: ", err2.Error())
		http.Error(w, resources.DecodingError+"of the POST response data.", http.StatusInternalServerError)
		return
	}

	// Update document with id and lastChange fields.
	_, err3 := client.Collection(resources.RegistrationsCollection).Doc(documentID).Set(ctx,
		postResponseMap, firestore.MergeAll)

	if err3 != nil {
		log.Println("The id and lastChange fields could not be set: ", err3.Error())
		http.Error(w, "An error occurred when setting the id and last change"+
			" timestamp of the new registration, Please try again. ", http.StatusInternalServerError)
	}

}
