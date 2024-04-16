package registrationsTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/functions/registrations"
	"reflect"
	"testing"
)

func TestDeleteDocumentWithRequestedId(t *testing.T) {
	type args struct {
		ctx          context.Context
		client       *firestore.Client
		requestedIds []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.DeleteDocumentWithRequestedId(tt.args.ctx, tt.args.client, tt.args.requestedIds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteDocumentWithRequestedId() = %v, expectedBody %v", got, tt.want)
			}
		})
	}
}

func TestFindDocumentWithId(t *testing.T) {
	type args struct {
		ctx        context.Context
		client     *firestore.Client
		documentId int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.FindDocumentWithId(tt.args.ctx, tt.args.client, tt.args.documentId); got != tt.want {
				t.Errorf("FindDocumentWithId() = %v, expectedBody %v", got, tt.want)
			}
		})
	}
}
