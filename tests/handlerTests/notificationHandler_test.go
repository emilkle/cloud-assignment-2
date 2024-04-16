package handlerTests

import (
	"countries-dashboard-service/handlers"
	"net/http"
	"testing"
)

func TestCallUrl(t *testing.T) {
	type args struct {
		url     string
		method  string
		content string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.CallUrl(tt.args.url, tt.args.method, tt.args.content)
		})
	}
}

func TestDefaultClientHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.DefaultClientHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestDefaultServerHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.DefaultServerHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestServiceHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.ServiceHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestWebhookHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.WebhookHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_webhookRequestDELETE(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//handlers.webhookRequestDELETE(tt.args.w, tt.args.r)
		})
	}
}

func Test_webhookRequestGET(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//handlers.webhookRequestGET(tt.args.w, tt.args.r)
		})
	}
}

func Test_webhookRequestPOST(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//handlers.webhookRequestPOST(tt.args.w, tt.args.r)
		})
	}
}
