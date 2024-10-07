package config

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const configName = "mochi"

var configExtensions = [2]string{"yaml", "yml"}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ErrNoConfig is the error returned when no config is found in the target directory.
var ErrNoConfig = errors.New("no config found in target")

// Config represents a config.
type Config struct {
	Decks []Deck `yaml:"decks" validate:"required,dive"` // sorted by longest Path (more specific first)
}

// Deck represents a sync config.
type Deck struct {
	Path string  `yaml:"path" validate:"required"`
	Name *string `yaml:"name" validate:"gt=0"`
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

	var config Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}

	if err := validate.Struct(&config); err != nil {
		return nil, err
	}

	config = cleanConfig(config)
	return &config, nil
}

func cleanConfig(config Config) Config {
	for i, deck := range config.Decks {
		path := filepath.Clean(filepath.Join("/", deck.Path))
		config.Decks[i].Path = path
	}

	slices.SortFunc(config.Decks, func(a, b Deck) int {
		return len(b.Path) - len(a.Path)
	})

	return config
}

// GetDeck returns the deck config that matches the path.
func (c *Config) GetDeck(path string) (Deck, bool) {
	for _, deck := range c.Decks {
		if deck.Path == path {
			return deck, true
		}
	}
	return Deck{}, false
}
