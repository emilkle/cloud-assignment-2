package handlerTests

import (
	"bytes"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"countries-dashboard-service/tests/functionTests"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var allWebhooks = `[
	{
		"ID": "1",
		"URL": "url1",
		"Country": "NO",
		"Event": "POST",
	},
	{
		"ID": "2",
		"URL": "url2",
		"Country": "EN",
		"Event": "POST",
	},
	{
		"ID": "3",
		"URL": "url3",
		"Country": "FI",
		"Event": "POST",
	}
]`

var singleWebhook = `[
	{
		"ID": "1",
		"URL": "url1",
		"Country": "NO",
		"Event": "POST",
	}
]`

var emulatorClient = functionTests.GetEmulatorClient()
var emulatorCtx = functionTests.GetEmulatorCtx()

func Test_webhookTrigger(t *testing.T) {
	firestoreEmulator.InitializeFirestoreEmulator()
	firestoreEmulator.PopulateFirestoreWithWebhooks()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	iter := emulatorClient.Collection(resources.WEBHOOK_COLLECTION).Documents(emulatorCtx)

	for {
		doc, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Fatalf("Failed to iterate over documents: %v", err1)
			return
		}
		_, err1 = doc.Ref.Delete(emulatorCtx)
		if err1 != nil {
			log.Printf("Failed to delete document: %v", err1)
		}
	}
	firestoreEmulator.PopulateFirestoreWithWebhooks()

	type args struct {
		httpMethod string
		w          http.ResponseWriter
		r          *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.WebhookTrigger(tt.args.httpMethod, tt.args.w, tt.args.r)
		})
	}
}

func TestCallUrl(t *testing.T) {
	type args struct {
		url     string
		id      string
		content string
		event   string
		country string
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			handlers.CallUrl(tt.args.url, tt.args.id, tt.args.content, tt.args.event, tt.args.country, w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("CallUrl() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestWebhookHandler(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		w.Header().Set("content-type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
		case http.MethodPut:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "REST method '"+r.Method+"' is not supported. Try"+
				" '"+http.MethodGet+", "+http.MethodPost+", "+http.MethodPut+" "+
				""+"or"+" "+http.MethodDelete+"' instead. ", http.StatusNotImplemented)
			return
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		method       string
		server       *httptest.Server
		expectedCode int
	}{
		// TODO: Add test cases.
		{"Method = GET (Status OK)", http.MethodGet, mockServer,
			http.StatusOK},
		{"Method = POST (Status OK)", http.MethodPost, mockServer,
			http.StatusOK},
		{"Method = PUT (Status Not allowed)", http.MethodPut, mockServer,
			http.StatusMethodNotAllowed},
		{"Method = DELETE (Status OK)", http.MethodDelete, mockServer,
			http.StatusOK},
		{"Method = OPTIONS (Status not implemented)", http.MethodOptions, mockServer,
			http.StatusNotImplemented},
		{"Method = HEAD (Status not implemented)", http.MethodHead, mockServer,
			http.StatusNotImplemented},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, mockServer.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("handler returned wrong status code for method %s: got %v want %v",
					tt.method, resp.StatusCode, tt.expectedCode)
			}
		})
	}
}

func Test_webhookRequestDELETE(t *testing.T) {
	notFoundIds := []string{"5", "6", "7"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		id := urlParts[4]

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		for _, notFoundId := range notFoundIds {
			if notFoundId == id {
				w.WriteHeader(http.StatusNotFound)
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
	}{
		{
			name:         "Delete webhook",
			method:       http.MethodDelete,
			path:         resources.REGISTRATIONS_PATH + "3",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Trying to delete an registration that is not registered",
			method:       http.MethodDelete,
			path:         resources.REGISTRATIONS_PATH + "6",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "No id given",
			method:       http.MethodDelete,
			path:         resources.REGISTRATIONS_PATH + "",
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, mockServer.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}

func Test_webhookRequestGET(t *testing.T) {
	notFoundIds := []string{"5", "6", "7"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		id := urlParts[4]

		// If trying to get all documents, no id is specified.
		if id == "" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, allWebhooks)
		}

		for _, notFoundId := range notFoundIds {
			if notFoundId == id {
				w.WriteHeader(http.StatusNotFound)
			}
		}

		if id == "3" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, singleWebhook)
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		expectedBody string
		expectedCode int
	}{
		// TODO: Add test cases.
		{
			name:         "Get all webhooks",
			method:       http.MethodGet,
			path:         resources.REGISTRATIONS_PATH,
			expectedBody: allWebhooks,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get webhook",
			method:       http.MethodGet,
			path:         resources.REGISTRATIONS_PATH + "3",
			expectedBody: singleWebhook,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get invalid registration",
			method:       http.MethodGet,
			path:         resources.REGISTRATIONS_PATH + "5",
			expectedBody: "",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err1 := http.NewRequest(tt.method, mockServer.URL+tt.path, nil)
			if err1 != nil {
				t.Fatalf("failed to create request: %v", err1)
			}

			resp, err2 := http.DefaultClient.Do(req)
			if err2 != nil {
				t.Fatalf("failed to send request: %v", err2)
			}

			body, err3 := io.ReadAll(resp.Body)
			if err3 != nil {
				t.Fatalf("failed to read response body: %v", err3)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("Expected response body %q but got %q.", tt.expectedBody, string(body))
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}

func Test_webhookRequestPOST(t *testing.T) {
	postRequestBody := map[string]interface{}{
		"ID":      "qwerty",
		"URL":     "someURL",
		"Country": "NO",
		"Event":   "POST",
	}

	postRequestBodyBytes, err1 := json.Marshal(postRequestBody)
	if err1 != nil {
		t.Fatalf("failed to marshal payload: %v", err1)
	}

	invalidRequestBody := map[string]interface{}{
		"ID":      "ola",
		"URL":     "otherURL",
		"Country": "SE",
		"Event":   "POST",
	}

	invalidRequestBodyBytes, err1 := json.Marshal(invalidRequestBody)
	if err1 != nil {
		t.Fatalf("failed to marshal payload: %v", err1)
	}

	postResponseBody := `{
		"id": "qwerty",
	}`

	postResponseBodyBytes := bytes.NewBufferString(postResponseBody)

	invalidResponse := "Invalid request body format"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var receivedPostRequestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&receivedPostRequestBody); err != nil {
			w.WriteHeader(http.StatusForbidden)
			t.Fatalf(fmt.Sprintf(resources.DECODING_ERROR+"of the POST request. Use this structure for your"+
				" POST request instead: \n%s", resources.JSON_STRUCT_POST_AND_PUT))
		}

		if reflect.DeepEqual(receivedPostRequestBody, postRequestBody) {
			w.WriteHeader(http.StatusCreated)
			w.Write(postResponseBodyBytes.Bytes())
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid request body format"))
		}

	}))
	defer mockServer.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		payload      []byte
		expectedBody string
		expectedCode int
	}{
		{
			name:         "Add a new webhook",
			method:       http.MethodPost,
			path:         resources.NOTIFICATIONS_PATH,
			payload:      postRequestBodyBytes,
			expectedBody: postResponseBody,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Wrong request format",
			method:       http.MethodPost,
			path:         resources.NOTIFICATIONS_PATH,
			payload:      invalidRequestBodyBytes,
			expectedBody: invalidResponse,
			expectedCode: http.StatusForbidden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Post(mockServer.URL, "application/json", bytes.NewBuffer(tt.payload))
			if err != nil {
				t.Fatalf("failed to send POST request: %v", err)
			}
			defer resp.Body.Close()

			var responseBodyBuffer bytes.Buffer
			_, err = responseBodyBuffer.ReadFrom(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			responseBody := responseBodyBuffer.String()

			if responseBody != tt.expectedBody {
				t.Errorf("Expected response: %s but got %s",
					tt.expectedBody, responseBody)
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d but got %d.", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}
