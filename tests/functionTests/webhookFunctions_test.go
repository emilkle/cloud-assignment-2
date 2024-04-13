package functionTests

import (
	"context"
	"countries-dashboard-service/functions"
	"countries-dashboard-service/resources"
	"reflect"
	"testing"
)

func TestAddWebhook(t *testing.T) {
	type args struct {
		webhookID string
		data      resources.WebhookPOST
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := functions.AddWebhook(tt.args.webhookID, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AddWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteWebhook(t *testing.T) {
	type args struct {
		ctx      context.Context
		client   *firestore.Client
		structID string
	}
	tests := []struct {
		name    string
		args    args
		want    *resources.WebhookPOSTResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.DeleteWebhook(tt.args.ctx, tt.args.client, tt.args.structID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteWebhook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := functions.GenerateID(); got != tt.want {
				t.Errorf("GenerateID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllWebhooks(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *firestore.Client
	}
	tests := []struct {
		name    string
		args    args
		want    []resources.WebhookGET
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.GetAllWebhooks(tt.args.ctx, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllWebhooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllWebhooks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWebhook(t *testing.T) {
	type args struct {
		ctx       context.Context
		client    *firestore.Client
		webhookID string
	}
	tests := []struct {
		name    string
		args    args
		want    *resources.WebhookGET
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := functions.GetWebhook(tt.args.ctx, tt.args.client, tt.args.webhookID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWebhook() got = %v, want %v", got, tt.want)
			}
		})
	}
}