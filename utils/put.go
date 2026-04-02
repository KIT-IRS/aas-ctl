package utils

import (
	"net/http"
)

// PutRequest function sends a PUT request to the specified endpoint with the provided data.
// Put requests are used to modify a resource entirely.
// In this case, for the modification of a submodel element, the whole updated submodel element JSON is needed.
func PutRequest(endpoint string, data []byte) ([]byte, error) {
	return HttpRequest(http.MethodPut, endpoint, data)
}
