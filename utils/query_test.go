package utils

import (
	"testing"

	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
	"github.com/stretchr/testify/assert"
)

func TestFindSubmodel(t *testing.T) {
	t.Skip("TODO: Integration test")
}

func TestFindSubmodelElement(t *testing.T) {
	tests := []struct {
		submodel  aastypes.ISubmodel
		elementId string
		element   aastypes.ISubmodelElement
		wantErr   bool
		errType   error
	}{
		{submodel, "URIOfTheProduct", propety, false, nil},
		{submodel, "ManufacturerName", mlp, false, nil},
		{submodel, "ContactInformation", smec, false, nil},
		{submodel, "", nil, true, &IdShortNotFoundError{}},
		{submodel, "RoleOfContactPerson", nil, true, &IdShortNotFoundError{}}, //element of SMEC
	}
	for _, tt := range tests {
		want := tt.element
		got, err := FindSubmodelElement(tt.submodel, tt.elementId)
		if tt.wantErr {
			assert.Error(t, err)
			assert.ErrorAs(t, err, &tt.errType)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, want, got)
	}
}

func TestGetSubmodelElement(t *testing.T) {
	tests := []struct {
		submodel   aastypes.ISubmodel
		elementIdx int
		element    aastypes.ISubmodelElement
		wantErr    bool
		errType    error
	}{
		{submodel, 0, propety, false, nil},
		{submodel, 4, mlp, false, nil},
		{submodel, 7, smec, false, nil},
		{submodel, -1, nil, true, &IdShortNotFoundError{}},
		{submodel, 99, nil, true, &IdShortNotFoundError{}},
	}
	for _, tt := range tests {
		want := tt.element
		got, err := GetSubmodelElement(tt.submodel, tt.elementIdx)
		if tt.wantErr {
			assert.Error(t, err)
			assert.ErrorAs(t, err, &tt.errType)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, want, got)
	}
}

func TestFindCollectionElement(t *testing.T) {
	tests := []struct {
		collection aastypes.ISubmodelElementCollection
		elementId  string
		element    aastypes.ISubmodelElement
		wantErr    bool
		errType    error
	}{
		{smec, "RoleOfContactPerson", smec_property, false, nil},
		{smec, "NationalCode", smec_mlp, false, nil},
		{smec, "Phone", smec_smec, false, nil},
		{smec, "", nil, true, &IdShortNotFoundError{}},
		{smec, "anything", nil, true, &IdShortNotFoundError{}},
	}
	for _, tt := range tests {
		want := tt.element
		got, err := findCollectionElement(tt.collection, tt.elementId)
		if tt.wantErr {
			assert.Error(t, err)
			assert.ErrorAs(t, err, &tt.errType)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, want, got)
	}
}
