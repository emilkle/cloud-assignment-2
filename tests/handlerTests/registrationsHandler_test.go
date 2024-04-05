package handlerTests

import (
	"countries-dashboard-service/resources"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RegistrationsHandlerMock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		w.WriteHeader(http.StatusOK)
	case http.MethodPut:
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "REST method '"+r.Method+"' is not supported. Try"+
			" '"+http.MethodGet+", "+http.MethodPost+", "+http.MethodPut+" "+
			""+"or"+" "+http.MethodDelete+"' instead. ", http.StatusNotImplemented)
		return
	}
}

func TestRegistrationsHandler(t *testing.T) {
	// Create a mock FirestoreService

	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		// TODO: Add test cases.
		{"Method = GET (Status OK)", http.MethodGet, resources.REGISTRATIONS_PATH,
			http.StatusOK},
		{"Method = POST (Status OK)", http.MethodPost, resources.REGISTRATIONS_PATH,
			http.StatusOK},
		{"Method = PUT (Status OK)", http.MethodPut, resources.REGISTRATIONS_PATH,
			http.StatusOK},
		{"Method = DELETE (Status OK)", http.MethodDelete, resources.REGISTRATIONS_PATH,
			http.StatusOK},
		{"Method = OPTIONS (Status not implemented)", http.MethodOptions, resources.REGISTRATIONS_PATH,
			http.StatusNotImplemented},
		{"Method = HEAD (Status not implemented)", http.MethodHead, resources.REGISTRATIONS_PATH,
			http.StatusNotImplemented},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			RegistrationsHandlerMock(w, request)

			if w.Code != tt.expectedCode {
				t.Errorf("For method %s, expected status code %d but got %d."+
					" Response body: %s", tt.method, tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
