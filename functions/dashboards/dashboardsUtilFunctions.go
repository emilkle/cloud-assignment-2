package dashboards

import (
	"io"
	"log"
	"net/http"
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
