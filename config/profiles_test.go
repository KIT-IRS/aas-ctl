package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var localhostProfile = Profile{
	Name:  "localhost",
	URL:   "http://localhost",
	Ports: Ports{Discovery: 8084, Registry: 8082, SmRegistry: 8083, Repository: 8081, SmRepository: 8081, ConceptDescriptions: 8081},
}

var exampleProfile = Profile{
	Name:  "example",
	URL:   "http://example.com/example",
	Ports: Ports{Discovery: 8084, Registry: 8082, SmRegistry: 8083, Repository: 8081, SmRepository: 8081, ConceptDescriptions: 8081},
}

func TestCreateProfileWithName(t *testing.T) {
	want := &Profile{Name: "test"}
	got := CreateProfileWithName("test")
	assert.Equal(t, want, got)
}

func TestGetName(t *testing.T) {
	tests := []struct {
		name    string
		profile *Profile
	}{
		{"localhost", &localhostProfile},
		{"example", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.name
		got := tt.profile.GetName()
		assert.Equal(t, want, got)
	}
}

func TestDiscovery(t *testing.T) {
	tests := []struct {
		discovery string
		profile   *Profile
	}{
		{"http://localhost:8084/lookup/shells", &localhostProfile},
		{"http://example.com/example:8084/lookup/shells", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.discovery
		got := tt.profile.Discovery()
		assert.Equal(t, want, got)
	}
}

func TestRegistry(t *testing.T) {
	tests := []struct {
		discovery string
		profile   *Profile
	}{
		{"http://localhost:8082/shell-descriptors", &localhostProfile},
		{"http://example.com/example:8082/shell-descriptors", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.discovery
		got := tt.profile.Registry()
		assert.Equal(t, want, got)
	}
}

func TestSmRegistry(t *testing.T) {
	tests := []struct {
		discovery string
		profile   *Profile
	}{
		{"http://localhost:8083/submodel-descriptors", &localhostProfile},
		{"http://example.com/example:8083/submodel-descriptors", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.discovery
		got := tt.profile.SmRegistry()
		assert.Equal(t, want, got)
	}
}

func TestRepository(t *testing.T) {
	tests := []struct {
		discovery string
		profile   *Profile
	}{
		{"http://localhost:8081/shells", &localhostProfile},
		{"http://example.com/example:8081/shells", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.discovery
		got := tt.profile.Repository()
		assert.Equal(t, want, got)
	}
}

func TestSmRepository(t *testing.T) {
	tests := []struct {
		discovery string
		profile   *Profile
	}{
		{"http://localhost:8081/submodels", &localhostProfile},
		{"http://example.com/example:8081/submodels", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.discovery
		got := tt.profile.SmRepository()
		assert.Equal(t, want, got)
	}
}

func TestConceptDescriptions(t *testing.T) {
	tests := []struct {
		discovery string
		profile   *Profile
	}{
		{"http://localhost:8081/concept-descriptions", &localhostProfile},
		{"http://example.com/example:8081/concept-descriptions", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.discovery
		got := tt.profile.ConceptDescriptions()
		assert.Equal(t, want, got)
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		formated string
		profile  *Profile
	}{
		{"Name: localhost\nURL: http://localhost", &localhostProfile},
		{"Name: example\nURL: http://example.com/example", &exampleProfile},
	}
	for _, tt := range tests {
		want := tt.formated
		got := tt.profile.format()
		assert.Equal(t, want, got)
	}
}
