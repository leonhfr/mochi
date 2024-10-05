package config

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configName = "mochi"

var configExtensions = [2]string{"yaml", "yml"}

// ErrNoConfig is the error returned when no config is found in the target directory.
var ErrNoConfig = errors.New("no config found in target")

// Config represents a config.
type Config struct {
	Sync []Sync `yaml:"sync"`
}

// Sync represents a sync config.
type Sync struct {
	Path string `yaml:"path"`
}

// Reader represents the interface to interact with a config file.
type Reader interface {
	Exists(string) bool
	Read(string) (io.ReadCloser, error)
}

// Parse parses the config in the target directory.
func Parse(reader Reader, target string) (*Config, error) {
	for _, ext := range configExtensions {
		path := filepath.Join(target, fmt.Sprintf("%s.%s", configName, ext))
		if !reader.Exists(path) {
			continue
		}

		return parseConfig(reader, path)
	}

	return nil, ErrNoConfig
}

func parseConfig(reader Reader, path string) (*Config, error) {
	r, err := reader.Read(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var config *Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
