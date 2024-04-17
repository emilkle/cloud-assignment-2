package registrationsTests

import (
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePUTRequest(t *testing.T) {
	SetupFirestoreDatabase()

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

	invalidErrorTargetCurrencies := errors.New("element:of 'targetCurrencies' field is not a string, " +
		"or the array is not a string array")

	lastChangeError := errors.New("code = InvalidArgument desc = Document name " +
		"\"projects/countries-dashboard-service/databases/(default)/documents/Registrations/" +
		"\" has invalid trailing \"/\".")

	tests := []struct {
		name         string
		data         resources.RegistrationsPOSTandPUT
		documentID   string
		wantErr      error
		expectedCode int
	}{
		// TODO: Add test cases.
		{
			name:         "The POST request body has a targetCurrencies field that is nil",
			data:         invalidRegistrationStructCurrenciesNil,
			documentID:   "",
			wantErr:      invalidErrorTargetCurrencies,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "The PUT request body has an invalid targetCurrency field where one" +
				" or more of the currencies are empty strings",
			data:         invalidRegistrationStructCurrenciesEmptyStrings,
			documentID:   "",
			wantErr:      invalidErrorTargetCurrencies,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "The last change timestamp could not be updated",
			data:         testRegistration,
			documentID:   "",
			wantErr:      lastChangeError,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			registrations.CreatePUTRequest(emulatorCtx, emulatorClient, w, tt.data, tt.documentID)
			if w.Code != tt.expectedCode {
				t.Errorf("Expected HTTP status code %d, but got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestGetDocumentID(t *testing.T) {
	SetupFirestoreDatabase()

	tests := []struct {
		name        string
		requestedId string
		want        string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "The document with the requested id was found, and the length of " +
				"the random document id is correct",
			requestedId: "1",
			want:        "FxObvU0Wpr2A1L9MT99z",
			wantErr:     false,
		},
		{
			name:        "The document with the requested id was not found",
			requestedId: "5",
			want:        "",
			wantErr:     true,
		},
		{
			name:        "The requested id could not be parsed to integer",
			requestedId: "k",
			want:        "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registrations.GetDocumentID(emulatorCtx, emulatorClient, tt.requestedId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDocumentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetDocumentID() got = %v, expected %v", len(got), len(tt.want))
			}
		})
	}
}
