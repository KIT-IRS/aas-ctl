package utils

import (
	"aas-ctl/config"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	aasjsonization "github.com/aas-core-works/aas-core3.0-golang/jsonization"
	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
)

// GetRequest function executes a http get request on the given endpoint and checks if the request was successful.
// First the get request gets executed, then it gets checked if the request was successful.
func GetRequest(endpoint string) ([]byte, error) {
	return HttpRequest(http.MethodGet, endpoint, []byte{})
}

// GetAny function calls getRequest and then unmarshals the response as any
func GetAny(endpoint string) (any, error) {
	response, err := GetRequest(endpoint)
	if err != nil {
		return nil, err
	}
	var raw any
	err = json.Unmarshal(response, &raw)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// GetJsonable function calls getRequest and then unmarshals the response as a jsonable (map[string]any)
func GetJsonable(endpoint string) (map[string]any, error) {
	response, err := GetRequest(endpoint)
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	err = json.Unmarshal(response, &raw)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// GetIdentifiable function returns the Identifiable matching the identifier.
// The identifiable may be a AAS as well as a SM.
// The first matching Identifiable is returned.
// The AAS repository is looked first.
func GetIdentifiable(identifier string) (aastypes.IIdentifiable, error) {
	identifiable, err := GetShell(identifier)
	if err == nil {
		return identifiable, nil
	}
	var identifibleNotFound *IdentifiableNotFoundError
	if !errors.As(err, &identifibleNotFound) {
		return nil, err
	}
	return GetSubmodel(identifier)
}

// GetAllShells function returns a slice of all AAS in the repository.
// Specifies the endpoint based on the current active profile and then calls GetReuqest.
func GetAllShells() ([]aastypes.IAssetAdministrationShell, error) {
	profile, err := config.GetActiveProfile()
	if err != nil {
		return nil, err
	}
	raw, err := GetJsonable(profile.Repository())
	if err != nil {
		return nil, err
	}
	rawShells := raw["result"].([]interface{})
	var shells []aastypes.IAssetAdministrationShell
	for i, r := range rawShells {
		shell, err := aasjsonization.AssetAdministrationShellFromJsonable(r)
		if err != nil {
			log.Printf("WARNING: Unable to load the Shell (index %d) from json; %v", i, err)
		}
		shells = append(shells, shell)
	}
	return shells, nil
}

// GetShell function takes a string identifier and tries to find a matching AAS in the repository.
// The identifier may be an id as well as an idShort.
// The function first tries to get a shell with the identifier as id, if this fails due to an IdNotFoundError, it tries to find a shell with the identifier as idShort.
// If successfull, a pointer to the corresponding AssetAdministrationShell struct is returned.
func GetShell(identifier string) (aastypes.IAssetAdministrationShell, error) {
	shell, err := GetShellById(identifier) // Try find shell by Id
	if err == nil {
		return shell, nil
	}
	var idNotFoundError *IdNotFoundError
	if !errors.As(err, &idNotFoundError) { // Error is not IdNotFoundError
		return nil, err
	}
	shell, err = GetShellByIdShort(identifier) // Try find shell by IdShort
	if err == nil {
		return shell, nil
	}
	var idShortNotFoundError *IdShortNotFoundError
	if !errors.As(err, &idShortNotFoundError) { // Error is not IdShortNotFoundError
		return nil, err
	}
	return nil, &IdentifiableNotFoundError{Identifier: identifier}
}

// GetShellById function takes a string identifier as global id and tries to get the corresponding shell in the active repository.
// If there is no shell with the given id it raises an IdNotFoundError.
// If successful a pointer to the corresponding AssedAdministrationShell struct is returned.
func GetShellById(id string) (aastypes.IAssetAdministrationShell, error) {
	id = base64.StdEncoding.EncodeToString([]byte(id))
	profile, err := config.GetActiveProfile()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%v/%v", profile.Repository(), id)
	raw, err := GetJsonable(endpoint)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == 404 {
			return nil, &IdNotFoundError{}
		}
		return nil, err
	}
	shell, err := aasjsonization.AssetAdministrationShellFromJsonable(raw)
	if err != nil {
		return nil, err
	}
	return shell, nil
}

// GetShellByIdShort function takes a string identifier as idShort and tries to find a corresponding shell.
// Therefore it queries all shells from the active repository and interates over the slice to check if there is a shell with matching idShort.
// If successful, a pointer to a corresponding AssetAdministrationShell struct is returned.
// If there are mutliple shells with the same idShort, the first matching shell is returned.
// If there is no matching shell, an IdShortNotFoundError is raised.
func GetShellByIdShort(idShort string) (aastypes.IAssetAdministrationShell, error) {
	shells, err := GetAllShells()
	if err != nil {
		return nil, err
	}
	for _, shell := range shells {
		if *shell.IDShort() == idShort {
			return shell, nil
		}
	}
	return nil, &IdShortNotFoundError{IdShort: idShort}
}

// GetSubmodel function takes a string identifier and tries to find a matching Submodel in the repository.
// The identifier may be an id as well as an idShort.
// The function first tries to get a Submodel with the identifier as id, if this fails due to an IdNotFoundError, it tries to find a shell with the identifier as idShort.
// If successfull, a pointer to the corresponding Submodel struct is returned.
func GetSubmodel(identifier string) (aastypes.ISubmodel, error) {
	sm, err := GetSubmodelById(identifier)
	if err == nil {
		return sm, nil
	}
	var idNotFoundError *IdNotFoundError
	if !errors.As(err, &idNotFoundError) {
		return nil, err
	}
	sm, err = GetSubmodelByIdShort(identifier)
	if err == nil {
		return sm, nil
	}
	var idShortNotFound *IdShortNotFoundError
	if !errors.As(err, &idShortNotFound) {
		return nil, err
	}
	return nil, &IdentifiableNotFoundError{Identifier: identifier}
}

// GetSubmodelById function takes a string identifier as global id and tries to get the corresponding Submodel in the active repository.
// If there is no Submodel with the given id it raises an IdNotFoundError.
// If successful a pointer to the corresponding Submodel struct is returned.
// The function first creates a RawSubmodel, from which it creates the Submodel afterwards.
func GetSubmodelById(id string) (aastypes.ISubmodel, error) {
	id = base64.StdEncoding.EncodeToString([]byte(id))
	profile, err := config.GetActiveProfile()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%v/%v", profile.SmRepository(), id)
	raw, err := GetJsonable(endpoint)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == 404 {
			return nil, &IdNotFoundError{}
		}
		return nil, err
	}
	submodel, err := aasjsonization.SubmodelFromJsonable(raw)
	if err != nil {
		return nil, err
	}
	return submodel, nil
}

// GetAllSubmodels function returns a slice of all Submodels in the repository.
// Specifies the endpoint based on the current active profile and then calls GetReuqest.
func GetAllSubmodels() ([]aastypes.ISubmodel, error) {
	profile, err := config.GetActiveProfile()
	if err != nil {
		return nil, err
	}
	raw, err := GetJsonable(profile.SmRepository())
	if err != nil {
		return nil, err
	}
	rawSms := raw["result"].([]interface{})
	var submodels []aastypes.ISubmodel
	for i, sm := range rawSms {
		submodel, err := aasjsonization.SubmodelFromJsonable(sm)
		if err != nil {
			log.Printf("WARNING: Unable to load the Submodel (index %d) from json; %v", i, err)
		}
		submodels = append(submodels, submodel)
	}
	return submodels, nil
}

// GetAllSubmodelsOfShell takes an string identifier and tries to find all Submodels of the corresponding shell.
// Tries to get the shell corresponding to the identifier.
// Iterates the submodels attribute of the shell to get the submodel ids.
// Requests the Submodels and add them in a slice.
// Returns the slice.
func GetAllSubmodelsOfShell(identifier string) ([]aastypes.ISubmodel, error) {
	shell, err := GetShell(identifier)
	if err != nil {
		return nil, err
	}
	var submodels []aastypes.ISubmodel
	for _, submodelRef := range shell.Submodels() {
		sm, err := GetSubmodelById(submodelRef.Keys()[0].Value())
		if err != nil {
			return nil, err
		}
		submodels = append(submodels, sm)
	}
	return submodels, nil
}

// GetSubmodelByIdShort function takes a string identifier as idShort and tries to find a corresponding Submodel in the active repository.
// If successful a pointer to a corresponding Submodel struct is returned.
// If there was no Submodel with the given idShort an IdShortNotFoundError is raised.
func GetSubmodelByIdShort(idShort string) (aastypes.ISubmodel, error) {
	submodels, err := GetAllSubmodels()
	if err != nil {
		return nil, err
	}
	for _, submodel := range submodels {
		if *submodel.IDShort() == idShort {
			return submodel, nil
		}
	}
	err = &IdShortNotFoundError{IdShort: idShort}
	return nil, err
}

// GetShellSubmodel function takes two string identifiers, one (id or idShort) for a shell and one (idShort) for a Submodel of the shell.
// Tries to get the shell from the current repository.
// Iterates the submodels of the shell in order to find one with matching idShort.
// If successful a pointer to a corresponding Submodel struct is returned.
// If there was no Submodel with the given idShort, a SmIdShortNotInShellErrors gets raised.
func GetShellSubmodel(shellId string, smIdShort string) (aastypes.ISubmodel, error) {
	shell, err := GetShell(shellId)
	if err != nil {
		return nil, err
	}
	for _, submodelRef := range shell.Submodels() {
		sm, err := GetSubmodelById(submodelRef.Keys()[0].Value())
		if err != nil {
			return nil, err
		}
		if *sm.IDShort() == smIdShort {
			return sm, nil
		}
	}
	err = &SmIdShortNotInShellError{shellIdShort: shellId, smIdShort: smIdShort}
	return nil, err
}

// getSubmodelFromReference takes a struct that implements the IReference interface and checks if it is a reference to a submodel.
// If it is, it calls the GetSubmodelById function with the id provided by the Reference and returns the result.
func getSubmodelFromReference(ref aastypes.IReference) (aastypes.ISubmodel, error) {
	key := ref.Keys()[0]
	if key.Type() != aastypes.KeyTypesSubmodel {
		return nil, errors.New("reference must be of type submodel")
	}
	return GetSubmodelById(key.Value())
}

// getEndpoint function takes a identifiable and creates its endpoint url
func getEndpoint(identifiable aastypes.IIdentifiable) (string, error) {
	id := base64.StdEncoding.EncodeToString([]byte(identifiable.ID()))
	profile, err := config.GetActiveProfile()
	if err != nil {
		return "", err
	}
	switch identifiable.(type) {
	case aastypes.IAssetAdministrationShell:
		return fmt.Sprintf("%v/%v", profile.Repository(), id), nil
	case aastypes.ISubmodel:
		return fmt.Sprintf("%v/%v", profile.SmRepository(), id), nil
	}
	return "", errors.New("something went wrong during the endpoint generation of the identifiable")
}

// getElementEndpoint function takes a Submodel and a SubmodelElement and creates the endpoint url of the SubmodelElement
func getElementEndpoint(sm aastypes.ISubmodel, e aastypes.ISubmodelElement) (string, error) {
	smEndpoint, err := getEndpoint(sm)
	if err != nil {
		return "", err
	}
	if e.IDShort() == nil {
		return "", errors.New("submodel element has no idShort, cannot create endpoint")
	}
	endpoint, err := url.JoinPath(smEndpoint, "submodel-elements", *e.IDShort())
	if err != nil {
		return "", err
	}
	return endpoint, nil
}
