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

// RegistrationsHandler handles HTTP requests related to registrations used in the countries dashboard service.
// It supports HTTP GET, POST, PUT and DELETE methods.
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		RegistrationRequestGET(w, r)
	case http.MethodPost:
		RegistrationRequestPOST(w, r)
		WebhookTrigger(http.MethodPost, w, r)
	case http.MethodPut:
		RegistrationRequestPUT(w, r)
		WebhookTrigger(http.MethodPut, w, r)
	case http.MethodDelete:
		RegistrationRequestDELETE(w, r)
		WebhookTrigger(http.MethodDelete, w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' is not supported. Try"+
			" '"+http.MethodGet+", "+http.MethodPost+", "+http.MethodPut+" "+
			""+"or"+" "+http.MethodDelete+"' instead. ", http.StatusNotImplemented)
		return
	}
}

// RegistrationRequestGET handles the HTTP GET request for registrations stored in the Firestore database.
// It is possible to get more documents at once by calling /dashboard/v1/registrations/1,2,3 for
// getting specific entries in specific orders, or by
// using /dashboard/v1/registrations/ to get all documents.
func RegistrationRequestGET(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Retrieve the 4th url-part that contains the id.
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	// Check if the query does not contain an id.
	if id == "" {
		// Fetch all the documents in the  firestore database and handle the error that it returns.
		registrationsResponses, err1 := registrations.GetAllRegisteredDocuments(ctx, client)
		if err1 != nil {
			log.Println("Could not retrieve all documents.", err1.Error())
			http.Error(w, "Could not retrieve all documents.", http.StatusInternalServerError)
			return
		}
		// Write the response using the standardResponseWriter
		standardResponseWriter(w, registrationsResponses)
		return
	}

	var registrationsResponses []resources.RegistrationsGET
	var notFoundIds []string

	// Splits the specified ids in the URL with a comma, and returns the document for each of the corresponding ids.
	// Each found document is then added to the registrationResponses array.
	registrationIds := strings.Split(id, ",")
	for _, registrationId := range registrationIds {
		registrationsResponse, err2 := registrations.CreateRegistrationsGET(ctx, client, registrationId)
		// Checks is the id is in the notFoundIds array by checking the error from the CreateRegistrationsGET function.
		// If the error is not nil it gets appended to the notFoundIds array.
		if err2 != nil {
			notFoundIds = append(notFoundIds, registrationId)
		}
		registrationsResponses = append(registrationsResponses, registrationsResponse)
	}

	// Returns an error if the notFoundIds array is longer than 0.
	if len(notFoundIds) > 0 {
		log.Println("Registration id(s) " + strings.Join(notFoundIds, ", ") +
			" could not be found.")
		http.Error(w, "Registration id(s) "+strings.Join(notFoundIds, ", ")+
			" could not be found.", http.StatusNotFound)
		return
	}

	// Writes the response using the standardResponseWriter.
	standardResponseWriter(w, registrationsResponses)
}

// RegistrationRequestPOST handles HTTP POST requests for registration data.
// It decodes the incoming request body into a Registration struct, creates a new document in the database,
// and returns a response indicating success or failure.
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

	// Decode request body into a Registration struct for POST responses.
	var postRegistrationBody resources.RegistrationsPOSTandPUT
	err1 := json.NewDecoder(r.Body).Decode(&postRegistrationBody)
	if err1 != nil {
		// Respond with decoding error if unable to decode request.
		log.Println("Unable to parse the POST request body: ", err1.Error())
		http.Error(w, fmt.Sprintf(resources.DECODING_ERROR+"of the POST request. Use this structure for your"+
			" POST request instead: \n%s", resources.JSON_STRUCT_POST_AND_PUT), http.StatusForbidden)
		return
	}

	// Create new document in the Firestore database.
	documentID, err2 := registrations.CreatePOSTRequest(ctx, client, w, postRegistrationBody)
	if err2 != nil {
		// Respond with error if unable to create new document.
		log.Println("Error when creating a new document: ", err2.Error())
		http.Error(w, "Error when creating a new document.", http.StatusInternalServerError)
		return
	}

	// Create response, write to the response body, and return an error, if any.
	postResponse, err3 := registrations.CreatePOSTResponse(ctx, client, w)
	if err3 != nil {
		// Respond with error if unable to create the response.
		log.Println("Error when creating the POST response: ", err3.Error())
		http.Error(w, "Error when creating the POST response.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	standardResponseWriter(w, postResponse)

	// Update request status in the database.
	registrations.UpdatePOSTRequest(ctx, client, w, documentID, postResponse)
}

// RegistrationRequestPUT handles HTTP PUT requests for updating registration data.
// It retrieves the document ID from the URL, decodes the incoming request body into a Registration struct,
// updates the document in the database, and returns a response indicating success or failure.
func RegistrationRequestPUT(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Retrieve the 4th url-part that contains the id.
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	// Check if the query does not contain an id.
	if id == "" {
		log.Println("No id was specified in this query.")
		http.Error(w, "No id was specified in this query, please write an "+
			"integer number in the query to use this service.", http.StatusBadRequest)
		return
	}

	// Decode request body into a Registration struct.
	var putRegistrationBody resources.RegistrationsPOSTandPUT
	err1 := json.NewDecoder(r.Body).Decode(&putRegistrationBody)
	if err1 != nil {
		// Respond with decoding error if unable to decode request body.
		log.Println("Unable to parse the PUT request body: ", err1.Error())
		http.Error(w, fmt.Sprintf(resources.DECODING_ERROR+"of the PUT request. Use this structure for your"+
			" PUT request instead: \n%s", resources.JSON_STRUCT_POST_AND_PUT), http.StatusForbidden)
		return
	}

	// Get document ID from the Firestore database.
	documentID, err2 := registrations.GetDocumentID(ctx, client, id)
	if err2 != nil {
		// Respond with error if unable to get document ID.
		log.Println("Unable to get the document ID of the given document.", err2.Error())
		http.Error(w, "Unable to get the document ID of the given document.", http.StatusInternalServerError)
		return
	}

	// Update document in the database with a documentID.
	registrations.CreatePUTRequest(ctx, client, w, putRegistrationBody, documentID)

	// Indicate if the PUT request was successful with an HTTP status code.
	log.Println("The requested document was successfully updated.")
	w.WriteHeader(http.StatusOK)
}

// RegistrationRequestDELETE handles HTTP DELETE requests for deleting registration data.
// It retrieves the document ID(s) from the URL, deletes the corresponding document(s) from the database,
// and returns a response indicating success or failure.
func RegistrationRequestDELETE(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Retrieve the 4th url-part that contains the id.
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	// Check if the query does not contain an id.
	if id == "" {
		log.Println("No id(s) were specified in this query.")
		http.Error(w, "No id(s) were specified in this query, please write an "+
			"integer number in the query to use this service.", http.StatusBadRequest)
		return
	}

	// Split multiple IDs if provided.
	registrationIds := strings.Split(id, ",")
	// Delete document(s) with requested ID(s) from the database.
	notFoundIds := registrations.DeleteDocumentWithRequestedId(ctx, client, registrationIds)

	// Respond with error if any document ID(s) not found.
	if len(notFoundIds) > 0 {
		log.Println("The document(s) with ID(s) " + strings.Join(notFoundIds, ", ") + " were not found.")
		http.Error(w, "The document(s) with ID(s) "+strings.Join(notFoundIds, ", ")+" were not found.",
			http.StatusNotFound)
		return
	}

	// Indicate if the documents were successfully deleted or not by returning an HTTP status code.
	log.Println("The requested document(s) were successfully deleted from the database.")
	http.Error(w, "The requested document(s) were successfully deleted "+
		"from the database.", http.StatusNoContent)
}

// standardResponseWriter writes the response data in JSON format to the HTTP response writer.
func standardResponseWriter(w http.ResponseWriter, response any) {
	// JSON encoding the given data.
	jsonData, err1 := json.Marshal(response)
	if err1 != nil {
		log.Println("Unable to encode the registration data: ", err1.Error())
		http.Error(w, resources.ENCODING_ERROR+"of the registrations data.", http.StatusInternalServerError)
		return
	}

	// Writing the JSON encoded data to the response body.
	_, err2 := w.Write(jsonData)
	if err2 != nil {
		log.Println("Unable to write the JSON data: ", err2.Error())
		http.Error(w, "Error while writing response.", http.StatusInternalServerError)
		return
	}
}
