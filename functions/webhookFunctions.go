package functions

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/database"
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
func AddWebhook(webhookID string, data resources.WebhookPOST) error {
	ctx := database.GetFirestoreContext()
	client := database.GetFirestoreClient()

	ref := client.Collection(resources.WEBHOOK_COLLECTION)

	// Encode data struct as JSON byte slice
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling webhook data: %v", err)
	}

	var webhook resources.WebhookGET
	err = json.Unmarshal(jsonData, &webhook)
	if err != nil {
		return fmt.Errorf("error unmarshalling webhook data: %v", err)
	}

	// Create a new document with a unique ID
	_, _, err = ref.Add(ctx, map[string]interface{}{
		"ID":      webhookID,
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

// Deletes a webhook from the database
func DeleteWebhook(ctx context.Context, client *firestore.Client, structID string) (*resources.WebhookPOSTResponse, error) {

	ref := client.Collection(resources.WEBHOOK_COLLECTION)
	query := ref.Where("ID", "==", structID).Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("document with webhook ID %s not found", structID)
	}

	_, err = docs[0].Ref.Delete(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Document with webhook ID %v successfully deleted\n", structID)
	return &resources.WebhookPOSTResponse{ID: structID}, nil
}

func CreateWebhookGET(ctx context.Context, client *firestore.Client, idParam string) (resources.WebhookGET, error) {
	// Parse ID parameter to integer.
	//idNumber, err1 := strconv.Atoi(idParam)
	//if err1 != nil {
	//	log.Println("This id could not be parsed, try another id.", err1.Error())
	//	return resources.WebhookGET{}, err1
	//}

	// Query Firestore for documents with matching ID.
	query := client.Collection(resources.WEBHOOK_COLLECTION).Where("id", "==", idParam).Limit(1)
	documents, err2 := query.Documents(ctx).GetAll()
	if err2 != nil {
		log.Println("Failed to fetch documents:", err2)
		return resources.WebhookGET{}, err2
	}

	// Check if any documents were found.
	if len(documents) == 0 {
		err3 := fmt.Errorf("no document found with ID: %s", idParam)
		log.Println(err3)
		return resources.WebhookGET{}, err3
	}

	// Construct RegistrationsGET struct from retrieved data.
	for _, document := range documents {
		data := document.Data()

		return resources.WebhookGET{
			ID:      idParam,
			URL:     data["URL"].(string),
			Country: data["Country"].(string),
			Event:   data["Event"].(string),
		}, nil
	}

	// Print the error message to the server log if the retrieving of the document fails.
	log.Println("Document with ID", idParam, "was not found.")
	return resources.WebhookGET{}, nil
}

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

		// Retrieve the lastChange timestamp from the document.
		/*
			lastChange, ok := data["lastChange"].(string)
			if !ok {
				log.Printf("The timestamp of the last change"+
					" %v could not be converted to string.", data["lastChange"])
				continue
			}

		*/

		// Construct RegistrationsGET response struct.
		webhookResponse := resources.WebhookGET{
			ID:      data["ID"].(string),
			URL:     data["URL"].(string),
			Country: data["Country"].(string),
			Event:   data["Event"].(string),
		}

		//registrationID := document.Ref.ID

		// Update all the id fields in for the Firestore documents after deleting a document in the middle of the
		// ascending order, to ensure that all registration documents will be found.
		//UpdateId(ctx, client, registrationID, registrationsResponse)

		webhookResponses = append(webhookResponses, webhookResponse)

		idIndex++
	}

	return webhookResponses, nil
}

/*
// Finds and returns all webhooks from database
func GetAllWebhooks() ([]resources.WebhookGET, error) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	ref := client.Collection(resources.WEBHOOK_COLLECTION)

	docs, err := ref.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error fetching documents: %s\n", err)
		return nil, err
	}

	var webhooks []resources.WebhookGET

	for _, doc := range docs {
		var webhook resources.WebhookGET
		if err := doc.DataTo(&webhook); err != nil {
			log.Printf("Error parsing %s", err)
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks, nil
}

*/

// Fetches a single webhook from the database
func GetWebhook(ctx context.Context, client *firestore.Client, webhookID string) (*resources.WebhookGET, error) {
	ref := client.Collection(resources.WEBHOOK_COLLECTION)
	query := ref.Where("ID", "==", webhookID).Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("document with webhook ID %s not found", webhookID)
	}

	// Extract data from the document
	var webhookData resources.WebhookGET
	for _, doc := range docs {
		if err := doc.DataTo(&webhookData); err != nil {
			return nil, err
		}
		// Since you're using Limit(1), you only expect one document, but just in case, break after processing one.
		break
	}

	return &webhookData, nil
}
