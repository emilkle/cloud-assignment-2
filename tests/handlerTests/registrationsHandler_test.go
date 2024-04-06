package handlerTests

import (
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Define an interface for handling registration requests
type RegistrationsHandler interface {
	HandleRegistrations(w http.ResponseWriter, r *http.Request)
}

// Real implementation of RegistrationsHandler
type RealRegistrationHandler struct{}

func (rh *RealRegistrationHandler) HandleRegistrations(w http.ResponseWriter, r *http.Request) {
	handlers.RegistrationsHandler(w, r)
}

// Mock implementation of RegistrationsHandler
type MockRegistrationsHandler struct {
	HandleRegistrationsFunc func(w http.ResponseWriter, r *http.Request)
}

func (mh *MockRegistrationsHandler) HandleRegistrations(w http.ResponseWriter, r *http.Request) {
	// Call the mock function provided by the test
	mh.HandleRegistrationsFunc(w, r)
}

func SetRegistrationsHandler(handler RegistrationsHandler) {
	registrationsHandler = handler
}

// Global variable to hold the registrations handler
var registrationsHandler RegistrationsHandler

// RegistrationsHandlerFunc is a function type that wraps the method HandleRegistrations
type RegistrationsHandlerFunc func(w http.ResponseWriter, r *http.Request)

// HandleRegistrations is the handler function for registrations
func (f RegistrationsHandlerFunc) HandleRegistrations(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func NewRegistrationsHandler(f func(w http.ResponseWriter, r *http.Request)) RegistrationsHandler {
	return RegistrationsHandlerFunc(f)
}

/*
type MockRegistrations struct{}

func (m *MockRegistrations) GetAllRegisteredDocuments() ([]resources.RegistrationsGET, error) {
	return []resources.RegistrationsGET{
		{
			Id:      1,
			Country: "Norway",
			IsoCode: "NO",
			Features: resources.Features{
				Temperature:      true,
				Precipitation:    true,
				Capital:          true,
				Coordinates:      false,
				Population:       true,
				Area:             false,
				TargetCurrencies: []string{"EUR", "USD", "SEK"},
			},
			LastChange: "20240229 14:07",
		},
		{
			Id:      2,
			Country: "Chad",
			IsoCode: "TD",
			Features: resources.Features{
				Temperature:      true,
				Precipitation:    false,
				Capital:          true,
				Coordinates:      false,
				Population:       true,
				Area:             false,
				TargetCurrencies: []string{"CFA", "RUB"},
			},
			LastChange: "20240323 15:20",
		},
		{
			Id:      3,
			Country: "Sweden",
			IsoCode: "SE",
			Features: resources.Features{
				Temperature:      true,
				Precipitation:    true,
				Capital:          true,
				Coordinates:      false,
				Population:       false,
				Area:             false,
				TargetCurrencies: []string{"NOK", "SEK", "USD", "DKK"},
			},
			LastChange: "20240324 10:57",
		},
		{
			Id:      4,
			Country: "Denmark",
			IsoCode: "DK",
			Features: resources.Features{
				Temperature:      false,
				Precipitation:    true,
				Capital:          true,
				Coordinates:      true,
				Population:       false,
				Area:             true,
				TargetCurrencies: []string{"NOK", "MYR", "JPY", "EUR"},
			},
			LastChange: "20240324 16:19",
		},
	}, nil
}

func (m *MockRegistrations) CreateRegistrationsGET(id string) (resources.RegistrationsGET, error) {
	denmark := resources.RegistrationsGET{
		Id:      4,
		Country: "Denmark",
		IsoCode: "DK",
		Features: resources.Features{
			Temperature:      false,
			Precipitation:    true,
			Capital:          true,
			Coordinates:      true,
			Population:       false,
			Area:             true,
			TargetCurrencies: []string{"NOK", "MYR", "JPY", "EUR"},
		},
		LastChange: "20240324 16:19",
	}

	return denmark, nil
}*/

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

func Test_RegistrationRequestDELETE(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		server *httptest.Server
		name   string
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, resources.REGISTRATIONS_PATH, nil)
			if err != nil {
				// Handle error
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.RegistrationsHandler)
			handler.ServeHTTP(rr, req)
			handlers.RegistrationRequestDELETE(tt.args.w, tt.args.r)

			if statusCode := rr.Code; statusCode != http.StatusOK {
				t.Errorf("Expected status code %d but got %d.", http.StatusOK, statusCode)
			}
		})
	}
}

/*
func Test_RegistrationRequestGET(t *testing.T) {
	tests := []struct {
		name               string
		server             *httptest.Server
		url                string
		expectedStatusCode int
		//mockDatabaseCalls  func(*mockRegistrations)

		//w      http.ResponseWriter
		//r      *http.Request
	}{
		// TODO: Add test cases.
		{name: "Positive server response", server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))},
		{
			name:               "Get all registrations",
			request:            httptest.NewRequest(http.MethodGet, resources.REGISTRATIONS_PATH, nil),
			expectedStatusCode: http.StatusOK,
			expectedError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			//if err != nil {
			//	t.Fatalf("Failed to create request: %v", err)
			//}

			//rr := httptest.NewRecorder()
			//handler := http.HandlerFunc(handlers.RegistrationRequestGET)
			//handler.ServeHTTP(rr, tt.request)

			rr := httptest.NewRecorder()
			handlers.RegistrationRequestGET(rr, tt.request)

			if statusCode := rr.Code; statusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedStatusCode, statusCode)
			}

		})
	}
}*/

func Test_RegistrationRequestPOST(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.RegistrationRequestPOST(tt.args.w, tt.args.r)
		})
	}
}

func Test_RegistrationRequestPUT(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.RegistrationRequestPUT(tt.args.w, tt.args.r)
		})
	}
}
