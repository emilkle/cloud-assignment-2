package dashboardTests

import (
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRetrieveTempAndPrecipitation tests the RetrieveTempAndPrecipitation function without making a real HTTP request to a rest API.
func TestRetrieveTempAndPrecipitation(t *testing.T) {
	// Create a local test server to simulate a successful HTTP response i.e. the mockJSONResponse.
	// A HTTP request will respond the mockJSONResponse
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `{
			"hourly": {
				"temperature_2m": [15.5, 15.7, 16.1],
				"precipitation": [0.1, 0.0, 0.2]
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

	// Create a local test server to simulate empty response error
	tsNoData := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"hourly": {"temperature_2m": [], "precipitation": []}}`)
	}))
	defer tsNoData.Close()

	// Struct for the test cases
	tests := []struct {
		name           string               // Name of the tests
		latitude       float64              // Parameter needed by the RetrieveTempAndPrecipitation function
		longitude      float64              // Parameter needed by the RetrieveTempAndPrecipitation function
		id             int                  // Parameter needed by the RetrieveTempAndPrecipitation function
		runTest        bool                 // Flag variable for identifying whether the RetrieveTempAndPrecipitation function is to be tested or not
		testServer     *httptest.Server     // Reference to the test server
		expectedResult resources.HourlyData // The expected result returned by RetrieveTempAndPrecipitation function from the test
		expectedError  string               // The expected error from the test
	}{
		{
			name:       "Successful Retrieval",
			latitude:   50.5,
			longitude:  14.5,
			id:         1,
			runTest:    true,
			testServer: tsSuccess,
			expectedResult: resources.HourlyData{
				Temperature:   []float64{15.5, 15.7, 16.1},
				Precipitation: []float64{0.1, 0.0, 0.2},
			},
			expectedError: "",
		},
		{
			name:          "Server Error Response",
			latitude:      52.0,
			longitude:     13.0,
			id:            2,
			runTest:       true,
			testServer:    tsServerError,
			expectedError: "HTTP error: 500 Internal Server Error",
		},
		{
			name:          "No Data Response",
			latitude:      52.0,
			longitude:     13.0,
			id:            4,
			runTest:       true,
			testServer:    tsNoData,
			expectedError: "no temperature and precipitation data returned",
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the test URL if the runTest flag is true to use the mock HTTP server
			if tt.runTest {
				dashboards.TestUrlRetrieveTempAndPrecipitation = tt.testServer.URL
			}

			// Call RetrieveTempAndPrecipitation function using the values from the test cases
			result, err := dashboards.RetrieveTempAndPrecipitation(tt.latitude, tt.longitude, tt.id, tt.runTest)

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
			if tt.expectedError == "" && !equalHourlyData(result, tt.expectedResult) {
				t.Errorf("%s: Expected result %+v, got %+v", tt.name, tt.expectedResult, result)
			}
		})
	}
}

// equalHourlyData helper function to compare two HourlyData objects
func equalHourlyData(actual, expected resources.HourlyData) bool {
	if len(actual.Temperature) != len(expected.Temperature) || len(actual.Precipitation) != len(expected.Precipitation) {
		return false
	}
	for i, v := range actual.Temperature {
		if v != expected.Temperature[i] {
			return false
		}
	}
	for i, v := range actual.Precipitation {
		if v != expected.Precipitation[i] {
			return false
		}
	}
	return true
}
