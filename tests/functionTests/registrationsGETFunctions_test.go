package functionTests

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/functions/registrations"
	"countries-dashboard-service/resources"
	"reflect"
	"testing"
)

var validData = map[string]interface{}{
	"id":      1,
	"country": "Norway",
	"isoCode": "NO",
	"features": map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             true,
		"targetCurrencies": []interface{}{"NOK", "USD"},
	},
	"lastChange": "20220101 15:07",
}

var invalidData = map[string]interface{}{
	"id":      1,
	"country": "Norway",
	"isoCode": "NO",
	"features": map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             true,
		"targetCurrencies": []string{"NOK", "USD"},
	},
	"lastChange": "20220101 15:07",
}

var want = resources.RegistrationsGET{
	Id:      1,
	Country: "Norway",
	IsoCode: "NO",
	Features: resources.Features{
		Temperature:      true,
		Precipitation:    true,
		Capital:          true,
		Coordinates:      true,
		Population:       true,
		Area:             true,
		TargetCurrencies: []string{"NOK", "USD"},
	},
	LastChange: "20220101 15:07",
}

var doNotWant = resources.RegistrationsGET{
	Id:      1,
	Country: "Norway",
	IsoCode: "NO",
	Features: resources.Features{
		Temperature:      true,
		Precipitation:    true,
		Capital:          true,
		Coordinates:      true,
		Population:       true,
		Area:             true,
		TargetCurrencies: nil,
	},
	LastChange: "20220101 15:07",
}

func TestCreateRegistrationsGET(t *testing.T) {
	tests := []struct {
		name    string
		idParam string
		want    resources.RegistrationsGET
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registrations.CreateRegistrationsGET(tt.idParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRegistrationsGET() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRegistrationsGET() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllRegisteredDocuments(t *testing.T) {
	tests := []struct {
		name    string
		want    []resources.RegistrationsGET
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registrations.GetAllRegisteredDocuments()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRegisteredDocuments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllRegisteredDocuments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTargetCurrencies(t *testing.T) {
	var currencies = []string{"NOK", "USD"}

	var featuresData = map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             true,
		"targetCurrencies": []interface{}{"NOK", "USD"},
	}

	var invalidFeaturesData = map[string]interface{}{
		"temperature":      true,
		"precipitation":    true,
		"capital":          true,
		"coordinates":      true,
		"population":       true,
		"area":             true,
		"targetCurrencies": []string{"NOK", "USD"},
	}

	tests := []struct {
		name         string
		featuresData map[string]interface{}
		want         []string
		wantErr      bool
	}{
		// TODO: Add test cases.
		{name: "Returns string array", featuresData: featuresData, want: currencies, wantErr: false},
		{name: "Returns error", featuresData: invalidFeaturesData, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.GetTargetCurrencies(tt.featuresData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTargetCurrencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CreateRegistrationsResponse(t *testing.T) {
	tests := []struct {
		name       string
		data       map[string]interface{}
		lastChange string
		idIndex    int
		want       resources.RegistrationsGET
		wantErr    bool
	}{
		// TODO: Add test cases.
		{name: "Positive test", data: validData, lastChange: "20220101 15:07", idIndex: 1,
			want: want, wantErr: false},
		{name: "Negative test", data: invalidData, lastChange: "20220101 15:07", idIndex: 1,
			want: doNotWant, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registrations.CreateRegistrationsResponse(tt.data, tt.lastChange, tt.idIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createRegistrationsResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UpdateId(t *testing.T) {
	type args struct {
		ctx         context.Context
		client      *firestore.Client
		documentID  string
		getResponse resources.RegistrationsGET
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registrations.UpdateId(tt.args.ctx, tt.args.client, tt.args.documentID, tt.args.getResponse)
		})
	}
}
