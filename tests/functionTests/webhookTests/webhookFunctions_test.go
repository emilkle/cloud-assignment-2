package webhookTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"countries-dashboard-service/tests/functionTests"
	"google.golang.org/api/iterator"
	"log"
	"reflect"
	"testing"
)

var emulatorClient = functionTests.GetEmulatorClient()
var emulatorCtx = functionTests.GetEmulatorCtx()

func TestAddWebhook(t *testing.T) {
	firestoreEmulator.InitializeFirestoreEmulator()
	firestoreEmulator.PopulateFirestoreWithWebhooks()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	iter := emulatorClient.Collection(resources.WEBHOOK_COLLECTION).Documents(emulatorCtx)

	for {
		doc, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Failed to iterate over documents: %v", err1)
			return
		}
		_, err1 = doc.Ref.Delete(emulatorCtx)
		if err1 != nil {
			log.Printf("Failed to delete document: %v", err1)
		}
	}
	firestoreEmulator.PopulateFirestoreWithWebhooks()

	tests := []struct {
		name             string
		idParam          string
		inputBody        resources.WebhookPOST
		expectedWebhook  *resources.WebhookGET
		expectedResponse *resources.WebhookGET
	}{
		{
			name:    "Create a single registration",
			idParam: "3",
			inputBody: resources.WebhookPOST{
				URL:     "URL3",
				Country: "NO",
				Event:   "POST",
			},
			expectedResponse: &resources.WebhookGET{
				ID:      "3",
				URL:     "URL3",
				Country: "NO",
				Event:   "POST",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := functions.AddWebhook(emulatorCtx, emulatorClient, tt.idParam, tt.inputBody)
			if err != nil {
				return
			}
			expectedWebhook, err := functions.GetWebhook(emulatorCtx, emulatorClient, "3")
			if !reflect.DeepEqual(expectedWebhook, tt.expectedResponse) {
				t.Errorf("AddWebhook() got = %v, expectedBody %v", expectedWebhook, tt.expectedResponse)
			}
		})
	}
}

func TestDeleteWebhook(t *testing.T) {
	firestoreEmulator.InitializeFirestoreEmulator()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	iter := emulatorClient.Collection(resources.WEBHOOK_COLLECTION).Documents(emulatorCtx)

	for {
		doc, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Failed to iterate over documents: %v", err1)
			return
		}
		_, err1 = doc.Ref.Delete(emulatorCtx)
		if err1 != nil {
			log.Printf("Failed to delete document: %v", err1)
		}
	}
	firestoreEmulator.PopulateFirestoreWithWebhooks()

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
		{
			name: "Successful deletion of webhook",
			args: args{
				ctx:      emulatorCtx,
				client:   emulatorClient,
				structID: "1",
			},
			want: &resources.WebhookPOSTResponse{
				ID: "1",
			},
		},
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
	firestoreEmulator.InitializeFirestoreEmulator()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	iter := emulatorClient.Collection(resources.WEBHOOK_COLLECTION).Documents(emulatorCtx)

	for {
		doc, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Failed to iterate over documents: %v", err1)
			return
		}
		_, err1 = doc.Ref.Delete(emulatorCtx)
		if err1 != nil {
			log.Printf("Failed to delete document: %v", err1)
		}
	}
	firestoreEmulator.PopulateFirestoreWithWebhooks()

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
			if !reflect.DeepEqual(got, tt.expectedBody) {
				t.Errorf("GetAllWebhooks() got = %v, want %v", got, tt.expectedBody)
			}
		})
	}
}

func TestGetWebhook(t *testing.T) {
	firestoreEmulator.InitializeFirestoreEmulator()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	iter := emulatorClient.Collection(resources.WEBHOOK_COLLECTION).Documents(emulatorCtx)

	for {
		doc, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Failed to iterate over documents: %v", err1)
			return
		}
		_, err1 = doc.Ref.Delete(emulatorCtx)
		if err1 != nil {
			log.Printf("Failed to delete document: %v", err1)
		}
	}
	firestoreEmulator.PopulateFirestoreWithWebhooks()

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
		{
			name: "Webhook get success",
			args: args{
				ctx:       emulatorCtx,
				client:    emulatorClient,
				webhookID: "1",
			},
			want: &resources.WebhookGET{
				ID:      "1",
				URL:     "URL1",
				Country: "NO",
				Event:   "POST",
			},
			wantErr: false,
		},
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
