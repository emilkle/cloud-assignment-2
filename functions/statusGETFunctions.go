package functions

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"fmt"
	"net/http"
	"time"
)

var client *firestore.Client
var ctx context.Context

// CheckEndpointStatus checks and returns the status of an endpoint.
// If the endpoint does not respond within 10 seconds it is timed out
func CheckEndpointStatus(url string) int {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Get(url)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return http.StatusServiceUnavailable
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()
	return response.StatusCode
}

// NumberOfRegisteredWebhooksGet fetches all webhooks stored in the webhook collection in the database
// and returns the number of webhooks
func NumberOfRegisteredWebhooksGet() int {
	client, ctx = dashboards.RecognizeEnvironmentVariableForClientContext(client, ctx)

	collection := client.Collection(resources.WEBHOOK_COLLECTION)
	webhooks, err := collection.Documents(ctx).GetAll()
	if err != nil {
		fmt.Printf("Failed to get all webhooks: %v", err)
	}
	return len(webhooks)
}

// CheckFirestoreStatus returns an HTTP status code that simulates the Firestore database's status.
func CheckFirestoreStatus() int {
	client, ctx = dashboards.RecognizeEnvironmentVariableForClientContext(client, ctx)
	collection := client.Collection(resources.WEBHOOK_COLLECTION)
	_, err := collection.Limit(1).Documents(ctx).GetAll()
	if err != nil {
		fmt.Printf("Error accessing Firestore: %v\n", err)
		// Simulate HTTP 503 Service Unavailable for any Firestore errors
		return http.StatusServiceUnavailable
	}
	// Simulate HTTP 200 OK if no errors were returned
	return http.StatusOK
}
