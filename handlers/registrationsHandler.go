package handlers

import (
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"encoding/json"
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
	w.Header().Set("content-type", "application/json")

	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]
	if id == "" {
		http.Error(w, "Missing id parameter. ", http.StatusBadRequest)
		return
	}

	var registrationsResponses []resources.RegistrationsGET

	registrationIds := strings.Split(id, ",")
	for _, registrationId := range registrationIds {
		registrationsResponse, err1 := functions.CreateRegistrationsGET(registrationId)
		if err1 != nil {
			http.Error(w, "Registration id "+registrationId+" could not be found.", http.StatusNotAcceptable)
			return
		}
		registrationsResponses = append(registrationsResponses, registrationsResponse)
	}

	// JSON encoding the registrations data.
	jsonData, err2 := json.Marshal(registrationsResponses)
	if err2 != nil {
		http.Error(w, resources.ENCODING_ERROR+"of the registrations data.", http.StatusInternalServerError)
		return
	}

	// Writing the JSON encoded data to the response.
	_, err3 := w.Write(jsonData)
	if err3 != nil {
		http.Error(w, "Error while writing response.", http.StatusInternalServerError)
		return
	}
}

func registrationRequestPOST(w http.ResponseWriter, r *http.Request) {

}

func registrationRequestPUT(w http.ResponseWriter, r *http.Request) {

}

func registrationRequestDELETE(w http.ResponseWriter, r *http.Request) {

}
