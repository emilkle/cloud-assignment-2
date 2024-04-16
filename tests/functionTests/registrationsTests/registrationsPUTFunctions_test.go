package registrationsTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"net/http"
	"testing"
)

func TestCreatePUTRequest(t *testing.T) {
	type args struct {
		ctx        context.Context
		client     *firestore.Client
		w          http.ResponseWriter
		data       resources.RegistrationsPOSTandPUT
		documentID string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registrations.CreatePUTRequest(tt.args.ctx, tt.args.client, tt.args.w, tt.args.data, tt.args.documentID)
		})
	}
}

func TestGetDocumentID(t *testing.T) {
	type args struct {
		ctx         context.Context
		client      *firestore.Client
		requestedId string
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
			got, err := registrations.GetDocumentID(tt.args.ctx, tt.args.client, tt.args.requestedId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDocumentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDocumentID() got = %v, expectedBody %v", got, tt.want)
			}
		})
	}
}
