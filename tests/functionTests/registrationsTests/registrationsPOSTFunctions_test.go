package registrationsTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	"io"
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

func TestCreatePOSTRequest(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *firestore.Client
		w      http.ResponseWriter
		data   resources.RegistrationsPOSTandPUT
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registrations.CreatePOSTRequest(tt.args.ctx, tt.args.client, tt.args.w, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePOSTRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreatePOSTRequest() got = %v, expectedBody %v", got, tt.want)
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

		jsonResponse, err := json.Marshal(postResponse)
		if err != nil {
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
	type args struct {
		ctx          context.Context
		client       *firestore.Client
		w            http.ResponseWriter
		documentID   string
		postResponse resources.RegistrationsPOSTResponse
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registrations.UpdatePOSTRequest(tt.args.ctx, tt.args.client, tt.args.w, tt.args.documentID, tt.args.postResponse)
		})
	}
}

func TestValidateDataTypes(t *testing.T) {
	type args struct {
		data resources.RegistrationsPOSTandPUT
		w    http.ResponseWriter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registrations.ValidateDataTypes(tt.args.data, tt.args.w); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDataTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
