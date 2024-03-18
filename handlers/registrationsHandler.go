package handlers

import "net/http"

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

func registrationRequestGET(w http.ResponseWriter, r *http.Request) {}

func registrationRequestPOST(w http.ResponseWriter, r *http.Request) {

}

func registrationRequestPUT(w http.ResponseWriter, r *http.Request) {

}

func registrationRequestDELETE(w http.ResponseWriter, r *http.Request) {

}
