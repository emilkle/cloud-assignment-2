package webhookTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"net/http"
	"reflect"
	"testing"
)

func TestAddWebhook(t *testing.T) {
	firestoreEmulator.PopulateFirestoreData()
	emulatorClient := firestoreEmulator.GetEmulatorClient()
	emulatorCtx := firestoreEmulator.GetEmulatorContext()

	tests := []struct {
		name         string
		idParam      string
		expectedBody resources.WebhookPOST
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "Create a single registration",
			idParam:      "1",
			expectedBody: resources.WebhookPOST{},
			wantErr:      false,
		},
		{
			name:         "Registration was not found",
			idParam:      "3",
			expectedBody: resources.WebhookPOST{},
			wantErr:      true,
		},
		{
			name:         "Invalid id string",
			idParam:      "sdfsddfs",
			expectedBody: resources.WebhookPOST{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := functions.AddWebhook(emulatorCtx, emulatorClient, tt.idParam, tt.expectedBody)
			/*if (err1 != nil) != tt.wantErr {
				t.Errorf("Could not find the document with id: " + tt.idParam)
				return
			}

			*/
			if !reflect.DeepEqual(got, tt.expectedBody) {
				t.Errorf("GetAllRegisteredDocuments() got = %v, expectedBody %v", got, tt.expectedBody)
			}
		})
	}
}

func TestDeleteWebhook(t *testing.T) {
	type args struct {
		ctx      context.Context
		client   *firestore.Client
		structID string
	}
	tests := []struct {
		name    string
		args    args
		want    *resources.WebhookPOSTResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.DeleteWebhook(tt.args.ctx, tt.args.client, tt.args.structID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteWebhook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name       string
		wantlength int
	}{
		{name: "Successful generation", wantlength: 36},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := functions.GenerateID()
			if len(got) != tt.wantlength {
				t.Errorf("GenerateID() length = %v, want %v", len(got), tt.wantlength)
			}
		})
	}
}

func TestGetAllWebhooks(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *firestore.Client
		w      http.ResponseWriter
		data   resources.WebhookGET
	}
	tests := []struct {
		name    string
		args    args
		want    []resources.WebhookGET
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.GetAllWebhooks(tt.args.ctx, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllWebhooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllWebhooks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWebhook(t *testing.T) {
	type args struct {
		ctx       context.Context
		client    *firestore.Client
		webhookID string
	}
	tests := []struct {
		name    string
		args    args
		want    *resources.WebhookGET
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.GetWebhook(tt.args.ctx, tt.args.client, tt.args.webhookID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWebhook() got = %v, want %v", got, tt.want)
			}
		})
	}
}
