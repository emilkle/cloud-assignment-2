package functions

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	uuid2 "github.com/google/uuid"
	"google.golang.org/api/iterator"
	"log"
)

// Generates a unique ID using google uuid library
func GenerateID() string {
	uuid, err := uuid2.NewRandom()
	if err != nil {
		// If something goes wrong returns a placeholder ID
		return "PlaceholderID"
	}
	return uuid.String()
}

// Adds a webhook to the firestore database
func AddWebhook(ctx context.Context, client *firestore.Client, webhookID string, data resources.WebhookPOST) error {
	ref := client.Collection(resources.WEBHOOK_COLLECTION)

	// Encode data struct as JSON byte slice
	jsonData, err1 := json.Marshal(data)
	if err1 != nil {
		return fmt.Errorf("error marshalling webhook data: %v", err1)
	}

	var webhook resources.WebhookGET
	err2 := json.Unmarshal(jsonData, &webhook)
	if err2 != nil {
		return fmt.Errorf("error unmarshalling webhook data: %v", err2)
	}

	// Create a new document with a unique ID
	_, _, err3 := ref.Add(ctx, map[string]interface{}{
		"ID":      webhookID,
		"URL":     webhook.URL,
		"Country": webhook.Country,
		"Event":   webhook.Event,
	})
	if err3 != nil {
		return fmt.Errorf("error adding webhook: %v", err3)
	}

	log.Printf("Webhook with ID %s successfully added", webhook.ID)
	return nil
}

// Deletes a webhook from the database
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

// GetAllWebhooks gets all webhooks from database
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

// Fetches a single webhook from the database
func GetWebhook(ctx context.Context, client *firestore.Client, webhookID string) (*resources.WebhookGET, error) {
	ref := client.Collection(resources.WEBHOOK_COLLECTION)
	query := ref.Where("ID", "==", webhookID).Limit(1)

	docs, err1 := query.Documents(ctx).GetAll()
	if err1 != nil {
		return nil, err1
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("document with webhook ID %s not found", webhookID)
	}

	// Extract data from the document
	var webhookData resources.WebhookGET
	for _, doc := range docs {
		if err2 := doc.DataTo(&webhookData); err2 != nil {
			return nil, err2
		}
		// Since you're using Limit(1), you only expect one document, but just in case, break after processing one.
		break
	}
	return &webhookData, nil
}
