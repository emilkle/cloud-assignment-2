package dashboardTests

import (
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRetrieveCapitalPopulationAndArea tests the RetrieveCapitalPopulationAndArea function without making a real HTTP request to a rest API.
func TestRetrieveCapitalPopulationAndArea(t *testing.T) {
	// Create a local test server to simulate a successful HTTP response i.e. the mockJSONResponse.
	// A HTTP request will respond the mockJSONResponse
	tsResponse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `[{
			"capital": ["Test Capital"],
			"population": 12345,
			"area": 6789.0
		}]`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer tsResponse.Close()

	// Create a local test server to simulate internal server error
	tsServerError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer tsServerError.Close()

	// Struct for different test cases
	tests := []struct {
		name           string                          // Name of the test
		isoCode        string                          // Variable needed by the RetrieveCapitalPopulationAndArea function
		id             int                             // Variable needed by the RetrieveCapitalPopulationAndArea function
		runTest        bool                            // Flag variable for identifying whether the RetrieveCapitalPopulationAndArea is tested or not
		testServer     *httptest.Server                // Reference to the test server
		expectedResult resources.CapitalPopulationArea // The expected result returned by RetrieveCapitalPopulationAndArea from the test
		expectedError  string                          // The expected error from the test
	}{
		{
			name:       "Successful Retrieval",
			isoCode:    "TEST",
			id:         1,
			runTest:    true,
			testServer: tsResponse,
			expectedResult: resources.CapitalPopulationArea{
				Capital:    []string{"Test Capital"},
				Population: 12345,
				Area:       6789.0,
			},
			expectedError: "",
		},
		{
			name:          "Server Error Response",
			isoCode:       "TEST",
			id:            4,
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
				dashboards.TestUrl = tt.testServer.URL
			}

			// Call RetrieveCapitalPopulationAndArea function using the values from the test cases
			result, err := dashboards.RetrieveCapitalPopulationAndArea(tt.isoCode, tt.id, tt.runTest)

			// Checks for an unexpected error
			if (err != nil) != (tt.expectedError != "") {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				return
			}

			// Checks if expected error is equal to actual error
			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("Expected error message '%s', got '%s'", tt.expectedError, err.Error())
				return
			}

			// Checks if actual result matches the expected result
			if tt.expectedError == "" && !isEqualResources(result, tt.expectedResult) {
				t.Errorf("Expected result to be %+v, got %+v", tt.expectedResult, result)
			}
		})
	}
}

// isEqualResources compares two instances of resources.CapitalPopulationArea
func isEqualResources(actual, expected resources.CapitalPopulationArea) bool {
	return len(actual.Capital) == len(expected.Capital) && actual.Capital[0] == expected.Capital[0] && actual.Population == expected.Population && actual.Area == expected.Area
}
