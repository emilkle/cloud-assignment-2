package handlers

import (
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
	//err1 := json.Unmarshal(content, &postRegistrationBody)
	if err1 != nil {
		http.Error(w, resources.DECODING_ERROR+"of the POST request.", http.StatusInternalServerError)
		return
	}

	postRegistration := resources.RegistrationsPOSTandPUT{
		Country: postRegistrationBody.Country,
		IsoCode: postRegistrationBody.IsoCode,
		Features: resources.Features{
			Temperature:      postRegistrationBody.Features.Temperature,
			Precipitation:    postRegistrationBody.Features.Precipitation,
			Capital:          postRegistrationBody.Features.Capital,
			Coordinates:      postRegistrationBody.Features.Coordinates,
			Population:       postRegistrationBody.Features.Population,
			Area:             postRegistrationBody.Features.Area,
			TargetCurrencies: postRegistrationBody.Features.TargetCurrencies,
		},
	}

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
		// Returns document ID in body
		log.Println("Document added to the registrations collection. " +
			"Identifier of the added document: " + documentId.ID)
		//http.Error(w, documentId.ID, http.StatusCreated)

		// TODO add the fields in the POST response to the new document. Like done here:
		// https://cloud.google.com/firestore/docs/samples/firestore-data-set-nested-fields
		postResponse, _ := functions.ParsePostResponse(client)

		standardResponseWriter(w, postResponse)
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
