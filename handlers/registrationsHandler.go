package handlers

import (
	"cloud.google.com/go/firestore"
	"countries-dashboard-service/database"
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		registrationRequestGET(w, r)
	case http.MethodPost:
		registrationRequestPOST(w, r)
	case http.MethodPut:
		registrationRequestPUT(w, r)
	case http.MethodDelete:
		registrationRequestDELETE(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' is not supported. Try"+
			" '"+http.MethodGet+", "+http.MethodPost+", "+http.MethodPut+" "+
			""+"or"+" "+http.MethodDelete+"' instead. ", http.StatusNotImplemented)
		return
	}
}

func registrationRequestGET(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	if id == "" {
		registrationsResponses, err1 := functions.GetAllRegisteredDocuments()
		if err1 != nil {
			http.Error(w, "Could not retrieve all documents. ", http.StatusInternalServerError)
			return
		}
		standardResponseWriter(w, registrationsResponses)
		return
	}

	var registrationsResponses []resources.RegistrationsGET

	registrationIds := strings.Split(id, ",")
	for _, registrationId := range registrationIds {
		registrationsResponse, err2 := functions.CreateRegistrationsGET(registrationId)
		if err2 != nil {
			http.Error(w, "Registration id "+registrationId+" could not be found.", http.StatusNotAcceptable)
			return
		}
		registrationsResponses = append(registrationsResponses, registrationsResponse)
	}
	standardResponseWriter(w, registrationsResponses)
}

func registrationRequestPOST(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Might use later for faster parsing of response body.
	/*
		content, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error when reading the body of the POST request.", http.StatusInternalServerError)
		}
	*/

	var postRegistrationBody resources.RegistrationsPOSTandPUT
	err1 := json.NewDecoder(r.Body).Decode(&postRegistrationBody)
	if err1 != nil {
		http.Error(w, resources.DECODING_ERROR+"of the POST request.", http.StatusInternalServerError)
		return
	}

	postRegistration := map[string]interface{}{
		"country": postRegistrationBody.Country,
		"isoCode": postRegistrationBody.IsoCode,
		"features": map[string]interface{}{
			"temperature":      postRegistrationBody.Features.Temperature,
			"precipitation":    postRegistrationBody.Features.Precipitation,
			"capital":          postRegistrationBody.Features.Capital,
			"coordinates":      postRegistrationBody.Features.Coordinates,
			"population":       postRegistrationBody.Features.Population,
			"area":             postRegistrationBody.Features.Area,
			"targetCurrencies": postRegistrationBody.Features.TargetCurrencies,
		}}

	newDocumentRef := client.Collection(resources.REGISTRATIONS_COLLECTION)
	if newDocumentRef == nil {
		http.Error(w, "The doccument reference is nil.", http.StatusInternalServerError)
		return
	}

	documentId, _, err2 := newDocumentRef.Add(ctx, postRegistration)
	if err2 != nil {
		log.Println("An error occurred when creating a new document:", err2.Error())
		http.Error(w, "An error occurred when creating a new document.", http.StatusInternalServerError)
		return
	} else {
		log.Println("Document added to the registrations collection. " +
			"Identifier of the added document: " + documentId.ID)
		//http.Error(w, documentId.ID, http.StatusCreated)

		postResponse, _ := functions.ParsePostResponse(client)

		standardResponseWriter(w, postResponse)

		postResponseMap := make(map[string]interface{})
		jsonString, err3 := json.Marshal(&postResponse)
		if err3 != nil {
			log.Println("Unable to marshal the POST response: ", err3.Error())
			http.Error(w, resources.ENCODING_ERROR+"of the POST response data.", http.StatusInternalServerError)
			return
		}
		err3 = json.Unmarshal(jsonString, &postResponseMap)
		if err3 != nil {
			log.Println("Unable to unmarshal the POST response: ", err3.Error())
			http.Error(w, resources.DECODING_ERROR+"of the POST response data.", http.StatusInternalServerError)
			return
		}

		// Update document with id and lastChange fields.
		_, err4 := client.Collection(resources.REGISTRATIONS_COLLECTION).Doc(documentId.ID).Set(ctx,
			postResponseMap, firestore.MergeAll)

		if err4 != nil {
			log.Println("The id and lastChange fields could not be set: ", err4.Error())
			http.Error(w, "An error occurred when setting the id and last change"+
				" timestamp of the new registration, Please try again. ", http.StatusInternalServerError)
		}
	}
}

func registrationRequestPUT(w http.ResponseWriter, r *http.Request) {

}

func registrationRequestDELETE(w http.ResponseWriter, r *http.Request) {

}

func standardResponseWriter(w http.ResponseWriter, response any) {
	// JSON encoding the given data.
	jsonData, err2 := json.Marshal(response)
	if err2 != nil {
		http.Error(w, resources.ENCODING_ERROR+"of the registrations data.", http.StatusInternalServerError)
		return
	}

	// Writing the JSON encoded data to the response body.
	_, err3 := w.Write(jsonData)
	if err3 != nil {
		http.Error(w, "Error while writing response.", http.StatusInternalServerError)
		return
	}
}
