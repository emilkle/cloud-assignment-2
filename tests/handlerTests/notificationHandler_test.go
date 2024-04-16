package handlerTests

/*
import (
	"countries-dashboard-service/handlers"
	"net/http"
	"testing"
)

var allWebhooks = `[
	{
		"ID": "1",
		"URL": "url1",
		"Country": "NO",
		"Event": "POST",
	},
	{
		"ID": "2",
		"URL": "url2",
		"Country": "EN",
		"Event": "POST",
	},
	{
		"ID": "3",
		"URL": "url3",
		"Country": "FI",
		"Event": "POST",
	}
]`

var singleWebhook = `[
	{
		"ID": "1",
		"URL": "url1",
		"Country": "NO",
		"Event": "POST",
	}
]`

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
		{
			name: "Successful POST request",
			args: args{url: "Someurl", method: http.MethodPost, content: "Some string content"},
		},
		{
			name: "Succ"
,,	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.CallUrl(tt.args.url, tt.args.method, tt.args.content)
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

func Test_webhookRequestDELETE(t *testing.T)
} {
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


*/
