package handlerTests

import (
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistrationsHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		// TODO: Add test cases.
		{"Method = GET (Status OK)", http.MethodGet, resources.REGISTRATIONS_PATH + "1",
			http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positiveRequest := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handlers.RegistrationsHandler(w, positiveRequest)

			if w.Code != tt.expectedCode {
				t.Errorf("For method %s, expected status code %d but got %d."+
					" Response body: %s", tt.method, tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
