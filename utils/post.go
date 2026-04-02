package utils

import (
	"net/http"
)

// PostRequest function sends a POST request to the specified endpoint with the provided data.
// Post requests are used to create a new resource or submit data to a server.
// In this case, it can be used to create new submodels or submodel elements.
func PostRequest(endpoint string, data []byte) ([]byte, error) {
	return HttpRequest(http.MethodPost, endpoint, data)
}
