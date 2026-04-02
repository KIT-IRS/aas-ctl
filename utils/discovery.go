package utils

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	aasjsonization "github.com/aas-core-works/aas-core3.0-golang/jsonization"
	aasstringification "github.com/aas-core-works/aas-core3.0-golang/stringification"
	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
)

// Function BuildDiscoveryTree builds the discovery tree for the given endpoint and returns it.
// The discovery tree is a slice of strings containing the elements on the next layer.
func BuildDiscoveryTree(endpoint string) ([]string, error) {
	jsonable, err := GetJsonable(endpoint)
	if err != nil {
		return nil, err
	}
	var deserializationError *aasjsonization.DeserializationError
	aas, err := aasjsonization.AssetAdministrationShellFromJsonable(jsonable)
	if err == nil {
		return buildAASDiscoveryTree(aas)
	} else if !errors.As(err, &deserializationError) {
		return nil, err
	}
	sm, err := aasjsonization.SubmodelFromJsonable(jsonable)
	if err == nil {
		return buildSMDiscoveryTree(sm)
	} else if !errors.As(err, &deserializationError) {
		return nil, err
	}
	element, err := aasjsonization.SubmodelElementFromJsonable(jsonable)
	if err == nil {
		return buildSMEDiscoveryTree(element)
	} else if errors.As(err, &deserializationError) {
		return nil, err
	}
	return nil, fmt.Errorf("unable to build discovery tree for endpoint %v", endpoint)
}

// buildAASDiscoveryTree function builds the discovery tree for an AAS endpoint.
// Function is neccessary because the path "*aas-endpoint*/$value" doesn't exist.
func buildAASDiscoveryTree(aas aastypes.IAssetAdministrationShell) ([]string, error) {
	var tree []string
	for _, ref := range aas.Submodels() {
		sm, err := getSubmodelFromReference(ref)
		if err != nil {
			log.Printf("WARNING: Unable to get Submodel from reference %v; %v", ref.Keys()[0].Value(), err)
		}
		tree = append(tree, formatIdentifiable(sm))
	}
	return tree, nil
}

// buildSMDiscoveryTree function builds the discovery tree for a Submodel endpoint.
// Function is neccessary because "*sm-endpoint*/$value" does return an unordered list of SubmodelElements, which also should be accessible via index.
func buildSMDiscoveryTree(sm aastypes.ISubmodel) ([]string, error) {
	var tree []string
	for i, sme := range sm.SubmodelElements() {
		if sme.IDShort() != nil {
			tree = append(tree, *sme.IDShort())
		} else {
			tree = append(tree, fmt.Sprintf("[%d]", i))
		}
	}
	return tree, nil
}

// buildSMEDiscoveryTree function builds the discovery tree for a SubmodelElement endpoint.
func buildSMEDiscoveryTree(sme aastypes.ISubmodelElement) ([]string, error) {
	switch sme := sme.(type) {
	case aastypes.IMultiLanguageProperty:
		return []string{formatMultiLanguageProperty(sme)}, nil
	case aastypes.IProperty:
		return []string{formatProperty(sme)}, nil
	case aastypes.IRange:
		return []string{formatRange(sme)}, nil
	case aastypes.ISubmodelElementList:
		return buildSMEListDiscoveryTree(sme)
	case aastypes.ISubmodelElementCollection: // Order is relevant, SEC must be after SEL
		return buildSMECollectionDiscoveryTree(sme)
	default:
		return nil, errors.New("unable to build discovery tree")
	}
}

// buildSMEListDiscoveryTree function builds the discovery tree for a SubmodelElementList endpoint.
func buildSMEListDiscoveryTree(sel aastypes.ISubmodelElementList) ([]string, error) {
	var tree []string
	for i, sme := range sel.Value() {
		if sme.IDShort() != nil {
			tree = append(tree, *sme.IDShort())
		} else {
			tree = append(tree, fmt.Sprintf("[%d]", i))
		}
	}
	return tree, nil
}

// buildSMECollectionDiscoveryTree function builds the discovery tree for a SubmodelElementCollection endpoint.
func buildSMECollectionDiscoveryTree(sec aastypes.ISubmodelElementCollection) ([]string, error) {
	var tree []string
	for i, sme := range sec.Value() {
		if sme.IDShort() != nil {
			tree = append(tree, *sme.IDShort())
		} else {
			tree = append(tree, fmt.Sprintf("[%d]", i))
		}
	}
	return tree, nil
}

// ResolveDiscovery function resolves the given args to a specific url endpoint.
// Possible argument combinations, here ID can be ID or IDShort:
// [ShellID] => [Submodel, Submodel, ...]
// [ShellID, SubmodelID] => [SubmodelElement, SubmodelElement, ...]
// [SubmdoelID] => [SubmodelElement, SubmodelElement, ...]
// [ShellID, SubmodelID, SubmodelElementID] => SubmodelElementValue
// [ShellID, SubmodelID, SubmodelElementIdx] => SubmodelElementValue
// [SubmodelID, SubmodelElementID] => SubmodelElementValue
// [SubmodelID, SubmodelElementIdx] => SubmodelElementValue
// [ShellID, SubmodelID, SubmodelElementCollectionID, SubmodelElementID] => SubmodelElementValue
// [ShellID, SubmodelID, SubmodelElementCollectionID, SubmodelElementIdx] => SubmodelElementValue
// list is not complete!
// Elements stored in a list can be accessed via Idx and IDShort (if the element has a IDShort).
// This results in the following requirements for args:
// args[0] must be an identifier for an identifiable (AAS or SM)
// args[1] may be an identifier for an identifiable or a referable (SM or SME)
// args[>1] may be an identifier for a referable or an index
func ResolveDiscovery(args []string) (string, error) {
	identifiable, err := GetIdentifiable(args[0])
	if err != nil {
		return "", err
	}
	// identifiable is either AAS or SM
	switch identfiable := identifiable.(type) {
	case aastypes.IAssetAdministrationShell:
		return resolveAASDiscovery(identfiable, args[1:])
	case aastypes.ISubmodel:
		return resolveSMDiscovery(identfiable, args[1:])
	}
	return "", errors.New("unknown error occured during discovery")
}

// resolveAASDiscovery function resolves the discovery path for the given AAS.
// The endpoint URL is returned.
func resolveAASDiscovery(aas aastypes.IAssetAdministrationShell, args []string) (string, error) {
	if len(args) == 0 {
		return getEndpoint(aas)
	}
	sm, err := GetShellSubmodel(aas.ID(), args[0])
	if err == nil {
		return resolveSMDiscovery(sm, args[1:])
	}
	var smIdShortNotInShell *SmIdShortNotInShellError
	if errors.As(err, &smIdShortNotInShell) && isInt(args[0]) {
		idx, _ := strconv.Atoi(args[0]) // err is nil because it passed isInt function
		sm, err = getSubmodelFromReference(aas.Submodels()[idx])
		if err != nil {
			return "", err
		}
		return resolveSMDiscovery(sm, args[1:])
	}
	return "", err
}

// resolveSMDiscovery function resolves the discovery oath for the given SM.
// The endpoint URL is returned.
func resolveSMDiscovery(sm aastypes.ISubmodel, args []string) (string, error) {
	smEndpoint, err := getEndpoint(sm)
	if err != nil {
		return "", err
	}
	if len(args) == 0 {
		return smEndpoint, nil
	}
	var sme aastypes.ISubmodelElement
	if isInt(args[0]) {
		idx, _ := strconv.Atoi(args[0]) // err is nil
		sme, err = GetSubmodelElement(sm, idx)
		if err != nil {
			return "", err
		}
	} else {
		sme, err = FindSubmodelElement(sm, args[0])
		if err != nil {
			return "", err
		}
	}
	elementEndpoint, err := resolveSMElementDiscovery(sme, args[1:])
	if err != nil {
		return "", err
	}
	return url.JoinPath(smEndpoint, "submodel-elements", elementEndpoint)
}

// resolveSMElementDiscovery resolves the discovery path for the given SME.
// The endpoint URL is returned.
func resolveSMElementDiscovery(sme aastypes.ISubmodelElement, args []string) (string, error) {
	if len(args) == 0 {
		return *sme.IDShort(), nil
	}
	switch sme := sme.(type) {
	case aastypes.ISubmodelElementCollection:
		var collectionElement aastypes.ISubmodelElement
		var err error
		if isInt(args[0]) {
			idx, _ := strconv.Atoi(args[0])
			collectionElement, err = getCollectionElement(sme, idx)
			if err != nil {
				return "", err
			}
		} else {
			collectionElement, err = findCollectionElement(sme, args[0])
			if err != nil {
				return "", err
			}
		}
		collectionEndpoint, err := resolveSMElementDiscovery(collectionElement, args[1:])
		if err != nil {
			return "", err
		}
		return *sme.IDShort() + "." + collectionEndpoint, nil
	}
	modelType, _ := aasstringification.ModelTypeToString((sme.ModelType()))
	return "", fmt.Errorf("discover path exceeded, %v element has no further sub-elements", modelType)
}
