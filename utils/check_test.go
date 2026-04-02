package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/aas-core-works/aas-core3.0-golang/jsonization"
	aastypes "github.com/aas-core-works/aas-core3.0-golang/types"
	"github.com/stretchr/testify/assert"
)

var submodel aastypes.ISubmodel = func() aastypes.ISubmodel {
	raw, _ := os.ReadFile("../test/submodel.json")
	var rawSm map[string]interface{}
	json.Unmarshal(raw, &rawSm)
	sm, _ := jsonization.SubmodelFromJsonable(rawSm)
	return sm
}()

func loadSubmodelElementFromFile(filepath string) aastypes.ISubmodelElement {
	raw, _ := os.ReadFile(filepath)
	var rawElement map[string]interface{}
	json.Unmarshal(raw, &rawElement)
	sme, _ := jsonization.SubmodelElementFromJsonable(rawElement)
	return sme
}

var propety aastypes.ISubmodelElement = loadSubmodelElementFromFile("../test/property.json")

var mlp aastypes.ISubmodelElement = loadSubmodelElementFromFile("../test/mlp.json")

var smec aastypes.ISubmodelElementCollection = loadSubmodelElementFromFile("../test/smec.json").(aastypes.ISubmodelElementCollection)
var smec_property aastypes.ISubmodelElement = loadSubmodelElementFromFile("../test/smec_property.json")
var smec_mlp aastypes.ISubmodelElement = loadSubmodelElementFromFile("../test/smec_mlp.json")
var smec_smec aastypes.ISubmodelElement = loadSubmodelElementFromFile("../test/smec_smec.json")

func TestRequireSingleArg(t *testing.T) {
	tests := []struct {
		args    []string
		wantErr bool
	}{
		{nil, true},
		{[]string{}, true},
		{[]string{"one"}, false},
		{[]string{"one", "two"}, true},
		{[]string{"one", "two", "three"}, true},
	}
	for _, tt := range tests {
		err := RequireSingleArg(tt.args)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestRequireMinArgs(t *testing.T) {
	tests := []struct {
		args    []string
		n       int
		wantErr bool
	}{
		{nil, 1, true},
		{[]string{}, 1, true},
		{[]string{"one"}, 1, false},
		{[]string{"one", "two", "three"}, 1, false},
		{nil, 3, true},
		{[]string{"one", "two"}, 3, true},
		{[]string{"one", "two", "three"}, 3, false},
		{nil, -1, true},
		{[]string{}, -1, true},
		{[]string{"one", "two", "three"}, -1, true},
	}
	for _, tt := range tests {
		err := RequireMinArgs(tt.args, tt.n)
		if tt.wantErr {
			assert.Error(t, err, fmt.Sprintf("RequireMinArgs(%v, %d); expected err, got nil", tt.args, tt.n))
		} else {
			assert.NoError(t, err, fmt.Sprintf("RequireMinArgs(%v, %d); expected nil, got err", tt.args, tt.n))
		}
	}
}

func TestIsInt(t *testing.T) {
	tests := []struct {
		str   string
		isInt bool
	}{
		{"", false},
		{"int", false},
		{"-int", false},
		{"0", true},
		{"-17", true},
		{"012345", true},
		{strconv.Itoa(math.MaxInt), true},
		{strconv.Itoa(math.MinInt), true},
		{"9223372036854775808", false},
		{"-9223372036854775809", false},
	}
	for _, tt := range tests {
		want := tt.isInt
		got := isInt(tt.str)
		assert.Equal(t, want, got)
	}
}

func TestElementIDInSubmodel(t *testing.T) {
	tests := []struct {
		sm              aastypes.ISubmodel
		elementId       string
		containsElement bool
	}{
		{submodel, "URIOfTheProduct", true},
		{submodel, "ContactInformation", true},
		{submodel, "FurtherDetailsOfContract", false},
		{submodel, "something", false},
		{submodel, "", false},
	}
	for _, tt := range tests {
		want := tt.containsElement
		got := elementIdInSubmodel(tt.sm, tt.elementId)
		assert.Equal(t, want, got)
	}
}

func TestElementIdxInSubmodel(t *testing.T) {
	tests := []struct {
		sm              aastypes.ISubmodel
		elementIdx      int
		containsElement bool
	}{
		{submodel, 0, true},
		{submodel, 8, true},
		{submodel, 9, false},
		{submodel, 100, false},
		{submodel, -1, false},
	}
	for _, tt := range tests {
		want := tt.containsElement
		got := elementIdxInSubmodel(tt.sm, tt.elementIdx)
		assert.Equal(t, want, got)
	}
}

func TestElementValueEquals(t *testing.T) {
	tests := []struct {
		element aastypes.ISubmodelElement
		value   string
		equal   bool
	}{
		{propety, "https://www.irs.kit.edu/Composite/GalTWIN/001", true},
		{propety, "", false},
		{mlp, "Karlsruher Institut für Technologie", true},
		{mlp, "de", false},
		{mlp, "", false},
		{smec, "anything", false},
		{smec, "", false},
		{smec, "0173-1#07-AAS931#001", false},
	}
	for _, tt := range tests {
		want := tt.equal
		got := elementValueEquals(tt.element, tt.value)
		assert.Equal(t, want, got)
	}
}

func TestElementIdxInCollection(t *testing.T) {
	tests := []struct {
		smec            aastypes.ISubmodelElementCollection
		elementIdx      int
		containsElement bool
	}{
		{smec, 0, true},
		{smec, 21, true},
		{smec, 22, false},
		{smec, 100, false},
		{smec, -1, false},
	}
	for _, tt := range tests {
		want := tt.containsElement
		got := elementIdxInCollection(tt.smec, tt.elementIdx)
		assert.Equal(t, want, got)
	}
}
