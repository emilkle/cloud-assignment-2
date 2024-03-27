package functions

import (
	"fmt"
	"io"
	"net/http"
)

func CheckEndpointStatus(url string) int {
	statusResponse, err := http.Get(url)
	if err != nil {
		return http.StatusServiceUnavailable
	}
	defer func(Body io.ReadCloser) {
		if err != nil {
			err := Body.Close()
			if err != nil {
				fmt.Printf("failed to close response body from endpoint: %s, during status check. %v", url, err)
			}
		}
	}(statusResponse.Body)
	return statusResponse.StatusCode
}
