package dashboardTests

import (
	"countries-dashboard-service/functions/dashboards"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// String made to Mock the JSON response from the rest api http:// 129.241.150.113:9090/ currency/NOK
var mockJSONHttpRequestResponse = `{
    "result": "success",
    "provider": "https://www.exchangerate-api.com",
    "documentation": "https://www.exchangerate-api.com/docs/free",
    "terms_of_use": "https://www.exchangerate-api.com/terms",
    "time_last_update_unix": 1712707351,
    "time_last_update_utc": "Wed, 10 Apr 2024 00:02:31 +0000",
    "time_next_update_unix": 1712794861,
    "time_next_update_utc": "Thu, 11 Apr 2024 00:21:01 +0000",
    "time_eol_unix": 0,
    "base_code": "NOK",
    "rates": {
        "NOK": 1,
        "USD": 0.093687,
        "EUR": 0.086289
    }
}`

// TestHttpRequest tests the HttpRequest function without making a real HTTP request to a rest API.
func TestHttpRequest(t *testing.T) {
	// Create a local test server to simulate a successful HTTP response i.e. the mockJSONResponse.
	// A HTTP request will respond the mockJSONResponse
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockJSONHttpRequestResponse)
	}))
	defer tsSuccess.Close()

	// Create a local test server to simulate not found error.
	// A HTTP request will produce a 404 status code.
	tsNotFound := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer tsNotFound.Close()

	// Struct for different test cases
	tests := []struct {
		name     string // Name of the test
		url      string // URL for the test, which in this case points to one of the test servers
		fetching string // variable needed by the HttpRequest function
		id       int    // Variable needed by the HttpRequest function
		wantErr  bool   // Indicate if a error is expected
		wantCode int    // Expected HTTP status code
	}{
		{
			name:     "HttpRequestSuccessful",
			url:      tsSuccess.URL,
			fetching: "test fetching success",
			id:       1,
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name:     "NotFoundResponse",
			url:      tsNotFound.URL,
			fetching: "testing 404 response",
			id:       1,
			wantErr:  false,
			wantCode: http.StatusNotFound,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call HttpRequest function using the values from the test cases
			got, err := dashboards.HttpRequest(tt.url, tt.fetching, tt.id)

			// Checks for an unexpected error
			if err != nil {
				t.Errorf("%s unexpected error: %v", tt.name, err)
				return
			}

			// Check if the expected status code is the same as the returned status code
			if got.StatusCode != tt.wantCode {
				t.Errorf("%s expected status code %d, got %d", tt.name, tt.wantCode, got.StatusCode)
			}
		})
	}
}
