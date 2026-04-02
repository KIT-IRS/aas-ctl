package utils

import (
	"bytes"
	"io"
	"net/http"
)

// HttpRequest function sends an HTTP request with the specified method to the given endpoint with the provided data.
func HttpRequest(method string, endpoint string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !IsSuccessStatus(resp.StatusCode) {
		err = &HTTPError{
			StatusCode: resp.StatusCode,
			StatusText: http.StatusText(resp.StatusCode),
			URL:        endpoint,
		}
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// IsSuccessStatus function checks if the HTTP status code indicates a successful response (2xx).
func IsSuccessStatus(code int) bool {
	return code >= 200 && code < 300
}
