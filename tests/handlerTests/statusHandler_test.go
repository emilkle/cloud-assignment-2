package handlerTests

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/handlers"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
)

// TestStatusHandler tests the status handler
func TestStatusHandler(t *testing.T) {
	firestoreEmulator.InitializeFirestoreEmulator()
	firestoreEmulator.PopulateFirestoreWithRegistrations()
	handlers.StartTime = time.Now()

	// Mock functions used in the status handler
	functions.CheckEndpointStatusFunc = func(url string) int {
		if url == resources.REST_COUNTRIES_PATH+"/alpha/no/" || url == resources.CURRENCY_PATH+"NOK/" {
			return http.StatusOK
		} else if url == resources.OPEN_METEO_PATH {
			return http.StatusOK
		}
		return http.StatusNotFound
	}

	functions.CheckFirestoreStatusFunc = func() int {
		return http.StatusOK
	}

	functions.NumberOfRegisteredWebhooksGetFunc = func(client *firestore.Client, ctx context.Context) int {
		return 2
	}

	// Define the test cases
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkBody      bool
	}{
		{"Valid GET Request", http.MethodGet, http.StatusOK, true},
		{"Invalid POST Request", http.MethodPost, http.StatusMethodNotAllowed, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "/status", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.StatusHandler)

			handler.ServeHTTP(rr, req)

			// Check if the status code is what was expected
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			// Check if response body is what was expected
			if tc.checkBody {
				expected, _ := json.Marshal(resources.StatusResponse{
					CountriesApi:   http.StatusOK,
					MeteoApi:       http.StatusOK,
					CurrencyApi:    http.StatusOK,
					NotificationDB: http.StatusOK,
					Webhooks:       2,
					Version:        "V1",
					Uptime:         math.Round(time.Since(handlers.StartTime).Seconds()),
				})
				if !bytes.Equal(rr.Body.Bytes(), expected) {
					t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expected))
				}
			}
		})
	}
}
