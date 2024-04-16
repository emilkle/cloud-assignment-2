package registrationsTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"fmt"
	"log"
	"reflect"
	"testing"
)

var emulatorClient *firestore.Client
var emulatorCtx context.Context

var validData = map[string]interface{}{
	"id":      1,
	"country": "Norway",
	"isoCode": "NO",
	"features": map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             false,
		"targetCurrencies": []interface{}{"EUR", "USD", "SEK"},
	},
	"lastChange": "20240229 14:07",
}

var invalidData = map[string]interface{}{
	"id":      1,
	"country": "Norway",
	"isoCode": "NO",
	"features": map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             true,
		"targetCurrencies": []string{"NOK", "USD"},
	},
	"lastChange": "20240229 14:07",
}

var want = resources.RegistrationsGET{
	Id:      1,
	Country: "Norway",
	IsoCode: "NO",
	Features: resources.Features{
		Temperature:      true,
		Precipitation:    true,
		Capital:          true,
		Coordinates:      true,
		Population:       true,
		Area:             false,
		TargetCurrencies: []string{"EUR", "USD", "SEK"},
	},
	LastChange: "20240229 14:07",
}

var doNotWant = resources.RegistrationsGET{
	Id:      1,
	Country: "Norway",
	IsoCode: "NO",
	Features: resources.Features{
		Temperature:      true,
		Precipitation:    true,
		Capital:          true,
		Coordinates:      true,
		Population:       true,
		Area:             true,
		TargetCurrencies: nil,
	},
	LastChange: "20240229 14:07",
}

func TestCreateRegistrationsGET(t *testing.T) {
	firestoreEmulator.PopulateFirestoreData()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	tests := []struct {
		name         string
		idParam      string
		expectedBody resources.RegistrationsGET
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "Create a single registration",
			idParam:      "1",
			expectedBody: want,
			wantErr:      false,
		},
		{
			name:         "Registration was not found",
			idParam:      "3",
			expectedBody: resources.RegistrationsGET{},
			wantErr:      true,
		},
		{
			name:         "Invalid id string",
			idParam:      "sdfsddfs",
			expectedBody: resources.RegistrationsGET{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err1 := registrations.CreateRegistrationsGET(emulatorCtx, emulatorClient, tt.idParam)
			if (err1 != nil) != tt.wantErr {
				t.Errorf("Could not find the document with id: " + tt.idParam)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedBody) {
				t.Errorf("GetAllRegisteredDocuments() got = %v, expectedBody %v", got, tt.expectedBody)
			}
		})
	}
}

func TestGetAllRegisteredDocuments(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		client       *firestore.Client
		expectedBody []resources.RegistrationsGET
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registrations.GetAllRegisteredDocuments(tt.ctx, tt.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRegisteredDocuments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedBody) {
				t.Errorf("GetAllRegisteredDocuments() got = %v, expectedBody %v", got, tt.expectedBody)
			}
		})
	}
}

func TestGetTargetCurrencies(t *testing.T) {
	var currencies = []string{"NOK", "USD"}

	var featuresData = map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             true,
		"targetCurrencies": []interface{}{"NOK", "USD"},
	}

	var invalidFeaturesData = map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             true,
		"targetCurrencies": []string{"NOK", "USD"},
	}

	tests := []struct {
		name         string
		featuresData map[string]interface{}
		want         []string
		wantErr      bool
	}{
		// TODO: Add test cases.
		{name: "Returns string array", featuresData: featuresData, want: currencies, wantErr: false},
		{name: "Returns error", featuresData: invalidFeaturesData, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.GetTargetCurrencies(tt.featuresData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTargetCurrencies() = %v, expectedBody %v", got, tt.want)
			}
		})
	}
}

func Test_CreateRegistrationsResponse(t *testing.T) {
	tests := []struct {
		name       string
		data       map[string]interface{}
		lastChange string
		idIndex    int
		want       resources.RegistrationsGET
		wantErr    bool
	}{
		// TODO: Add test cases.
		{name: "Positive test", data: validData, lastChange: "20240229 14:07", idIndex: 1,
			want: want, wantErr: false},
		{name: "Negative test", data: invalidData, lastChange: "20240229 14:07", idIndex: 1,
			want: doNotWant, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.CreateRegistrationsResponse(tt.data, tt.lastChange, tt.idIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createRegistrationsResponse() = %v, expectedBody %v", got, tt.want)
			}
		})
	}
}

func Test_UpdateId(t *testing.T) {
	firestoreEmulator.PopulateFirestoreData()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	firstDocumentID := ""
	docRef := emulatorClient.Collection(resources.REGISTRATIONS_COLLECTION).Limit(1)
	docs, err := docRef.Documents(emulatorCtx).GetAll()
	if err != nil {
		log.Println("Failed to retrieve documents: ", err.Error())
		return
	}

	if len(docs) > 0 {
		firstDocumentID = docs[0].Ref.ID
		fmt.Println("First document ID:", firstDocumentID)
	} else {
		fmt.Println("No documents found in the collection")
	}

	tests := []struct {
		name        string
		documentID  string
		getResponse resources.RegistrationsGET
		expectedId  int
		expectedErr bool
	}{
		// TODO: Add test cases.
		{
			name:        "Valid document id",
			documentID:  firstDocumentID,
			getResponse: want,
			expectedId:  1,
			expectedErr: false,
		},
		{
			name:        "Document id is invalid",
			documentID:  "",
			getResponse: resources.RegistrationsGET{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registrations.UpdateId(emulatorCtx, emulatorClient, tt.documentID, tt.getResponse)
			mockResponse := resources.RegistrationsGET{
				Id: tt.getResponse.Id,
			}
			if tt.expectedId != mockResponse.Id {
				t.Errorf("Expected id %v, but got %v", tt.expectedId, mockResponse.Id)
			}
		})
	}
}
