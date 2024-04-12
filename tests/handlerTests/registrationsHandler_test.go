package handlerTests

import (
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
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

func TestRegistrationsHandler(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
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
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		method       string
		server       *httptest.Server
		expectedCode int
	}{
		// TODO: Add test cases.
		{"Method = GET (Status OK)", http.MethodGet, mockServer,
			http.StatusOK},
		{"Method = POST (Status OK)", http.MethodPost, mockServer,
			http.StatusOK},
		{"Method = PUT (Status OK)", http.MethodPut, mockServer,
			http.StatusOK},
		{"Method = DELETE (Status OK)", http.MethodDelete, mockServer,
			http.StatusOK},
		{"Method = OPTIONS (Status not implemented)", http.MethodOptions, mockServer,
			http.StatusNotImplemented},
		{"Method = HEAD (Status not implemented)", http.MethodHead, mockServer,
			http.StatusNotImplemented},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, mockServer.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("handler returned wrong status code for method %s: got %v want %v",
					tt.method, resp.StatusCode, tt.expectedCode)
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

func Test_RegistrationRequestGET(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	defer mockServer.Close()

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
		{
			name:               "Get all registrations",
			server:             mockServer,
			url:                mockServer.URL,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, tt.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			handlers.RegistrationRequestGET(w, request)

			if statusCode := w.Code; statusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedStatusCode, statusCode)
			}

		})
	}
}

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
