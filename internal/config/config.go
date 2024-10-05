package config

import (
	"errors"
	"fmt"
	"os"
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

// Parse parses the config in the target directory.
func Parse(target string) (Config, error) {
	for _, ext := range configExtensions {
		path := buildPath(target, ext)
		if !fileExists(path) {
			continue
		}
		var config Config
		err := parseYAMLFile(path, &config)
		if err != nil {
			return config, err
		}
		return config, nil
	}
	return Config{}, ErrNoConfig
}

func buildPath(target, ext string) string {
	return filepath.Join(target, fmt.Sprintf("%s.%s", configName, ext))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func parseYAMLFile(path string, v any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(v); err != nil {
		return err
	}
	return nil
}
