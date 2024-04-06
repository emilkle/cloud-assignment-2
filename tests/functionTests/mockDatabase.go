package functionTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"github.com/stretchr/testify/mock"
	"log"
)

// FirestoreClient defines the methods used by Firestore client.
type FirestoreClient interface {
	Set(ctx context.Context, docRef *firestore.DocumentRef, data interface{},
		opts ...firestore.SetOption) (*firestore.WriteResult, error)
	Collection(path string) *firestore.CollectionRef
	// Add other methods used by Firestore client if necessary
}

// MockFirestoreClient is a mock implementation of FirestoreClient for testing.
type MockFirestoreClient struct {
	SetFunc func(ctx context.Context, docRef *firestore.DocumentRef, data interface{},
		opts ...firestore.SetOption) (*firestore.WriteResult, error)
	mock.Mock
	// Add other mock functions if necessary
}

func (m *MockFirestoreClient) Collection(path string) *firestore.CollectionRef {
	//TODO implement me
	//panic("implement me")
	return &firestore.CollectionRef{
		Path: path,
	}
}

// Set implements the Set method for the mock Firestore client.
func (m *MockFirestoreClient) Set(ctx context.Context, docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	return m.SetFunc(ctx, docRef, data, opts...)
}

// UpdateId updates the document's id field.
func UpdateId(ctx context.Context, client FirestoreClient, documentID string, getResponse resources.RegistrationsGET) {
	// Update the document's id field.
	_, err := client.Set(ctx,
		client.Collection(resources.REGISTRATIONS_COLLECTION).Doc(documentID),
		map[string]interface{}{"id": getResponse.Id},
		firestore.MergeAll)

	if err != nil {
		log.Println("The id field could not be set: ", err.Error())
	}
}
