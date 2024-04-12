package handlerTests

import (
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var allRegistrations = `[
    {
        "id": 1,
        "country": "Norway",
        "iso_code": "NO",
        "features": {
            "temperature": true,
            "precipitation": true,
            "capital": true,
            "coordinates": false,
            "population": true,
            "area": false,
            "target_currencies": ["EUR", "USD", "SEK"]
        },
        "last_change": "20240229 14:07"
    },
    {
        "id": 2,
        "country": "Chad",
        "iso_code": "TD",
        "features": {
            "temperature": true,
            "precipitation": false,
            "capital": true,
            "coordinates": false,
            "population": true,
            "area": false,
            "target_currencies": ["CFA", "RUB"]
        },
        "last_change": "20240323 15:20"
    },
    {
        "id": 3,
        "country": "Sweden",
        "iso_code": "SE",
        "features": {
            "temperature": true,
            "precipitation": true,
            "capital": true,
            "coordinates": false,
            "population": false,
            "area": false,
            "target_currencies": ["NOK", "SEK", "USD", "DKK"]
        },
        "last_change": "20240324 10:57"
    },
    {
        "id": 4,
        "country": "Denmark",
        "iso_code": "DK",
        "features": {
            "temperature": false,
            "precipitation": true,
            "capital": true,
            "coordinates": true,
            "population": false,
            "area": true,
            "target_currencies": ["NOK", "MYR", "JPY", "EUR"]
        },
        "last_change": "20240324 16:19"
    }
]`

var singleRegistration = `{
		"id": 3,
			"country": "Sweden",
			"iso_code": "SE",
			"features": {
			"temperature": true,
				"precipitation": true,
				"capital": true,
				"coordinates": false,
				"population": false,
				"area": false,
				"target_currencies": ["NOK", "SEK", "USD", "DKK"]
	},
	"last_change": "20240324 10:57"
	}`

var twoRegistrations = `[
{
"id": 1,
"country": "Norway",
"iso_code": "NO",
"features": {
"temperature": true,
"precipitation": true,
"capital": true,
"coordinates": false,
"population": true,
"area": false,
"target_currencies": ["EUR", "USD", "SEK"]
},
"last_change": "20240229 14:07"
},
{
"id": 2,
"country": "Chad",
"iso_code": "TD",
"features": {
"temperature": true,
"precipitation": false,
"capital": true,
"coordinates": false,
"population": true,
"area": false,
"target_currencies": ["CFA", "RUB"]
},
"last_change": "20240323 15:20"
}, ]`

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
	notFoundIds := []string{"5", "6", "7"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		id := urlParts[4]

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		for _, notFoundId := range notFoundIds {
			if notFoundId == id {
				w.WriteHeader(http.StatusNotFound)
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		// TODO: Add test cases.
		{
			name:         "Delete one registration",
			method:       http.MethodDelete,
			path:         resources.REGISTRATIONS_PATH + "3",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Delete multiple registrations",
			method:       http.MethodDelete,
			path:         resources.REGISTRATIONS_PATH + "1,2,3",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Trying to delete an registration that is not registered",
			method:       http.MethodDelete,
			path:         resources.REGISTRATIONS_PATH + "6",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "No id given",
			method:       http.MethodDelete,
			path:         resources.REGISTRATIONS_PATH + "",
			expectedCode: http.StatusBadRequest,
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, mockServer.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}

func Test_RegistrationRequestGET(t *testing.T) {
	notFoundIds := []string{"5", "6", "7"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		id := urlParts[4]

		// If trying to get all documents, no id is specified.
		if id == "" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, allRegistrations)
		}

		for _, notFoundId := range notFoundIds {
			if notFoundId == id {
				w.WriteHeader(http.StatusNotFound)
			}
		}

		if id == "3" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, singleRegistration)
		}

		if id == "1,2" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, twoRegistrations)
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		expectedBody string
		expectedCode int
	}{
		// TODO: Add test cases.
		{
			name:         "Get all registrations",
			method:       http.MethodGet,
			path:         resources.REGISTRATIONS_PATH,
			expectedBody: allRegistrations,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get one registration",
			method:       http.MethodGet,
			path:         resources.REGISTRATIONS_PATH + "3",
			expectedBody: singleRegistration,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get invalid registration",
			method:       http.MethodGet,
			path:         resources.REGISTRATIONS_PATH + "5",
			expectedBody: "",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Get multiple individual registrations",
			method:       http.MethodGet,
			path:         resources.REGISTRATIONS_PATH + "1,2",
			expectedBody: twoRegistrations,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err1 := http.NewRequest(tt.method, mockServer.URL+tt.path, nil)
			if err1 != nil {
				t.Fatalf("failed to create request: %v", err1)
			}

			resp, err2 := http.DefaultClient.Do(req)
			if err2 != nil {
				t.Fatalf("failed to send request: %v", err2)
			}

			body, err3 := io.ReadAll(resp.Body)
			if err3 != nil {
				t.Fatalf("failed to read response body: %v", err3)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("Expected response body %q but got %q.", tt.expectedBody, string(body))
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedCode, resp.StatusCode)
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
