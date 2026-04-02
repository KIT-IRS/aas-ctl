package utils

import (
	"fmt"
	"strconv"

	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
)

// RequireSingleArg function checks if exactly one argument is provided.
// If not, an error is returned.
func RequireSingleArg(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, got %v", len(args))
	}
	return nil
}

// RequireMinArgs function checks if there are at least n args provided.
// If not, an error is returned.
func RequireMinArgs(args []string, n int) error {
	if n < 1 {
		return fmt.Errorf("n must be >= 1, got %d", n)
	}
	if len(args) < n {
		return fmt.Errorf("expected at least %v argument(s), got %v", n, len(args))
	}
	return nil
}

// isInt function determines wehter a string can be interpreted as integer by strconv.Atoi()
// Maybe move to anoter file?
func isInt(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}

// Determines if the provided Shell contains a submodel with the given IDShort
func smIdShortInShell(shell aastypes.IAssetAdministrationShell, smIdShort string) (bool, error) {
	for _, ref := range shell.Submodels() {
		sm, err := getSubmodelFromReference(ref)
		if err != nil {
			return false, err
		}
		if sm.IDShort() != nil && *sm.IDShort() == smIdShort {
			return true, nil
		}
	}
	return false, nil
}

// elementIdInSubmodel determines if a Submodel has a SubmodelElement with a given IDShort
func elementIdInSubmodel(sm aastypes.ISubmodel, elementId string) bool {
	for _, element := range sm.SubmodelElements() {
		if element.IDShort() != nil && *element.IDShort() == elementId {
			return true
		}
	}
	return false
}

// elementIdxInSubmodel determines if a Submodel has a SubmodelElement with a given index.
// This is required since IDShort is an optional field.
func elementIdxInSubmodel(sm aastypes.ISubmodel, elementIdx int) bool {
	return 0 <= elementIdx && elementIdx < len(sm.SubmodelElements())
}

// elementValueEquals if the value of a SubmodelElement equals a given value.
// This method is required, since not all SubmodelElements have a value attribute.
// Currently this function is only implemented for MultiLanguageProperty and Property.
func elementValueEquals(element aastypes.ISubmodelElement, value string) bool {
	switch element := element.(type) {
	case aastypes.IMultiLanguageProperty:
		return mlpValueEquals(element, value)
	case aastypes.IProperty:
		return *element.Value() == value
	default:
		return false
	}
}

// mlpValueEquals function checks if any of the LangString in the MultiLanguageProperty equals the given value.
func mlpValueEquals(mlp aastypes.IMultiLanguageProperty, value string) bool {
	for _, ls := range mlp.Value() {
		if ls.Text() == value {
			return true
		}
	}
	return false
}

// elementIdxInCollection determines if a SubmodelElementCollection has a SubmodelElement at a given index.
func elementIdxInCollection(smc aastypes.ISubmodelElementCollection, idx int) bool {
	return 0 <= idx && idx < len(smc.Value())
}
