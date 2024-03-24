package functionTests

import (
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"reflect"
	"testing"
)

func TestCreateRegistrationsGET(t *testing.T) {
	correctCurrencies := []string{"EUR", "USD", "SEK"}
	tests := []struct {
		name    string
		idParam string
		want    resources.RegistrationsGET
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Positive test", "1", resources.RegistrationsGET{
			Id:      1,
			Country: "Norway",
			IsoCode: "NO",
			Features: resources.Features{
				Temperature:      true,
				Precipitation:    true,
				Capital:          true,
				Coordinates:      true,
				Population:       true,
				Area:             false,
				TargetCurrencies: correctCurrencies,
			},
			LastChange: "20240229 14:07",
		},
			false},
		{"Negative test", "0", resources.RegistrationsGET{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.CreateRegistrationsGET(tt.idParam)
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
