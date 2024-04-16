package registrationsTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/functions/registrations"
	"reflect"
	"testing"
)

func TestDeleteDocumentWithRequestedId(t *testing.T) {
	firestoreEmulator.PopulateFirestoreData()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	tests := []struct {
		name         string
		requestedIds []string
		notFoundIds  []string
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "All document successfully deleted",
			requestedIds: []string{"1", "2"},
			notFoundIds:  nil,
			wantErr:      false,
		},
		{
			name:         "Multiple documents could not be found",
			requestedIds: []string{"3", "4"},
			notFoundIds:  []string{"3", "4"},
			wantErr:      true,
		},
		{
			name:         "One document could not be found",
			requestedIds: []string{"3"},
			notFoundIds:  []string{"3"},
			wantErr:      true,
		},
		{
			name:         "Invalid id requested",
			requestedIds: []string{"test1", "test2"},
			notFoundIds:  nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.DeleteDocumentWithRequestedId(emulatorCtx, emulatorClient, tt.requestedIds); !reflect.DeepEqual(got, tt.notFoundIds) {
				t.Errorf("DeleteDocumentWithRequestedId() = %v, expected %v", got, tt.notFoundIds)
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
