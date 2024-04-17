package registrationsTests

import (
	"cloud.google.com/go/firestore"
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"countries-dashboard-service/tests/functionTests"
	"log"
	"reflect"
	"testing"
)

func TestDeleteDocumentWithRequestedId(t *testing.T) {
	functionTests.SetupFirestoreDatabase(resources.REGISTRATIONS_PATH)

	tests := []struct {
		name         string
		requestedIds []string
		notFoundIds  []string
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "All documents successfully deleted",
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
	functionTests.SetupFirestoreDatabase(resources.REGISTRATIONS_PATH)

	tests := []struct {
		name           string
		documentId     int
		want           bool
		wantErr        bool
		successfulTest bool
	}{
		// TODO: Add test cases.
		{
			name:           "Successfully found and deleted the document",
			documentId:     3,
			want:           true,
			wantErr:        false,
			successfulTest: true,
		},
		{
			name:       "Could not iterate over query result",
			documentId: 0,
			want:       false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		if tt.successfulTest {
			var newDocumentRef *firestore.DocumentRef
			var _ *firestore.WriteResult
			var err1 error
			docs, err2 := emulatorClient.Collection(resources.REGISTRATIONS_COLLECTION).Documents(emulatorCtx).GetAll()
			if err2 != nil {
				log.Println("Failed to retrieve documents: ", err2.Error())
				return
			}

			if len(docs) < 3 {
				newDocumentRef, _, err1 = emulatorClient.Collection(resources.REGISTRATIONS_COLLECTION).
					Add(emulatorCtx, invalidDocument3)
				if err1 != nil {
					log.Printf("An error occurred when creating a new document: %v", err1.Error())
				} else {
					log.Printf(
						"Document added to the registrations collection. Identifier of the added document: %v",
						newDocumentRef.ID)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.FindDocumentWithId(emulatorCtx, emulatorClient, tt.documentId); got != tt.want {
				t.Errorf("FindDocumentWithId() = %v, expected %v", got, tt.want)
			}
		})
	}
}
