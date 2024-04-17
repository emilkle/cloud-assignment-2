package registrationsTests

import (
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"net/http/httptest"
	"testing"
)

func TestCreatePUTRequest(t *testing.T) {
	SetupFirestoreDatabase()

	tests := []struct {
		name       string
		data       resources.RegistrationsPOSTandPUT
		documentID string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			registrations.CreatePUTRequest(emulatorCtx, emulatorClient, w, tt.data, tt.documentID)
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
