package registrationsTests

import (
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var allRegistrations = []resources.RegistrationsGET{
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
}

var testRegistration = resources.RegistrationsPOSTandPUT{
	Country: "Spain",
	IsoCode: "ES",
	Features: resources.Features{
		Temperature:      false,
		Precipitation:    true,
		Capital:          false,
		Coordinates:      true,
		Population:       false,
		Area:             true,
		TargetCurrencies: []string{"EUR", "NOK"},
	},
}

func TestCreatePOSTRequest(t *testing.T) {
	SetupFirestoreDatabase()
	invalidRegistrationStructCountry := resources.RegistrationsPOSTandPUT{
		Country: "",
		IsoCode: "",
		Features: resources.Features{
			Temperature:      false,
			Precipitation:    false,
			Capital:          false,
			Coordinates:      false,
			Population:       false,
			Area:             false,
			TargetCurrencies: nil,
		},
	}

	invalidRegistrationStructCurrencies := resources.RegistrationsPOSTandPUT{
		Country: "Norway",
		IsoCode: "NO",
		Features: resources.Features{
			Temperature:      false,
			Precipitation:    false,
			Capital:          false,
			Coordinates:      false,
			Population:       false,
			Area:             false,
			TargetCurrencies: []string{"", ""},
		},
	}

	invalidErrorCountry := errors.New("'country' field is not a string")

	invalidErrorCurrencies := errors.New("element:of 'targetCurrencies' field is not a string, " +
		"or the array is not a string array")

	tests := []struct {
		name    string
		data    resources.RegistrationsPOSTandPUT
		want    string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "The POST request is valid and the length of the document id is correct," +
				" and thus the document is added",
			data:    testRegistration,
			want:    "1234567890polikjas23",
			wantErr: nil,
		},
		{
			name:    "The POST request body has an invalid country field format",
			data:    invalidRegistrationStructCountry,
			want:    "",
			wantErr: invalidErrorCountry,
		},
		{
			name:    "The POST request body has an invalid targetCurrency field format",
			data:    invalidRegistrationStructCurrencies,
			want:    "",
			wantErr: invalidErrorCurrencies,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			got, err := registrations.CreatePOSTRequest(emulatorCtx, emulatorClient, w, tt.data)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreatePOSTRequest() error = %v, \n wantErr = %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("CreatePOSTRequest() got document ID with length = %v, "+
					"expected document ID with length  %v", len(got), len(tt.want))
			}
		})
	}
}

func TestCreatePOSTResponse(t *testing.T) {
	expectedResponse := resources.RegistrationsPOSTResponse{
		Id:         5,
		LastChange: "20240229 14:07",
	}

	expectedJsonResponse, err := json.Marshal(expectedResponse)
	if err != nil {
		fmt.Println(resources.ENCODING_ERROR)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextId := len(allRegistrations) + 1
		postResponse := resources.RegistrationsPOSTResponse{
			Id:         nextId,
			LastChange: "20240229 14:07",
		}

		jsonResponse, err1 := json.Marshal(postResponse)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))

	tests := []struct {
		name         string
		method       string
		expectedBody []byte
		expectedCode int
	}{
		// TODO: Add test cases.
		{
			name:         "The new registration has the next id in line",
			method:       http.MethodPost,
			expectedBody: expectedJsonResponse,
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err1 := http.NewRequest(tt.method, mockServer.URL, nil)
			if err1 != nil {
				t.Fatalf("failed to create request: %v", err1)
			}

			resp, err2 := http.DefaultClient.Do(req)
			if err2 != nil {
				t.Fatalf("failed to send request: %v", err2)
			}
			defer resp.Body.Close()

			body, err3 := io.ReadAll(resp.Body)
			if err3 != nil {
				t.Fatalf("failed to read response body: %v", err3)
			}

			if string(body) != string(tt.expectedBody) {
				t.Errorf("Expected response body %q but got %q.", tt.expectedBody, string(body))
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}

func TestUpdatePOSTRequest(t *testing.T) {
	SetupFirestoreDatabase()

	var postRegistration = map[string]interface{}{
		"country": "Denmark",
		"isoCode": "DK",
		"features": map[string]interface{}{
			"temperature":      true,
			"precipitation":    true,
			"capital":          true,
			"coordinates":      true,
			"population":       true,
			"area":             false,
			"targetCurrencies": []interface{}{"EUR", "USD", "SEK"},
		},
	}

	newDocumentRef, _, err1 := emulatorClient.Collection(resources.REGISTRATIONS_COLLECTION).Add(emulatorCtx,
		postRegistration)
	if err1 != nil {
		log.Println("An error occurred when creating a new document:", err1.Error())
	}

	postResponse := resources.RegistrationsPOSTResponse{
		Id:         3,
		LastChange: "20240405 18:07",
	}

	documentID := newDocumentRef.ID

	tests := []struct {
		name         string
		documentID   string
		postResponse resources.RegistrationsPOSTResponse
		expectedCode int
	}{
		// TODO: Add test cases.
		{
			name:         "The POST request body was successfully updated",
			documentID:   documentID,
			postResponse: postResponse,
			expectedCode: http.StatusOK,
		},
		{
			name:         "The lastChange and id fields could not be updated",
			documentID:   "",
			postResponse: postResponse,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			registrations.UpdatePOSTRequest(emulatorCtx, emulatorClient, w, tt.documentID, tt.postResponse)
			if w.Code != tt.expectedCode {
				t.Errorf("Expected HTTP status code %d, but got %d", http.StatusInternalServerError, w.Code)
			}
		})
	}
}

func TestValidateDataTypes(t *testing.T) {
	invalidRegistrationStructIsoCodes := resources.RegistrationsPOSTandPUT{
		Country: "Norway",
		IsoCode: "",
		Features: resources.Features{
			Temperature:      false,
			Precipitation:    false,
			Capital:          false,
			Coordinates:      false,
			Population:       false,
			Area:             false,
			TargetCurrencies: []string{"", ""},
		},
	}

	invalidRegistrationStructCurrenciesNil := resources.RegistrationsPOSTandPUT{
		Country: "Norway",
		IsoCode: "NO",
		Features: resources.Features{
			Temperature:      false,
			Precipitation:    false,
			Capital:          false,
			Coordinates:      false,
			Population:       false,
			Area:             false,
			TargetCurrencies: nil,
		},
	}

	invalidRegistrationStructCurrenciesEmptyStrings := resources.RegistrationsPOSTandPUT{
		Country: "Norway",
		IsoCode: "NO",
		Features: resources.Features{
			Temperature:      false,
			Precipitation:    false,
			Capital:          false,
			Coordinates:      false,
			Population:       false,
			Area:             false,
			TargetCurrencies: []string{"", ""},
		},
	}

	invalidErrorIsoCode := errors.New("'isoCode' field is not a string")

	invalidErrorTargetCurrencies := errors.New("element:of 'targetCurrencies' field is not a string, " +
		"or the array is not a string array")

	tests := []struct {
		name    string
		data    resources.RegistrationsPOSTandPUT
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "The POST request body has no invalid field formats",
			data:    testRegistration,
			wantErr: nil,
		},
		{
			name:    "The POST request body has an invalid isoCode field format",
			data:    invalidRegistrationStructIsoCodes,
			wantErr: invalidErrorIsoCode,
		},
		{
			name:    "The POST request body has a targetCurrencies field that is nil",
			data:    invalidRegistrationStructCurrenciesNil,
			wantErr: invalidErrorTargetCurrencies,
		},
		{
			name: "The POST request body has an invalid targetCurrency field where one" +
				" or more of the currencies are empty strings",
			data:    invalidRegistrationStructCurrenciesEmptyStrings,
			wantErr: invalidErrorTargetCurrencies,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if err := registrations.ValidateDataTypes(tt.data, w); err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ValidateDataTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
