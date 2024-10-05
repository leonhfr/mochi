package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/leonhfr/mochi/internal/file"
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
		if !file.Exists(path) {
			continue
		}
		var config Config
		err := file.ParseYAML(path, &config)
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
