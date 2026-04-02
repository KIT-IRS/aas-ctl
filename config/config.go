/*
Package config implements structs and functions necessary for the access and selection of diffent profiles.
*/
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

// Config struct represents the configuration denoted in the cofig file.
type Config struct {
	ActiveProfile *Profile  `mapstructure:"activeProfile" structures:"activeProfile" env:"activeProfile" json:"activeProfile"`
	Profiles      []Profile `mapstructure:"profiles" structs:"profiles" env:"profiles" json:"profiles"`
}

// Path to the configuration file
var ConfigFile = func() string {
	path, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	return path
}()

// getConfigPath function returns the path of the configfile.
func getConfigPath() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(dirname, ".aas", "config.json"), nil
}

// configFileExists function checks if the ConfigFile ~/.aas/config.json exists
func configFileExists() bool {
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// LoadConfig reads the configuration from the file.
// If the config file does not exist, a new (empty) config is created
func LoadConfig() (*Config, error) {
	if !configFileExists() {
		return &Config{}, nil
	}
	return loadConfigFrom(ConfigFile)
}

// loadConfigFrom function loads the configuration from the given source path
func loadConfigFrom(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetActiveProfile function returns a pointer to a Profile struct that contains the active profile information.
func GetActiveProfile() (IProfile, error) {
	configs, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return configs.ActiveProfile, nil
}

// Save writes the configuration to the file.
func (c *Config) Save() error {
	file, err := os.Create(ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var encoder = json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(c)
}

// Adds the profile to the config.
// If the config has no active profile, the new profile is the selected one
func (c *Config) AddProfile(new *Profile) error {
	for _, p := range c.Profiles {
		if p.GetName() == new.GetName() {
			return fmt.Errorf("profile with name %v already in config", new.GetName())
		}
	}
	c.Profiles = append(c.Profiles, *new)
	if c.ActiveProfile == nil {
		c.ActiveProfile = new
	}
	return nil
}
