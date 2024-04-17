package functions

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"fmt"
	uuid2 "github.com/google/uuid"
	"google.golang.org/api/iterator"
	"log"
)

// GenerateID Generates a unique ID using google uuid library and returns it
func GenerateID() string {
	uuid, err := uuid2.NewRandom()
	if err != nil {
		// If something goes wrong returns a placeholder ID
		return "PlaceholderID"
	}
	return uuid.String()
}

// AddWebhook adds a webhook to the firestore database
// It takes a context, Firestore client, the ID of the webhook to add and a struct of type WebhookPOST.
// It returns a log of successful addition of webhook and an error.
func AddWebhook(ctx context.Context, client *firestore.Client, webhookID string, data resources.WebhookPOST) error {
	ref := client.Collection(resources.WEBHOOK_COLLECTION)

	// Create a new document with a unique ID
	_, _, err3 := ref.Add(ctx, map[string]interface{}{
		"ID":      webhookID,
		"URL":     data.URL,
		"Country": data.Country,
		"Event":   data.Event,
	})
	if err3 != nil {
		return fmt.Errorf("error adding webhook: %v", err3)
	}
	return nil
}

// DeleteWebhook removes a webhook from the database
// It takes a context, Firestore client, and the ID of the webhook to delete.
// It returns a log of successful deletion of webhook or an error.
func DeleteWebhook(ctx context.Context, client *firestore.Client, structID string) (*resources.WebhookPOSTResponse, error) {
	// Reference the webhooks collection in firestore and query document with corresponding id
	ref := client.Collection(resources.WEBHOOK_COLLECTION)
	query := ref.Where("ID", "==", structID).Limit(1)

	// Get all documents matching the query
	docs, err1 := query.Documents(ctx).GetAll()
	if err1 != nil {
		return nil, err1
	}

	// Check if any documents were found
	if len(docs) == 0 {
		return nil, fmt.Errorf("document with webhook ID %s not found", structID)
	}

	// Delete the first document
	_, err2 := docs[0].Ref.Delete(ctx)
	if err2 != nil {
		return nil, err2
	}

	log.Printf("Document with webhook ID %v successfully deleted\n", structID)
	return &resources.WebhookPOSTResponse{ID: structID}, nil
}

// GetAllWebhooks retrieves all webhooks from database
// It takes a context and a Firestore client.
// It returns an array of webhooks, or an error if not found or any other error occurs.
func GetAllWebhooks(ctx context.Context, client *firestore.Client) ([]resources.WebhookGET, error) {
	// Iterate over documents in ascending order of lastChange timestamp.
	iter := client.Collection(resources.WEBHOOK_COLLECTION).Documents(ctx)
	var webhookResponses []resources.WebhookGET
	idIndex := 1

	// Iterate through documents and construct RegistrationsGET structs.
	for {
		document, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			return nil, err1
		}
		data := document.Data()

		// Construct RegistrationsGET response struct.
		webhookResponse := resources.WebhookGET{
			ID:      data["ID"].(string),
			URL:     data["URL"].(string),
			Country: data["Country"].(string),
			Event:   data["Event"].(string),
		}

		webhookResponses = append(webhookResponses, webhookResponse)
		idIndex++
	}

	return webhookResponses, nil
}

// GetWebhook retrieves a webhook document from Firestore based on its ID.
// It takes a context, Firestore client, and the ID of the webhook to retrieve.
// It returns the webhook data if found, or an error if not found or any other error occurs.
func GetWebhook(ctx context.Context, client *firestore.Client, webhookID string) (*resources.WebhookGET, error) {
	// Make reference to webhook collection in firestore database and query doc with specified id
	ref := client.Collection(resources.WEBHOOK_COLLECTION)
	query := ref.Where("ID", "==", webhookID).Limit(1)

	// Retrieve documents matching query
	docs, err1 := query.Documents(ctx).GetAll()
	if err1 != nil {
		return nil, err1
	}

	// If not found return an error
	if len(docs) == 0 {
		return nil, fmt.Errorf("document with webhook ID %s not found", webhookID)
	}

	// Extract data from the document
	var webhookData resources.WebhookGET
	for _, doc := range docs {
		if err2 := doc.DataTo(&webhookData); err2 != nil {
			return nil, err2
		}

		break
	}
	return &webhookData, nil
}
