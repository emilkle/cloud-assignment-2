package dashboards

import (
	"cloud.google.com/go/firestore"
	"context"
	"countries-dashboard-service/database"
	"countries-dashboard-service/firestoreEmulator"
	"io"
	"log"
	"net/http"
	"os"
)

// HttpRequest performs an HTTP GET request to the specified URL
func HttpRequest(url, fetching string, id int) (*http.Response, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch %s for dashboard with id: %d. Error: %s", fetching, id, err)
		return nil, err
	}
	return response, nil
}

// CloseResponseBody closes the response body and logs any errors
func CloseResponseBody(body io.ReadCloser, fetching string, id int) {
	err := body.Close()
	if err != nil {
		log.Printf("failed to close response body while fetching %s for dashboard with ID %d. Error: %s", fetching, id, err)
	}
}

// ConstructUrlForApiOrTest checks if a function is to be used for testing or not and constructs the url based on that
func ConstructUrlForApiOrTest(urlPath, testUrl string, runTest bool) string {
	url := ""
	if runTest == false {
		url = urlPath
	} else if runTest == true {
		url = testUrl
	}
	return url
}

// RecognizeEnvironmentVariableForClientContext checks if the environment variable is set to use the Firestore emulator
func RecognizeEnvironmentVariableForClientContext(client *firestore.Client, ctx context.Context) (*firestore.Client, context.Context) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "8081" {
		client = firestoreEmulator.GetEmulatorClient()
		ctx = firestoreEmulator.GetEmulatorContext()
	} else {
		client = database.GetFirestoreClient()
		ctx = database.GetFirestoreContext()
	}
	return client, ctx
}
