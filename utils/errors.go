/*
Package utils contains some functions used by other packages.
Contains the implementation of the REST-API requests.
*/

package utils

import (
	"fmt"
)

// HTTPError represents a custom error for HTTP responses.
type HTTPError struct {
	StatusCode int    // HTTP status code
	StatusText string // Status text (e.g., "404 Not Found")
	URL        string // URL that caused the error
}

// Error method makes HTTPError implement the error interface.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error: %d %s (URL: %s)", e.StatusCode, e.StatusText, e.URL)
}

// IdNotFoundError is raised, if the requested (global) id can not be found in the repository.
type IdNotFoundError struct {
	Id string
}

// Error method makes IdNotFoundError implement the error interface.
func (e *IdNotFoundError) Error() string {
	return fmt.Sprintf("Unable to find element with Id %v", e.Id)
}

// IdShortNotFoundError is raised, if the requested idShort is not found in the given namespace.
// May be in the global namespace as well as an AAS or a Submodel namespace.
type IdShortNotFoundError struct {
	IdShort string
}

// Error method makes IdShortNotFoundError implement the error interface.
func (e *IdShortNotFoundError) Error() string {
	return fmt.Sprintf("Unable to find element with IdShort %v", e.IdShort)
}

// IdentifiableNotFoundError is raised if neither an Id nor an IdShort matching the identifer has ben found.
type IdentifiableNotFoundError struct {
	Identifier string
}

// Error method makes IdentifiableNotFoundError implement the error interface.
func (e *IdentifiableNotFoundError) Error() string {
	return fmt.Sprintf("Unable to find identifiable with Id or IdShort %v", e.Identifier)
}

// SmIdShortNotInShellError gets raised, if a AAS does not contain a Submodel with the requested idShort.
type SmIdShortNotInShellError struct {
	shellIdShort string
	smIdShort    string
}

// Error method makes SmIdShortNotFound implement the error interface.
func (e *SmIdShortNotInShellError) Error() string {
	return fmt.Sprintf("Unable to find submodel with IdShort %v in shell with IdShort %v", e.smIdShort, e.shellIdShort)
}
