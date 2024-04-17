package webhookTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"countries-dashboard-service/tests/functionTests"
	"reflect"
	"testing"
)

var emulatorClient = functionTests.GetEmulatorClient()
var emulatorCtx = functionTests.GetEmulatorCtx()

var allWebhhoks = []resources.WebhookPOST{
	{
		URL:     "URL1",
		Country: "NO",
		Event:   "POST",
	},
	{
		URL:     "URL2",
		Country: "NO",
		Event:   "POST",
	},
}

var webhook = resources.WebhookPOST{
	URL:     "URL1",
	Country: "NO",
	Event:   "POST",
}

func TestAddWebhook(t *testing.T) {
	functionTests.SetupFirestoreDatabase(resources.WEBHOOK_COLLECTION)
	tests := []struct {
		name         string
		idParam      string
		expectedBody resources.WebhookPOST
	}{
		// TODO: Add test cases.
		{
			name:         "Create a single registration",
			idParam:      "aasflksjdfglksjdf",
			expectedBody: webhook,
		},
		{
			name:         "Registration was not found",
			idParam:      "3",
			expectedBody: resources.WebhookPOST{},
		},
		{
			name:         "Invalid id string",
			idParam:      "sdfsddfs",
			expectedBody: resources.WebhookPOST{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := functions.AddWebhook(emulatorCtx, emulatorClient, tt.idParam, tt.expectedBody)

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

var allWebhhoksGET = []resources.WebhookGET{
	{
		ID:      "1",
		URL:     "URL1",
		Country: "NO",
		Event:   "POST",
	},
	{
		ID:      "2",
		URL:     "URL2",
		Country: "NO",
		Event:   "POST",
	},
}

func TestGetAllWebhooks(t *testing.T) {
	functionTests.SetupFirestoreDatabase(resources.WEBHOOK_COLLECTION)

	tests := []struct {
		name         string
		expectedBody []resources.WebhookGET
		wantErr      bool
		invalidTest  bool
	}{
		{
			name:         "Get all documents",
			expectedBody: allWebhhoksGET,
			wantErr:      false,
			invalidTest:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.GetAllWebhooks(emulatorCtx, emulatorClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllWebhooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantErr) {
				t.Errorf("GetAllWebhooks() got = %v, want %v", got, tt.wantErr)
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
