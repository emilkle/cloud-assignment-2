package registrationsTests

import (
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	SetupFirestoreDatabase()

	emulatorDocs, err := registrations.GetAllRegisteredDocuments(emulatorCtx, emulatorClient)
	if err != nil {
		log.Println("Could not fetch all documents: ", err.Error())
	}

	expectedResponse := resources.RegistrationsPOSTResponse{
		Id:         len(emulatorDocs) + 1,
		LastChange: time.Now().Format("20060102 15:04"),
	}

	tests := []struct {
		name           string
		expectedBody   resources.RegistrationsPOSTResponse
		collectionName string
		wantErr        error
		invalidTest    bool
	}{
		// TODO: Add test cases.
		{
			name:           "The new registration has the next id in line",
			expectedBody:   expectedResponse,
			collectionName: resources.RegistrationsCollection,
			wantErr:        nil,
			invalidTest:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			got, err1 := registrations.CreatePOSTResponse(emulatorCtx, emulatorClient, w)
			if err1 != nil && err1.Error() != tt.wantErr.Error() {
				t.Errorf("CreatePOSTResponse() error = %v, \n wantErr = %v", err1, tt.wantErr)
			}
			if got.Id != tt.expectedBody.Id {
				t.Errorf("CreatePOSTResponse() got ID = %v, want ID = %v", got.Id, tt.expectedBody.Id)
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

	invalidJsonMarshal := make(chan bool)

	invalidJsonUnmarshal := `{
			'2222'
		}`

	newDocumentRef, _, err2 := emulatorClient.Collection(resources.RegistrationsCollection).Add(emulatorCtx,
		postRegistration)
	if err2 != nil {
		log.Println("An error occurred when creating a new document:", err2.Error())
	}

	postResponse := resources.RegistrationsPOSTResponse{
		Id:         3,
		LastChange: "20240405 18:07",
	}

	documentID := newDocumentRef.ID

	tests := []struct {
		name         string
		documentID   string
		postResponse any
		expectedCode int
	}{
		// TODO: Add test cases.
		{
			name:         "The POST request body could not be marshaled",
			documentID:   documentID,
			postResponse: invalidJsonMarshal,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "The POST request body could not be unmarshaled",
			documentID:   documentID,
			postResponse: invalidJsonUnmarshal,
			expectedCode: http.StatusInternalServerError,
		},
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
				t.Errorf("Expected HTTP status code %d, but got %d", tt.expectedCode, w.Code)
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
