package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

const (
	envName         = "CHARMIL_CONFIG_PATH_RHOAS"
	defaultFileName = "rhoas_config.json"

	// TestPath can be used for testing purposes
	TestPath = "mock_location.json"
)

// CfgHandler defines the fields required to manage config.
type CfgHandler struct {
	// Pointer to an instance of the host CLI config struct
	Cfg *Config

	// Path of the local config file
	FilePath string

	// Extension of the local config file
	fileExt string
}

// NewHandler links the specified arguments to a
// new instance of config handler and returns a pointer to it.
func NewHandler(cfg *Config) (*CfgHandler, error) {

	path, err := location()
	if err != nil {
		return nil, err
	}

	h := &CfgHandler{
		FilePath: path,
		Cfg:      cfg,
		fileExt:  filepath.Ext(path),
	}

	return h, nil
}

// Load reads config values from the local config file
// (using the file path linked to the handler) and stores
// them into the linked instance of host CLI config struct.
func (h *CfgHandler) Load() error {

	// Reads the local config file
	buf, err := readFile(h.FilePath)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	// Stores values (read from file) to the host config struct instance
	err = Unmarshal(buf, h.Cfg, h.fileExt)
	if err != nil {
		return err
	}

	return nil
}

// Save writes config values from the linked instance
// of host CLI config struct to the local config file
// (using the file path linked to the handler).
func (h *CfgHandler) Save() error {
	if h.FilePath == TestPath {
		return nil
	}

	// Stores the host CLI config as a byte array
	buf, err := Marshal(h.Cfg, h.fileExt)
	if err != nil {
		return err
	}

	// Writes the current host CLI config to the local config file
	err = writeFile(h.FilePath, buf)
	if err != nil {
		return err
	}

	return nil
}

// Marshal converts the passed object into byte data, based on the specified file format
func Marshal(in interface{}, fileExt string) ([]byte, error) {
	var marshalFunc func(in interface{}) ([]byte, error)

	switch fileExt {
	case ".yaml", ".yml":
		marshalFunc = yaml.Marshal

	case ".toml":
		marshalFunc = toml.Marshal

	case ".json":
		buf, err := json.MarshalIndent(in, "", "  ")
		if err != nil {
			return nil, err
		}
		return buf, nil

	default:
		return nil, fmt.Errorf("Unsupported file extension \"%v\"", fileExt)
	}

	buf, err := marshalFunc(in)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Unmarshal converts the passed byte data into a struct
func Unmarshal(in []byte, out interface{}, fileExt string) error {
	var unmarshalFunc func(in []byte, out interface{}) (err error)

	switch fileExt {
	case ".yaml", ".yml":
		unmarshalFunc = yaml.Unmarshal
	case ".json":
		unmarshalFunc = json.Unmarshal
	case ".toml":
		unmarshalFunc = toml.Unmarshal
	default:
		return fmt.Errorf("Unsupported file extension \"%v\"", fileExt)
	}

	err := unmarshalFunc(in, out)
	if err != nil {
		return err
	}

	return nil
}

// readFile reads the file specified by filePath and returns its contents.
func readFile(filePath string) ([]byte, error) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// writeFile writes data to the file specified by filePath.
func writeFile(filePath string, data []byte) error {
	err := ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

// location gets the path to the config file
func location() (string, error) {
	var path string

	if envCfgPath := os.Getenv(envName); envCfgPath != "" {
		path = envCfgPath
	} else {
		defaultDirPath, err := defaultDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(defaultDirPath, defaultFileName)
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		e := os.MkdirAll(filepath.Dir(path), 0700)
		if e != nil {
			return "", e
		}
	}

	return path, nil
}

// defaultDir returns the default parent directory of the config file
func defaultDir() (string, error) {
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userCfgDir, "charmil"), nil
}
