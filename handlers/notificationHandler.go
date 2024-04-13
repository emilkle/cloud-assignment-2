package handlers

import (
	"bytes"
	"countries-dashboard-service/database"
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	id := functions.GenerateID()
	err = functions.AddWebhook(id, webhook)
	if err != nil {
		log.Println("handle the error", err)
	}

	// Generate response from id
	response := resources.WebhookPOSTResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
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
		webhookResponses, err1 := functions.GetAllWebhooks(ctx, client)
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
	webhookResponse, err := functions.GetWebhook(ctx, client, id)
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

	response, err := functions.DeleteWebhook(ctx, client, id)
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

/*
Invokes the web service to trigger event. Currently only responds to POST requests.
*/
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := database.GetFirestoreContext()
	client := database.GetFirestoreClient()

	str, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error during decoding message content. Error: " + string(str))
	}

	switch r.Method {
	case http.MethodPost:
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Fatal("Error decoding request body: ", err)
		}

		urlFromRequest, ok := data["url"].(string)
		if !ok {
			log.Fatal("Missing 'url' field in request body")
		}

		var webhooks, err2 = functions.GetAllWebhooks(ctx, client)
		if err2 != nil {
			log.Println("Error fetching all webhooks: ", err2)
			http.Error(w, "Error fetching webhooks", http.StatusInternalServerError)
			return
		}

		var matchingWebhook *resources.WebhookGET
		for _, v := range webhooks {
			if v.URL == urlFromRequest {
				matchingWebhook = &v
				break
			}
		}

		if matchingWebhook == nil {
			log.Println("No matching webhook found for URL: ", urlFromRequest)
			http.Error(w, "No matching webhook found", http.StatusNotFound)
			return
		}
		go CallUrl(matchingWebhook.URL, matchingWebhook.Event, string(str), matchingWebhook, w)

		/*
			log.Println("Received POST request...")
			// Iterate through registered webhooks and invoke based on registered URL, method, and with received content
			var webhooks, _ = functions.GetAllWebhooks(ctx, client)
			if err != nil {
				log.Println("TODO Handle this error later", err)
			}
			for _, v := range webhooks {
				log.Println("Trigger event: Call to service endpoint with method " + v.Event +
					" and content '" + string(str) + "'.")

			}

		*/
	default:
		// Edit in the correct constant when done implementing webhooks
		http.Error(w, "Method "+r.Method+" not supported for "+resources.WebhookInv, http.StatusMethodNotAllowed)
	}
}

/*
Calls given URL with given content and awaits response (status and body).
*/
func CallUrl(url string, method string, content string, webhook *resources.WebhookGET, w io.Writer) {
	log.Println("Attempting invocation of url " + url + " with content '" + content + "'.")
	//res, err := http.Post(url, "text/plain", bytes.NewReader([]byte(content)))
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(content)))
	if err != nil {
		log.Printf("%v", "Error during request creation. Error:", err)
		return
	}

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

	if res.StatusCode == http.StatusOK {
		country := webhook.Country
		id := webhook.ID

		response := resources.WebhookInvocation{
			ID:      id,
			Country: country,
			Event:   webhook.Event,
			Time:    time.Now().Format("20060102 04:05"),
		}

		err3 := json.NewEncoder(w).Encode(response)
		if err3 != nil {
			log.Println("Error encoding response: ", err3)
		}
	} else {
		// Handle unsuccessful response (e.g., log error, return a different response object)
		log.Println("Webhook invocation unsuccessful. Received status code:", res.StatusCode)

	}

	// Read the response
	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response. Error:", err)
		return
	}

	log.Println("Webhook " + url + " invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}

/*
Default handler for server side displaying service information.
*/
func DefaultServerHandler(w http.ResponseWriter, r *http.Request) {

	// Define content type, so browser renders links correctly
	w.Header().Add("content-type", "text/html")

	// Prepare output returned to client
	output := "This service offers the following endpoints: <br>The " + resources.NOTIFICATIONS_PATH + " endpoint provides the registration functionality for webhooks (" + http.MethodPost + " and " + http.MethodGet + "), <br>The " +
		resources.WebhookInv + " endpoint triggers the invocation of registered webhooks when called (with arbitrary payload)." +
		"<br>The payload structure for the webhook registration via " +
		http.MethodPost + " is the following JSON structure: {\"url\": \"http://targetHost:targetPort/pathTobeInvoked\", \"event\": \"POST\"}<br>" +
		"Please see the associated Readme for more information."

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}

/*
Default handler for client side displaying service information.
*/
func DefaultClientHandler(w http.ResponseWriter, r *http.Request) {

	// Define content type, so browser renders links correctly
	w.Header().Add("content-type", "text/html")

	// Prepare output returned to client
	output := "This service reacts on the following endpoint: <br>The /pathToBeInvoked endpoint reacts to invocation with any payload, " +
		"but can variably be configured to perform integrity checks (assessible via console output only)" +
		"<br>Please see the associated Readme for more information."

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
