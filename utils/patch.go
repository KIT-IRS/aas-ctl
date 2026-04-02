package utils

import (
	"net/http"
	"strings"
)

// PatchRequest function sends a PATCH request to the specified endpoint with the provided data.
// Patch requests are used to apply partial modifications to a resource.
// In this case, the ".../$value" can be used to update the value of a submodel element just as property, without the need to completely edit the whole element json.
func PatchRequest(endpoint string, data []byte) ([]byte, error) {
	return HttpRequest(http.MethodPatch, ensureValueSuffix(endpoint), data)
}

func ensureValueSuffix(endpoint string) string {
	if !strings.HasSuffix(endpoint, "/$value") {
		return endpoint + "/$value"
	}
	return endpoint
}
