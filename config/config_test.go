package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfig = Config{
	ActiveProfile: &localhostProfile,
	Profiles:      []Profile{exampleProfile, localhostProfile},
}

func TestGetConfigPath(t *testing.T) {
	_, err := getConfigPath()
	assert.NoError(t, err) // path can not be tested, just ensure there is no error
}

func TestLoadConfigFrom(t *testing.T) {
	want := &testConfig
	got, err := loadConfigFrom("./config.json")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestAddProfile(t *testing.T) {
	tests := []struct {
		config             *Config
		profile            *Profile
		newProfileIsActive bool
		wantErr            bool
	}{
		{&Config{}, &localhostProfile, true, false},
		{&testConfig, &exampleProfile, false, true},
		{&testConfig, &Profile{Name: "new"}, false, false},
	}
	for _, tt := range tests {
		cfg := tt.config
		err := cfg.AddProfile(tt.profile)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Contains(t, cfg.Profiles, *tt.profile)
		if tt.newProfileIsActive {
			assert.Equal(t, tt.profile, cfg.ActiveProfile)
		}
	}
}
