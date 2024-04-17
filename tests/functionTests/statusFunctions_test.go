package functionTests

import (
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/functions"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckEndpointStatus(t *testing.T) {
	// Struct for test cases
	tests := []struct {
		name               string // Name of the test
		mockResponseStatus int    // A mocked status response returned by the test server
		want               int    // The status that is wanted to be returned in the test
		shutdownServer     bool   // Determine whether to simulate a network error by shutting down the server
	}{
		{
			name:               "StatusOK",
			mockResponseStatus: http.StatusOK,
			want:               http.StatusOK,
		},
		{
			name:               "StatusInternalServerError",
			mockResponseStatus: http.StatusInternalServerError,
			want:               http.StatusInternalServerError,
		},
		{
			name:               "NetworkError",
			mockResponseStatus: http.StatusOK, // This value will be ignored because the server will be shut down
			want:               http.StatusServiceUnavailable,
			shutdownServer:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creates a mock server for testing
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockResponseStatus)
			}))

			// Shut down server to simulate network error
			if tt.shutdownServer {
				server.Close()
			}

			// Call the CheckEndpointStatus function
			got := functions.CheckEndpointStatus(server.URL)

			// Check if the retrieved status is different from the wanted response
			if got != tt.want {
				t.Errorf("%s: CheckEndpointStatus() = %v, want %v", tt.name, got, tt.want)
			}

			if !tt.shutdownServer {
				server.Close()
			}
		})
	}
}

// TestNumberOfRegisteredWebhooksGet Checks if the function retrieves the correct number of webhooks from firestore
func TestNumberOfRegisteredWebhooksGet(t *testing.T) {
	firestoreEmulator.InitializeFirestoreEmulator()
	firestoreEmulator.PopulateFirestoreWithWebhooks()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	tests := []struct {
		name string
		want int
	}{
		{
			name: "Retrieve Webhook Count",
			want: 2,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := functions.NumberOfRegisteredWebhooksGet(emulatorClient, emulatorCtx); got != tt.want {
				t.Errorf("NumberOfRegisteredWebhooksGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
