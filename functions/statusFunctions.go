package functions

import (
	"fmt"
	"net/http"
	"time"
)

// CheckEndpointStatus checks and returns the status of an endpoint.
// If the endpoint does not respond within 10 seconds it is timed out
func CheckEndpointStatus(url string) int {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Get(url)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return http.StatusServiceUnavailable
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()
	return response.StatusCode
}
