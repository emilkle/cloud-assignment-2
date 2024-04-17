package dashboardTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/firestoreEmulator"
	"countries-dashboard-service/functions/dashboards"
	"countries-dashboard-service/resources"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

var emulatorClient *firestore.Client
var emulatorCtx context.Context

var testRegistrationGet = resources.RegistrationsGET{
	Id:      1,
	Country: "Norway",
	Features: resources.Features{
		Area:             false,
		Capital:          true,
		Coordinates:      true,
		Population:       true,
		Precipitation:    true,
		TargetCurrencies: []string{"EUR", "USD", "SEK"},
		Temperature:      true,
	},
	LastChange: "20240229 14:07",
}

var testFeatures = resources.Features{
	Area:             false,
	Capital:          true,
	Coordinates:      true,
	Population:       true,
	Precipitation:    true,
	TargetCurrencies: []string{"EUR", "USD", "SEK"},
	Temperature:      true,
}

var testFeaturesData = map[string]interface{}{
	"area":             false,
	"capital":          true,
	"coordinates":      true,
	"population":       true,
	"precipitation":    true,
	"temperature":      true,
	"targetCurrencies": []interface{}{"EUR", "USD", "SEK"}, // Ensure this matches expected type
}

var testDashboardData = []map[string]interface{}{
	{
		"country": "Norway",
		"features": map[string]interface{}{
			"area":             false,
			"capital":          true,
			"coordinates":      true,
			"population":       true,
			"precipitation":    true,
			"targetCurrencies": []interface{}{"EUR", "USD", "SEK"},
			"temperature":      true,
		},
		"id":         int64(1),
		"isoCode":    "NO",
		"lastChange": "20240229 14:07",
	},
}

var validDashboard = map[string]interface{}{
	"country": "Norway",
	"isoCode": "NO",
	"features": map[string]interface{}{
		"temperature":   2.0,
		"precipitation": 1.0,
		"capital":       "Oslo",
		"coordinates": map[string]interface{}{
			"latitude":  62.0,
			"longitude": 10.0,
		},
		"population": 5379475,
		"area":       385180.0,
		"targetCurrencies": map[string]float64{
			"NOK": 10.5,
			"USD": 1.1,
		},
	},
	"lastChange": "20240229 14:07",
}

var invalidDashboard = map[string]interface{}{
	"country": "Norway",
	"isoCode": "NO",
	"features": map[string]interface{}{
		"temperature":   2.0,
		"precipitation": 1.0,
		"capital":       "Oslo",
		"coordinates": map[string]interface{}{
			"latitude":  62.0,
			"longitude": 10.0,
		},
		"population": 5379475,
		"area":       385180.0,
		"targetCurrencies": map[string]float64{
			"NOK": 10.5,
			"USD": 1.1,
		},
	},
	"lastChange": "20240229 14:07",
}

var wantedDashboard = resources.DashboardsGet{
	Country: "Norway",
	IsoCode: "NO",
	FeatureValues: resources.FeatureValues{
		Temperature:   5.2,
		Precipitation: 2.0,
		Capital:       "Oslo",
		Coordinates: resources.CoordinatesValues{
			Latitude:  62.0,
			Longitude: 10.0,
		},
		Population: 5379475,
		Area:       0,
		TargetCurrencies: map[string]float64{
			"EUR": 0.086312,
			"USD": 0.998935,
			"SEK": 0.091928,
		},
	},
	LastRetrieval: time.Now().Format("20060102 15:04"),
}

var notWantedDashboard = resources.DashboardsGet{
	Country: "Norway",
	IsoCode: "NO",
	FeatureValues: resources.FeatureValues{
		Temperature:   1.5,
		Precipitation: 0.1,
		Capital:       "Oslo",
		Coordinates: resources.CoordinatesValues{
			Latitude:  62.0,
			Longitude: 13.0,
		},
		Population: 5436758,
		Area:       189086,
		TargetCurrencies: map[string]float64{
			"EUR": 1.0,
			"USD": 1.0,
			"SEK": 1.0,
		},
	},
	LastRetrieval: "20240229 14:07",
}

var testTimeTempAndPrecipitation = resources.ForecastResponse{
	Hourly: resources.HourlyData{
		Time: []string{
			"2024-04-16T00:00",
			"2024-04-16T01:00",
			"2024-04-16T02:00",
			"2024-04-16T03:00",
			"2024-04-16T04:00",
		},
		Temperature: []float64{
			5.5, 6.0, 5.8, 5.2, 3.5,
		},
		Precipitation: []float64{
			1.0, 1.5, 3.0, 1.5, 3,
		},
	},
}

var testHourlyData = resources.HourlyData{
	Time: []string{
		"2024-04-16T00:00",
		"2024-04-16T01:00",
		"2024-04-16T02:00",
		"2024-04-16T03:00",
		"2024-04-16T04:00",
	},
	Temperature: []float64{
		5.5, 6.0, 5.8, 5.2, 3.5,
	},
	Precipitation: []float64{
		1.0, 1.5, 3.0, 1.5, 3,
	},
}

var featureValues = resources.FeatureValues{
	Temperature:   -1.5,
	Precipitation: 0.0,
	Capital:       "Oslo",
	Coordinates: resources.CoordinatesValues{
		Latitude:  62.0,
		Longitude: 10.0,
	},
	Population: 5379475,
	Area:       0,
	TargetCurrencies: map[string]float64{
		"EUR": 0.086312,
		"USD": 0.998935,
		"SEK": 0.091928,
	},
}

var capPopArea = resources.CapitalPopulationArea{
	Capital:    []string{"Oslo"},
	Population: 5379475,
	Area:       0,
}

var coordinateResponse = resources.CoordinatesResponse{
	Results: []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{
		{Latitude: 62.0, Longitude: 10.0},
	},
}

var coordinateValues = resources.CoordinatesValues{
	Latitude:  62.0,
	Longitude: 10.0,
}

var testTargetCurrencyValues = resources.TargetCurrencyValues{
	TargetCurrencies: map[string]float64{
		"EUR": 0.086312,
		"USD": 0.998935,
		"SEK": 0.091928,
	},
}

// TestRetrieveDashboardData tests the RetrieveDashboardData function
func TestRetrieveDashboardData(t *testing.T) {
	firestoreEmulator.PopulateFirestoreWithRegistrations()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	type args struct {
		dashboardId string
	}
	tests := []struct {
		name          string
		args          args
		wantDocument  []map[string]interface{}
		wantIntegerId int
		wantErr       bool
	}{
		{
			name: "TestRetrieveDashboardData_Successful",
			args: args{"1"},
			wantDocument: []map[string]interface{}{
				{
					"country": "Norway",
					"features": map[string]interface{}{
						"area":             false,
						"capital":          true,
						"coordinates":      true,
						"population":       true,
						"precipitation":    true,
						"targetCurrencies": []interface{}{"EUR", "USD", "SEK"},
						"temperature":      true,
					},
					"id":         int64(1),
					"isoCode":    "NO",
					"lastChange": "20240229 14:07",
				},
			},
			wantIntegerId: 1,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDocument, gotIntegerId, err := dashboards.RetrieveDashboardData(emulatorClient, emulatorCtx, tt.args.dashboardId)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveDashboardData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDocumentData(gotDocument), tt.wantDocument) {
				t.Errorf("RetrieveDashboardData() got = %v, want %v", gotDocumentData(gotDocument), tt.wantDocument)
			}
			if gotIntegerId != tt.wantIntegerId {
				t.Errorf("RetrieveDashboardData() got1 = %v, want %v", gotIntegerId, tt.wantIntegerId)
			}
		})
	}
}

// gotDocumentData extracts the data from each document snapshot for direct comparison
func gotDocumentData(docs []*firestore.DocumentSnapshot) []map[string]interface{} {
	var results []map[string]interface{}
	for _, doc := range docs {
		if doc != nil {
			results = append(results, doc.Data())
		}
	}
	return results
}

func TestRetrieveDashboardGet(t *testing.T) {
	firestoreEmulator.PopulateFirestoreWithRegistrations()
	emulatorClient = firestoreEmulator.GetEmulatorClient()
	emulatorCtx = firestoreEmulator.GetEmulatorContext()

	testServerCapPopArea := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `[{
			"capital": ["Oslo"],
			"population": 5379475,
			"area": 323802.0
		}]`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer testServerCapPopArea.Close()

	testServerCurrencyExchange := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `{
        "rates": {
            "EUR": 0.086312,
			"USD": 0.998935,
			"SEK": 0.091928
        }
    }`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer testServerCurrencyExchange.Close()

	testServerTempPercip := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `{
      "hourly": {
        "time": [
          "2024-04-16T00:00", "2024-04-16T01:00", "2024-04-16T02:00",
          "2024-04-16T03:00", "2024-04-16T04:00"
        ],
        "temperature_2m": [
          5.5, 6.0, 5.8, 5.2, 3.5
        ],
        "precipitation": [
          1.0, 1.5, 3.0, 1.5, 3
        ]
      }
    }`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer testServerTempPercip.Close()

	testServerCoordinates := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `{
      "results": [
        {
          "latitude": 62,
          "longitude": 10
        }
      ]
    }`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer testServerCoordinates.Close()

	dashboards.TestUrlRetrieveCapitalPopulationAndArea = testServerCapPopArea.URL
	dashboards.TestUrlRetrieveCurrencyExchangeRates = testServerCurrencyExchange.URL
	dashboards.TestUrlRetrieveTempAndPrecipitation = testServerTempPercip.URL
	dashboards.TestUrlRetrieveCoordinates = testServerCoordinates.URL

	type args struct {
		dashboardId string
	}
	tests := []struct {
		name    string
		args    args
		want    resources.DashboardsGet
		wantErr bool
	}{
		{
			name:    "TestRetrieveDashboardGet_Successful",
			args:    args{"1"},
			want:    wantedDashboard,
			wantErr: false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dashboards.RetrieveDashboardGet(emulatorClient, emulatorCtx, tt.args.dashboardId, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveDashboardGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveDashboardGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetrieveTargetCurrenciesAndExchangeRates(t *testing.T) {
	testServerCurrencyExchange := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockJSONResponse := `{
        "rates": {
            "EUR": 0.086312,
			"USD": 0.998935,
			"SEK": 0.091928
        }
    }`
		fmt.Fprintln(w, mockJSONResponse)
	}))
	defer testServerCurrencyExchange.Close()

	dashboards.TestUrlRetrieveCurrencyExchangeRates = testServerCurrencyExchange.URL

	type args struct {
		featuresData map[string]interface{}
		id           int
	}
	tests := []struct {
		name    string
		args    args
		want    resources.TargetCurrencyValues
		wantErr bool
	}{
		{
			name:    "TestRetrieveTargetCurrenciesAndExchangeRates_Successful",
			args:    args{featuresData: testFeaturesData, id: 1},
			want:    testTargetCurrencyValues,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dashboards.RetrieveTargetCurrenciesAndExchangeRates(tt.args.featuresData, tt.args.id, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveTargetCurrenciesAndExchangeRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveTargetCurrenciesAndExchangeRates() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateMeanTemperatureAndPrecipitation(t *testing.T) {
	type args struct {
		tempAndPrecip resources.HourlyData
	}
	tests := []struct {
		name           string
		args           args
		wantMeanTemp   float64
		wantMeanPrecip float64
	}{
		{
			name:           "TestCalculateMeanTemperatureAndPrecipitation_Successful",
			args:           args{tempAndPrecip: testHourlyData},
			wantMeanTemp:   5.2,
			wantMeanPrecip: 2.0,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := dashboards.CalculateMeanTemperatureAndPrecipitation(tt.args.tempAndPrecip)
			if got != tt.wantMeanTemp {
				t.Errorf("CalculateMeanTemperatureAndPrecipitation() got = %v, want %v", got, tt.wantMeanTemp)
			}
			if got1 != tt.wantMeanPrecip {
				t.Errorf("CalculateMeanTemperatureAndPrecipitation() got = %v, want %v", got1, tt.wantMeanPrecip)
			}
		})
	}
}
