package handlers

import (
	"bytes"
	"countries-dashboard-service/database"
	"countries-dashboard-service/functions/notifications"
	"countries-dashboard-service/resources"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Initialize signature (via init())
var SignatureKey = "X-SIGNATURE"

// var Mac hash.Hash
var Secret []byte

// WebhookHandler handles webhook registration (POST), lookup (GET) requests and deletion (DELETE) requests.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		webhookRequestPOST(w, r)
		//webhookTrigger(w, r)
	case http.MethodGet:
		webhookRequestGET(w, r)
	case http.MethodDelete:
		webhookRequestDELETE(w, r)
	default:
		http.Error(w, "Method "+r.Method+" not supported for "+resources.NOTIFICATIONS_PATH, http.StatusMethodNotAllowed)
	}
}

// webhookRequestPOST handles the HTTP POST request for webhooks to be stored in the Firestore database.
func webhookRequestPOST(w http.ResponseWriter, r *http.Request) {
	webhook := resources.WebhookPOST{}
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		http.Error(w, "Something went wrong during decoding: "+err.Error(), http.StatusBadRequest)
	}

	// Generate ID and assign it to webhook struct
	id := notifications.GenerateID()
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()
	err = notifications.AddWebhook(ctx, client, id, webhook)
	if err != nil {
		log.Println("handle the error", err)
	}

	// Generate response from id
	response := resources.WebhookPOSTResponse{ID: id}

	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Encode response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Something went during encoding: "+err.Error(), http.StatusBadRequest)
	}
	log.Println("Webhook with url " + webhook.URL + " and ID " + id + " has been registered.")
}

// webhookRequestGET handles the HTTP GET request for webhooks stored in the Firestore database.
// It is possible to get all documents at once by calling /dashboard/v1/notifications/ .
// For getting specific entries /dashboard/v1/registrations/{id} is used.
func webhookRequestGET(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Retrieve the 4th url-part that contains the id.
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	// Check if the query does not contain an id.
	if id == "" {
		// Fetch all the documents in the  firestore database and handle the error that it returns.
		webhookResponses, err1 := notifications.GetAllWebhooks(ctx, client)
		if err1 != nil {
			http.Error(w, "Could not retrieve all documents.", http.StatusInternalServerError)
			return
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(webhookResponses)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	//var webhookResponses []resources.WebhookGET
	//var notFoundIds []string
	webhookResponse, err := notifications.GetWebhook(ctx, client, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(webhookResponse)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
	}
}

// webhookRequestDELETE handles HTTP DELETE requests for deleting webhooks.
// It retrieves the webhook id from the URL, deletes the corresponding document from the database,
// and returns a response indicating success or failure.
func webhookRequestDELETE(w http.ResponseWriter, r *http.Request) {
	client := database.GetFirestoreClient()
	ctx := database.GetFirestoreContext()

	// Retrieve the 4th url-part that contains the id.
	urlParts := strings.Split(r.URL.Path, "/")
	id := urlParts[4]

	// Check if the query does not contain an id.
	if id == "" {
		log.Println("No id(s) were specified in this query.")
		http.Error(w, "No id(s) were specified in this query, please write an "+
			"integer number in the query to use this service.", http.StatusBadRequest)
		return
	}

	response, err := notifications.DeleteWebhook(ctx, client, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error deleting webhook: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
	}
	log.Println("Webhook with id: " + id + " has been deleted from the database.")
	http.Error(w, "The requested webhook were successfully deleted from the database.", http.StatusNoContent)

}

// CallUrl sends an HTTP POST request to the specified URL with a payload containing
// information about the webhook invocation. It also performs content-based validation
// by generating a signature header based on the content and including it in the request.
func CallUrl(url string, id string, content string, event, country string, w io.Writer) {
	// Prepare payload containing webhook information
	payload := map[string]interface{}{
		"ID":      id,
		"Country": country,
		"Event":   event,
		"time":    time.Now().Format(time.RFC3339),
	}

	// Marshal payload to json
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling payload: ", err)
	}

	// Create request using method POST
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Error creating request: ", err)
	}

	req.Header.Set("Content-Type", "application/json")

	/// BEGIN: HEADER GENERATION FOR CONTENT-BASED VALIDATION

	// Hash content (for content-based validation; not relevant for URL-based validation)
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(content))
	if err != nil {
		log.Printf("%v", "Error during content hashing. Error:", err)
		return
	}
	// Convert hash to string & add to header to transport to client for validation
	req.Header.Add(SignatureKey, hex.EncodeToString(mac.Sum(nil)))

	/// END: CONTENT-BASED VALIDATION

	// Perform invocation
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request. Error:", err)
		return
	}

	// Read response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
	}

	// Log invocation of webhook
	log.Println("Webhook with url " + url + " and ID " + id + " invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(body))
}

// WebhookTrigger is an HTTP handler function responsible for triggering webhooks based on the incoming HTTP request.
// It fetches all registered webhooks from the database, identifies the ones with matching URLs, and triggers them asynchronously.
func WebhookTrigger(httpMethod string, w http.ResponseWriter, r *http.Request) {
	// Establish firestore context and client
	ctx := database.GetFirestoreContext()
	client := database.GetFirestoreClient()

	// Extract url from incoming HTTP request
	endpointUrl := r.URL.String()
	urlFromRequest := resources.RootPath + endpointUrl
	//urlFromRequest = "https://webhook.site/20d8180f-b4d4-479e-9aa6-32d970dd21ae" // Delete this line when done developing

	// Fetch all webhooks from database
	var webhooks, err = notifications.GetAllWebhooks(ctx, client)
	if err != nil {
		log.Println("Error fetching all webhooks: ", err)
		http.Error(w, "Error fetching webhooks", http.StatusInternalServerError)
		return
	}

	// Filter webhooks to find the ones with matching urls
	var matchingWebhooks []*resources.WebhookGET
	for _, v := range webhooks {
		if v.URL == urlFromRequest {
			webhookCopy := v
			matchingWebhooks = append(matchingWebhooks, &webhookCopy)
		}
	}

	// Return error if no matching webhooks are found
	if matchingWebhooks == nil {
		log.Println("No matching webhook found for URL: ", urlFromRequest)
		http.Error(w, "No matching webhook found", http.StatusNotFound)
		return
	}

	method := ""
	switch httpMethod {
	case http.MethodGet:
		method = resources.EventInvoke
	case http.MethodPut:
		method = resources.EventChange
	case http.MethodPost:
		method = resources.EventRegister
	case http.MethodDelete:
		method = resources.EventDelete
	}

	// Invoke the matching webhook(s)
	for _, webhook := range matchingWebhooks {
		go CallUrl(webhook.URL, webhook.ID, webhook.URL, method, webhook.Country, w)
	}

}
