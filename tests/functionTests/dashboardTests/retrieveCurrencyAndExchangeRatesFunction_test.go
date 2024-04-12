package dashboardTests

import (
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRetrieveCurrencyExchangeRates tests the RetrieveCurrencyExchangeRates function without making a real HTTP request to a rest API.
func TestRetrieveCurrencyExchangeRates(t *testing.T) {
	// Create a local test server to simulate a successful HTTP response i.e. the mockJSONResponse.
	// A HTTP request will respond the mockJSONResponse
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `{
			"rates": {
				"USD": 0.90,
				"EUR": 0.95
			}
		}`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer tsSuccess.Close()

	// Create a local test server to simulate internal server error
	tsServerError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer tsServerError.Close()

	// Struct for different test cases
	tests := []struct {
		name           string                         // Name of the test
		id             int                            // Parameter needed by the RetrieveCurrencyExchangeRates function
		runTest        bool                           // Flag variable for identifying whether the RetrieveCurrencyExchangeRates function is to be tested or not
		testServer     *httptest.Server               // Reference to the test server
		expectedResult resources.TargetCurrencyValues // The expected result returned by RetrieveCurrencyExchangeRates function from the test
		expectedError  string                         // The expected error from the test
	}{
		{
			name:       "Successful Retrieval",
			id:         1,
			runTest:    true,
			testServer: tsSuccess,
			expectedResult: resources.TargetCurrencyValues{
				TargetCurrencies: map[string]float64{"USD": 0.90, "EUR": 0.95},
			},
			expectedError: "",
		},
		{
			name:          "Server Error Response",
			id:            2,
			runTest:       true,
			testServer:    tsServerError,
			expectedError: "HTTP error: 500 Internal Server Error",
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the test URL if the runTest flag is true to use the mock HTTP server
			if tt.runTest {
				dashboards.TestUrlRetrieveCurrencyExchangeRates = tt.testServer.URL
			}

			// Call RetrieveCurrencyExchangeRates function using the values from the test cases
			result, err := dashboards.RetrieveCurrencyExchangeRates(tt.id, tt.runTest)

			// Check for an unexpected error
			if (err != nil) != (tt.expectedError != "") {
				t.Errorf("%s: Expected error %v, got %v", tt.name, tt.expectedError, err)
				return
			}

			// Check if the expected error message matches the actual error
			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("%s: Expected error message '%s', got '%s'", tt.name, tt.expectedError, err.Error())
				return
			}

			// Check if the actual result matches the expected result
			if tt.expectedError == "" && !equalCurrencyValues(result, tt.expectedResult) {
				t.Errorf("%s: Expected result to be %+v, got %+v", tt.name, tt.expectedResult, result)
			}
		})
	}
}

// Helper function to compare two TargetCurrencyValues maps
func equalCurrencyValues(actual, expected resources.TargetCurrencyValues) bool {
	if len(actual.TargetCurrencies) != len(expected.TargetCurrencies) {
		return false
	}
	for key, actualValue := range actual.TargetCurrencies {
		expVal, ok := expected.TargetCurrencies[key]
		if !ok || actualValue != expVal {
			return false
		}
	}
	return true
}
