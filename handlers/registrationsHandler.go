package handlers

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		RegistrationRequestGET(w, r)
	case http.MethodPost:
		RegistrationRequestPOST(w, r)
	case http.MethodPut:
		RegistrationRequestPUT(w, r)
	case http.MethodDelete:
		RegistrationRequestDELETE(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' is not supported. Try"+
			" '"+http.MethodGet+", "+http.MethodPost+", "+http.MethodPut+" "+
			""+"or"+" "+http.MethodDelete+"' instead. ", http.StatusNotImplemented)
		return
	}
}

// It is possible to get more documents at once by calling /dashboard/v1/registrations/1,2,3 for
// getting specific entries in specific orders, or by
// using /dashboard/v1/registrations/ to get all documents.
func RegistrationRequestGET(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	if id == "" {
		registrationsResponses, err1 := registrations.GetAllRegisteredDocuments()
		if err1 != nil {
			http.Error(w, "Could not retrieve all documents.", http.StatusInternalServerError)
			return
		}
		standardResponseWriter(w, registrationsResponses)
		return
	}

	var registrationsResponses []resources.RegistrationsGET
	var notFoundIds []string

	registrationIds := strings.Split(id, ",")
	for _, registrationId := range registrationIds {
		registrationsResponse, err2 := registrations.CreateRegistrationsGET(registrationId)
		if err2 != nil {
			notFoundIds = append(notFoundIds, registrationId)
		}
		registrationsResponses = append(registrationsResponses, registrationsResponse)
	}

	if len(notFoundIds) > 0 {
		http.Error(w, "Registration id(s) "+strings.Join(notFoundIds, ", ")+
			" could not be found.", http.StatusNotFound)
		return
	}

	standardResponseWriter(w, registrationsResponses)
}

func RegistrationRequestPOST(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, fmt.Sprintf(resources.DECODING_ERROR+"of the POST request. Use this structure for your"+
			" POST request instead: \n%s", resources.JSON_STRUCT_POST_AND_PUT), http.StatusInternalServerError)
		return
	}

	documentID, err2 := registrations.CreatePOSTRequest(ctx, client, w, postRegistrationBody)
	if err2 != nil {
		http.Error(w, "Error when creating a new document.", http.StatusInternalServerError)
		return
	}

	postResponse, _ := registrations.CreatePOSTResponse()
	w.WriteHeader(http.StatusCreated)

	standardResponseWriter(w, postResponse)

	registrations.UpdatePOSTRequest(ctx, client, w, documentID, postResponse)
}

func RegistrationRequestPUT(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	if id == "" {
		log.Println("No id was specified in this query.")
		http.Error(w, "No id was specified in this query, please write an "+
			"integer number in the query to use this service.", http.StatusBadRequest)
		return
	}

	var putRegistrationBody resources.RegistrationsPOSTandPUT
	err1 := json.NewDecoder(r.Body).Decode(&putRegistrationBody)
	if err1 != nil {
		http.Error(w, fmt.Sprintf(resources.DECODING_ERROR+"of the PUT request. Use this structure for your"+
			" PUT request instead: \n%s", resources.JSON_STRUCT_POST_AND_PUT), http.StatusInternalServerError)
		return
	}

	documentID, err2 := registrations.GetDocumentID(ctx, client, id)
	if err2 != nil {
		http.Error(w, "Error when updating the field(s) of the given document.", http.StatusInternalServerError)
		return
	}

	registrations.CreatePUTRequest(ctx, client, w, putRegistrationBody, documentID)

	log.Println("The requested document was successfully updated.")
	w.WriteHeader(http.StatusOK)
}

func RegistrationRequestDELETE(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	if id == "" {
		log.Println("No id(s) were specified in this query.")
		http.Error(w, "No id(s) were specified in this query, please write an "+
			"integer number in the query to use this service.", http.StatusBadRequest)
		return
	}

	registrationIds := strings.Split(id, ",")
	notFoundIds := registrations.DeleteDocumentWithRequestedId(ctx, client, registrationIds)

	if len(notFoundIds) > 0 {
		log.Println("The document(s) with ID(s) " + strings.Join(notFoundIds, ", ") + " were not found.")
		http.Error(w, "The document(s) with ID(s) "+strings.Join(notFoundIds, ", ")+" were not found.",
			http.StatusNotFound)
		return
	}

	log.Println("The requested documents were successfully deleted from the database.")
	http.Error(w, "The requested documents were successfully deleted "+
		"from the database.", http.StatusNoContent)
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
