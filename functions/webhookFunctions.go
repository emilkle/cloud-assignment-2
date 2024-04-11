package functions

import (
	"countries-dashboard-service/database"
	"countries-dashboard-service/resources"
	"encoding/json"
	"fmt"
	uuid2 "github.com/google/uuid"
	"log"
)

// Generates a unique ID using google uuid library
func generateID() string {
	uuid, err := uuid2.NewRandom()
	if err != nil {
		// If something goes wrong returns a placeholder ID
		return "PlaceholderID"
	}
	return uuid.String()
}

// Adds a webhook to the firestore database
func addWebhook(data []byte) error {
	ctx := database.GetFirestoreContext()
	client := database.GetFirestoreClient()

	ref := client.Collection(resources.WEBHOOK_COLLECTION)

	var webhook resources.WebhookGET
	err := json.Unmarshal(data, &webhook)
	if err != nil {
		return fmt.Errorf("error unmarshalling webhook data: %v", err)
	}

	// Generate ID and assign it to webhook struct
	id := generateID()
	webhook.ID = id

	// Create a new document with a unique ID
	_, _, err = ref.Add(ctx, map[string]interface{}{
		"ID":      webhook.ID,
		"URL":     webhook.URL,
		"Country": webhook.Country,
		"Event":   webhook.Event,
	})
	if err != nil {
		return fmt.Errorf("error adding webhook: %v", err)
	}

	log.Printf("Webhook with ID %s successfully added", webhook.ID)
	return nil
}

func deleteWebhook(structID string) error {
	ctx := database.GetFirestoreContext()
	client := database.GetFirestoreClient()

	ref := client.Collection(resources.WEBHOOK_COLLECTION)
	query := ref.Where("ID", "==", structID).Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		return fmt.Errorf("Document with webhook ID %s not found", structID)
	}

	_, err = docs[0].Ref.Delete(ctx)
	if err != nil {
		return err
	}

	log.Println("Document with webhook ID %s successfully deleted", structID)
	return nil
}

func getAllWebhooks() []resources.WebhookGET {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	ref := client.Collection(resources.WEBHOOK_COLLECTION)

	docs, err := ref.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error fetching documents: %s\n", err)
	}

	var webhooks []resources.WebhookGET

	for _, doc := range docs {
		var webhook resources.WebhookGET
		if err := doc.DataTo(&webhook); err != nil {
			log.Printf("Error parsing %s", err)
			return nil
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks
}
