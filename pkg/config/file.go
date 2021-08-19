package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NewFile creates a new config type
func NewFile() IConfig {
	cfg := &File{}

	return cfg
}

// File is a type which describes a config file
type File struct{}

const errorFormat = "%v: %w"

const EnvName = "RHOASCONFIG"

// Load loads the configuration from the configuration file. If the configuration file doesn't exist
// it will return an empty configuration object.
func (c *File) Load() (*Config, error) {
	file, err := c.Location()
	if err != nil {
		return nil, err
	}

	// #nosec G304
	data, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		return &Config{}, nil
	} else if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, "unable to parse config", err)
	}
	return &cfg, nil
}

// Save saves the given configuration to the configuration file.
func (c *File) Save(cfg *Config) error {
	file, err := c.Location()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("%v: %w", "unable to marshal config", err)
	}

	err = ioutil.WriteFile(file, data, 0o600)
	if err != nil {
		return fmt.Errorf(errorFormat, "unable to save config", err)
	}
	return nil
}

// Remove removes the configuration file.
func (c *File) Remove() error {
	file, err := c.Location()
	if err != nil {
		return err
	}
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return nil
	}
	err = os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}

// Location gets the path to the config file
func (c *File) Location() (path string, err error) {
	if rhoasConfig := os.Getenv(EnvName); rhoasConfig != "" {
		path = rhoasConfig
	} else {
		rhoasCfgDir, e := DefaultDir()
		if e != nil {
			return "", e
		}
		path = filepath.Join(rhoasCfgDir, "plugin_config.json")
		if e != nil {
			return "", e
		}
	}

	if _, err = os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		e := os.MkdirAll(filepath.Dir(path), 0700)
		if e != nil {
			return "", e
		}
	}

	return path, nil
}

// Checks if config has custom location
func HasCustomLocation() bool {
	rhoasConfig := os.Getenv(EnvName)
	return rhoasConfig != ""
}

// DefaultDir returns the default parent directory of the config file
func DefaultDir() (string, error) {
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userCfgDir, "rhoas"), nil
}
