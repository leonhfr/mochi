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

const (
	configName       = "mochi"
	defaultRateLimit = 50
	defaultRootName  = "Root Deck"
)

var configExtensions = [2]string{"yaml", "yml"}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ErrNoConfig is the error returned when no config is found in the target directory.
var ErrNoConfig = errors.New("no config found in target")

// Config represents a config.
type Config struct {
	RateLimit int    `yaml:"rateLimit"` // requests per second
	RootName  string `yaml:"rootName"`
	SkipRoot  bool   `yaml:"skipRoot"`
	Decks     []Deck `yaml:"decks" validate:"required,dive"` // sorted by longest Path (more specific first)
}

// Deck represents a sync config.
type Deck struct {
	Path   string `yaml:"path" validate:"required"`
	Name   string `yaml:"name"`
	Parser string `yaml:"parser" validate:"parser"`
}

// Reader represents the interface to interact with a config file.
type Reader interface {
	Exists(string) bool
	Read(string) (io.ReadCloser, error)
}

// Parse parses the config in the target directory.
func Parse(reader Reader, target string, parsers []string) (*Config, error) {
	if err := validate.RegisterValidation("parser", getValidatorFunc(parsers)); err != nil {
		return nil, err
	}

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

	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if err := validate.Struct(&config); err != nil {
		return nil, err
	}

	config = cleanConfig(config)
	return &config, nil
}

func cleanConfig(config Config) Config {
	if config.RateLimit <= 0 {
		config.RateLimit = defaultRateLimit
	}

	if config.RootName == "" {
		config.RootName = defaultRootName
	}

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
	if path == "/" && c.SkipRoot {
		return Deck{}, false
	} else if path == "/" {
		return Deck{Path: "/", Name: c.RootName}, true
	}

	for _, deck := range c.Decks {
		if deck.Path == path {
			return deck, true
		}
	}
	return Deck{}, false
}

func getValidatorFunc(parsers []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return fl.Field().IsZero() || slices.Contains(parsers, fl.Field().String())
	}
}
