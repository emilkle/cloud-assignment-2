package dashboardTests

import (
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRetrieveCoordinates(t *testing.T) {
	// Create a local test server to simulate a successful HTTP response i.e. the mockJSONResponse.
	// A HTTP request will respond the mockJSONResponse
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `{
			"results": [
				{
					"latitude": 12.34,
					"longitude": 56.78
				}
			]
		}`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer tsSuccess.Close()

	// Create a local test server to simulate internal server error
	tsServerError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer tsServerError.Close()

	// Define your test cases
	tests := []struct {
		name           string
		country        string
		id             int
		runTest        bool
		testServer     *httptest.Server
		expectedResult resources.CoordinatesValues
		expectedError  string
	}{
		{
			name:       "Successful Retrieval",
			country:    "TestCountry",
			id:         1,
			runTest:    true,
			testServer: tsSuccess,
			expectedResult: resources.CoordinatesValues{
				Latitude:  12.34,
				Longitude: 56.78,
			},
			expectedError: "",
		},
		{
			name:          "Server Error Response",
			country:       "TestCountry",
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
				dashboards.TestUrlRetrieveCoordinates = tt.testServer.URL
			}

			// Call RetrieveCoordinates function using the values from the test cases
			result, err := dashboards.RetrieveCoordinates(tt.country, tt.id, tt.runTest)

			// Check for an unexpected error
			if (err != nil) != (tt.expectedError != "") {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				return
			}

			// Check if the expected error message matches the actual error
			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("Expected error message '%s', got '%s'", tt.expectedError, err.Error())
				return
			}

			// Check if the actual result matches the expected result
			if tt.expectedError == "" && (result.Latitude != tt.expectedResult.Latitude || result.Longitude != tt.expectedResult.Longitude) {
				t.Errorf("Expected result to be %+v, got %+v", tt.expectedResult, result)
			}
		})
	}
}
