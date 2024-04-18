package handlerTests

import (
	"countries-dashboard-service/handlers"
	"countries-dashboard-service/resources"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var expectedDashboard = resources.DashboardsGet{
	Country: "Norway",
	IsoCode: "NO",
	FeatureValues: resources.FeatureValues{
		Temperature:   2.0,
		Precipitation: 1.0,
		Capital:       "Oslo",
		Coordinates: resources.CoordinatesValues{
			Latitude:  62.0,
			Longitude: 10.0,
		},
		Population: 5379475,
		Area:       385180.0,
		TargetCurrencies: map[string]float64{
			"EUR": 0.086312,
			"USD": 0.998935,
			"SEK": 0.091928,
		},
	},
	LastRetrieval: "20240229 14:07",
}

// TestDashboardsHandler tests the DashboardsHandler
func TestDashboardsHandler(t *testing.T) {
	handlers.SkipRealCallOfRetrieveDashboardGet = true

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, resources.DashboardsPath)

		w.Header().Set("Content-Type", "application/json")

		switch id {
		case "1":
			w.WriteHeader(http.StatusOK)
			response, err := json.Marshal(expectedDashboard)
			if err != nil {
				t.Fatal("Failed to marshal expected response:", err)
			}
			w.Write(response)
		case "":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Dashboard not found"}`))
		case "1,2":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"Cannot retrieve more than one dashboard, too many IDs specified."}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Dashboard not found"}`))
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expected       resources.DashboardsGet // Compare the actual structure
	}{
		{
			name:           "GET Dashboard",
			method:         http.MethodGet,
			path:           resources.DashboardsPath + "1",
			expectedStatus: http.StatusOK,
			expected:       expectedDashboard,
		},
		{
			name:           "Invalid Method",
			method:         http.MethodPost,
			path:           resources.DashboardsPath + "1",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Multiple IDs Error",
			method:         http.MethodGet,
			path:           resources.DashboardsPath + "1,2",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Dashboard Not Found",
			method:         http.MethodGet,
			path:           resources.DashboardsPath + "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(tt.method, mockServer.URL+tt.path, nil)
			responseRecorder := httptest.NewRecorder()
			handlers.DashboardsHandler(responseRecorder, request)

			if responseRecorder.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var got resources.DashboardsGet
				if err := json.NewDecoder(responseRecorder.Body).Decode(&got); err != nil {
					t.Fatal("Failed to decode response:", err)
				}
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Unexpected response: got %+v, want %+v", got, tt.expected)
				}
			}
		})
	}
}
