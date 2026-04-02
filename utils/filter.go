package utils

import (
	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
)

// IFilter defines the interface a Filter struct must implement.
// Filter structs can be used to filter a list of AssetAdministrationShells for some requested criteria.
type IFilter interface {
	Apply([]aastypes.IAssetAdministrationShell) ([]aastypes.IAssetAdministrationShell, error)
}

// Filter is the base instantiation of the IFilter interface.
// All fields are optional.
// If the filter is applied on a list of AAS, the filter creates a new list with only the AAS that fulfill all criteria.
type Filter struct {
	smIdShort  *string
	elementId  *string
	elementIdx *int
	value      *string
}

// Set the IDShort of a Submodel the filtered AAS must contain
func (f *Filter) SetSmIdShort(idShort string) {
	f.smIdShort = &idShort
}

// Set the IDShort of a SubmodelElement the filtered AAS must contain
// If given in combination with a Submodel IDShort, this Submodel must contain that SubmodelElement
func (f *Filter) SetElementID(id string) {
	f.elementId = &id
}

// Set the index of a SubmodelElement the filtered AAS must contain
// If given in combination with a Submodel IDShort, this Submodel must contain that SubmodelElement
func (f *Filter) SetElementIdx(idx int) {
	f.elementIdx = &idx
}

// Set the value of a SubmodelElement the filtered AAS must contain
// If given in combination wit an index or IDShort of a SubmodelElement, the value of this SubmodelElement must equal the given value.
func (f *Filter) SetValue(value string) {
	f.value = &value
}

// SearchFilterFromFlags function creates a Filter struct from the provided FlagsSearch struct and returns it.
func SearchFilterFromFlags(flags *FlagsSearch) IFilter {
	f := &Filter{}
	if flags.SMID != "" {
		f.SetSmIdShort(flags.SMID)
	}
	if flags.ElementID != "" {
		f.SetElementID(flags.ElementID)
	}
	if flags.ElementIdx != -1 {
		f.SetElementIdx(flags.ElementIdx)
	}
	if flags.Value != "" {
		f.SetValue(flags.Value)
	}
	return f
}

// Apply the filter on a list of AAS.
func (f *Filter) Apply(shells []aastypes.IAssetAdministrationShell) ([]aastypes.IAssetAdministrationShell, error) {
	var err error
	if f.smIdShort != nil {
		shells, err = f.applySm(shells)
		if err != nil {
			return nil, err
		}
	}
	if f.elementId != nil && f.elementIdx == nil {
		shells, err = f.applyElementID(shells)
		if err != nil {
			return nil, err
		}
	}
	if f.elementId == nil && f.elementIdx != nil {
		shells, err = f.applyElementIdx(shells)
		if err != nil {
			return nil, err
		}
	}
	if f.value != nil {
		shells, err = f.applyValue(shells)
		if err != nil {
			return nil, err
		}
	}
	return shells, nil
}

// applySm filters a list of AAS.
// It returns a list of all AAS of the original list, that contain a Submodel with the given IDShort.
func (f *Filter) applySm(shells []aastypes.IAssetAdministrationShell) ([]aastypes.IAssetAdministrationShell, error) {
	var filtered []aastypes.IAssetAdministrationShell
	for _, shell := range shells {
		containsId, err := smIdShortInShell(shell, *f.smIdShort)
		if err != nil {
			continue
		}
		if containsId {
			filtered = append(filtered, shell)
		}
	}
	return filtered, nil
}

// applyElementID filters a list of AAS.
// If only the elementId is given, it returns all AAS that contain a Submodel that has a SubmodelElement with the given IDShort.
// If additionaly a Submodel IDShort is given, it returns only the AAS where the Submodel with the given IDShort contains the required SubmodelElement.
func (f *Filter) applyElementID(shells []aastypes.IAssetAdministrationShell) ([]aastypes.IAssetAdministrationShell, error) {
	var filtered []aastypes.IAssetAdministrationShell
	for _, shell := range shells {
		if f.smIdShort != nil {
			sm, err := FindSubmodel(shell, *f.smIdShort)
			if err != nil {
				continue
			}
			if elementIdInSubmodel(sm, *f.elementId) {
				filtered = append(filtered, shell)
			}
		} else {
			for _, ref := range shell.Submodels() {
				sm, err := getSubmodelFromReference(ref)
				if err != nil {
					continue
				}
				if elementIdInSubmodel(sm, *f.elementId) {
					filtered = append(filtered, shell)
				}
			}
		}
	}
	return filtered, nil
}

// applyElementID filters a list of AAS.
// If only the elementIdx is given, it returns all AAS that contain a Submodel that has a SubmodelElement with the given IDShort.
// If additionaly a Submodel IDShort is given, it returns only the AAS where the Submodel with the given IDShort contains the required SubmodelElement.
func (f *Filter) applyElementIdx(shells []aastypes.IAssetAdministrationShell) ([]aastypes.IAssetAdministrationShell, error) {
	var filtered []aastypes.IAssetAdministrationShell
	for _, shell := range shells {
		if f.smIdShort != nil {
			sm, err := FindSubmodel(shell, *f.smIdShort)
			if err != nil {
				continue
			}
			if elementIdxInSubmodel(sm, *f.elementIdx) {
				filtered = append(filtered, shell)
			}
		} else {
			for _, ref := range shell.Submodels() {
				sm, err := getSubmodelFromReference(ref)
				if err != nil {
					continue
				}
				if elementIdxInSubmodel(sm, *f.elementIdx) {
					filtered = append(filtered, shell)
				}
			}
		}
	}
	return filtered, nil
}

// applyValue filters a list of AAS.
// If only the value is given, it filters all AAS where any Submodel has any SubmodelElement with that value.
// If a Submodel IDShort is provided, it returns all AAS where the requested Submodel contains a SubmodelElement with the given value.
// If a Submodel IDShort and a SubmodelElement IDShort/index is provided, it returns all AAS where this SubmodelElement equals the given value.
func (f *Filter) applyValue(shells []aastypes.IAssetAdministrationShell) ([]aastypes.IAssetAdministrationShell, error) {
	var filtered []aastypes.IAssetAdministrationShell
	for _, shell := range shells {
		if f.smIdShort != nil {
			sm, err := FindSubmodel(shell, *f.smIdShort)
			if err != nil {
				continue
			}
			if f.elementId != nil {
				element, err := FindSubmodelElement(sm, *f.elementId)
				if err != nil {
					continue
				}
				if elementValueEquals(element, *f.value) {
					filtered = append(filtered, shell)
				}
			} else if f.elementIdx != nil {
				element, err := GetSubmodelElement(sm, *f.elementIdx)
				if err != nil {
					continue
				}
				if elementValueEquals(element, *f.value) {
					filtered = append(filtered, shell)
				}
			} else {
				for _, element := range sm.SubmodelElements() {
					if elementValueEquals(element, *f.value) {
						filtered = append(filtered, shell)
						break
					}
				}
			}
		} else {
			for _, ref := range shell.Submodels() {
				sm, err := getSubmodelFromReference(ref)
				if err != nil {
					continue
				}
				if f.elementId != nil {
					element, err := FindSubmodelElement(sm, *f.elementId)
					if err != nil {
						continue
					}
					if elementValueEquals(element, *f.value) {
						filtered = append(filtered, shell)
					}
				} else if f.elementIdx != nil {
					element, err := GetSubmodelElement(sm, *f.elementIdx)
					if err != nil {
						continue
					}
					if elementValueEquals(element, *f.value) {
						filtered = append(filtered, shell)
					}
				} else {
					for _, element := range sm.SubmodelElements() {
						if elementValueEquals(element, *f.value) {
							filtered = append(filtered, shell)
							break
						}
					}
				}
			}
		}
	}
	return filtered, nil
}
