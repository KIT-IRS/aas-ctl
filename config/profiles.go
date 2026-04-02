package config

import (
	"fmt"

	"github.com/fatih/color"
)

type IProfile interface {
	GetName() string
	Discovery() string
	Registry() string
	SmRegistry() string
	Repository() string
	SmRepository() string
	ConceptDescriptions() string
	Print()
	PrintActive()
}

// Profile struct defines an AAS repository.
// Each profile is identified by its Name attribute.
type Profile struct {
	Name  string `mapstructure:"name" structs:"name" env:"name" json:"name"`
	URL   string `mapstructure:"url" structs:"url" env:"url" json:"url"`
	Ports Ports  `mapstructure:"ports" structs:"ports" env:"ports" json:"ports"`
}

// Ports struct is used to specify the ports of the different endpoints for each profile.
type Ports struct {
	Discovery           int `mapstructure:"discovery" structs:"discovery" env:"discovery" json:"discovery"`
	Registry            int `mapstructure:"registry" structs:"registry" env:"registry" json:"registry"`
	SmRegistry          int `mapstructure:"sm-registry" structs:"sm-registry" env:"sm-registry" json:"sm-registry"`
	Repository          int `mapstructure:"repository" structs:"repository" env:"repository" json:"repository"`
	SmRepository        int `mapstructure:"sm-repository" structs:"sm-repository" env:"sm-repository" json:"sm-repository"`
	ConceptDescriptions int `mapstructure:"concept-descriptions" structs:"concept-descriptions" env:"concept-descriptions" json:"concept-descriptions"`
}

// CreateProfileWithName function creates a new Profile struct with a given name
func CreateProfileWithName(name string) *Profile {
	return &Profile{Name: name}
}

// GetName function returns the name of the Profile
func (p *Profile) GetName() string {
	return p.Name
}

// Discovery function returns the AAS Discovery URL of the Profile
func (p *Profile) Discovery() string {
	return fmt.Sprintf("%v:%d/lookup/shells", p.URL, p.Ports.Discovery)
}

// Registry function returns The AAS Registry URL of the Profile
func (p *Profile) Registry() string {
	return fmt.Sprintf("%v:%d/shell-descriptors", p.URL, p.Ports.Registry)
}

// SmRegistry function returns the Submodel Registry URL of the Profile
func (p *Profile) SmRegistry() string {
	return fmt.Sprintf("%v:%d/submodel-descriptors", p.URL, p.Ports.SmRegistry)
}

// Repository function returns the AAS Repository URL of the Profile
func (p *Profile) Repository() string {
	return fmt.Sprintf("%v:%d/shells", p.URL, p.Ports.Repository)
}

// SmRepository function returns the Submodel Repository URL of the Profile
func (p *Profile) SmRepository() string {
	return fmt.Sprintf("%v:%d/submodels", p.URL, p.Ports.SmRepository)
}

// ConceptDescriptions function returns the Concept Description Repository URL of the Profile
func (p *Profile) ConceptDescriptions() string {
	return fmt.Sprintf("%v:%d/concept-descriptions", p.URL, p.Ports.ConceptDescriptions)
}

// Print function prints the Profile to the Stdout
func (p *Profile) Print() {
	fmt.Println(p.format())
}

// PrintActive prints the Profile in green to the Stdout
func (p *Profile) PrintActive() {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green(p.format()))
}

// format function formats the Profile to a printable string
func (p *Profile) format() string {
	return fmt.Sprintf("Name: %v\nURL: %v", p.Name, p.URL)
}
