package utils

import (
	"errors"
	"fmt"

	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
)

// FindSubmodel function searches a Submodel with a given IDShort in a AAS
// if no SM with matching IDShort is found, an IdShortNotFoundError is returned
func FindSubmodel(shell aastypes.IAssetAdministrationShell, smIdShort string) (aastypes.ISubmodel, error) {
	for _, ref := range shell.Submodels() {
		sm, err := getSubmodelFromReference(ref)
		if err != nil {
			return nil, err
		}
		if sm.IDShort() != nil && *sm.IDShort() == smIdShort {
			return sm, nil
		}
	}
	return nil, &IdShortNotFoundError{IdShort: smIdShort}
}

// FindSubmodelElement function searches for a SME with a given IDShort in a SM
// if no SME with a matching IDShort is found an IdShortNotFoundError is returned
func FindSubmodelElement(sm aastypes.ISubmodel, elementId string) (aastypes.ISubmodelElement, error) {
	for _, e := range sm.SubmodelElements() {
		if e.IDShort() != nil && *e.IDShort() == elementId {
			return e, nil
		}
	}
	return nil, &IdShortNotFoundError{IdShort: elementId}
}

// GetSubmodelElemnt function returns the SME at the given index, if it exists
// if the index does not exist, an error is returned
func GetSubmodelElement(sm aastypes.ISubmodel, elementIdx int) (aastypes.ISubmodelElement, error) {
	if len(sm.SubmodelElements()) == 0 {
		return nil, errors.New("submodel does not contain any SubmodelElements")
	}
	if elementIdx < 0 || elementIdx >= len(sm.SubmodelElements()) {
		return nil, fmt.Errorf("index %d is out of bounds (0-%d)", elementIdx, len(sm.SubmodelElements())-1)
	}
	return sm.SubmodelElements()[elementIdx], nil
}

// findCollectionElement function searches for a SME with a given IDShort in a SMC
// if no SME with a matching IDShort is found, an IdShortNotFoundError is returned
func findCollectionElement(smc aastypes.ISubmodelElementCollection, idShort string) (aastypes.ISubmodelElement, error) {
	for _, element := range smc.Value() {
		if element.IDShort() != nil && *element.IDShort() == idShort {
			return element, nil
		}
	}
	return nil, &IdShortNotFoundError{IdShort: idShort}
}

// getCollectionElement function returns the SME at the given index in a SMC, if it exists;
// if the index does not exist, an error is returned.
func getCollectionElement(smc aastypes.ISubmodelElementCollection, idx int) (aastypes.ISubmodelElement, error) {
	if elementIdxInCollection(smc, idx) {
		return smc.Value()[idx], nil
	}
	if len(smc.Value()) == 0 {
		return nil, errors.New("SubmodelElementCollection does not contain any SubmodelElements")
	}
	return nil, fmt.Errorf("index %d is out of bounds (0-%d)", idx, len(smc.Value())-1)
}
