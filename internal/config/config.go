package config

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const (
	configName       = "mochi"
	defaultRateLimit = 50
	defaultRootName  = "Root Deck"
)

var configExtensions = [2]string{"yaml", "yml"}

// ErrNoConfig is the error returned when no config is found in the target directory.
var ErrNoConfig = errors.New("no config found in target")

// Config represents a config.
type Config struct {
	RateLimit  int        `yaml:"rateLimit"` // requests per second
	RootName   string     `yaml:"rootName"`
	SkipRoot   bool       `yaml:"skipRoot"`
	Decks      []Deck     `yaml:"decks" validate:"required,dive"` // sorted by longest Path (more specific first)
	Vocabulary Vocabulary `yaml:"vocabulary"`                     // map[vocabulary name]template id
}

// Deck represents a sync config.
type Deck struct {
	Path   string `yaml:"path" validate:"required"`
	Name   string `yaml:"name"`
	Parser string `yaml:"parser"`
}

// Vocabulary represents the vocabulary templates.
type Vocabulary map[string]string

// Reader represents the interface to read a config file.
type Reader interface {
	Read(string) (io.ReadCloser, error)
}

// Parse parses the config in the target directory.
func Parse(reader Reader, target string, parsers []string) (*Config, error) {
	for _, ext := range configExtensions {
		path := filepath.Join(target, fmt.Sprintf("%s.%s", configName, ext))
		rc, err := reader.Read(path)
		if err == fs.ErrNotExist {
			continue
		} else if err != nil {
			return nil, err
		}
		defer rc.Close()

		return parseConfig(rc, parsers)
	}

	return nil, ErrNoConfig
}

func parseConfig(r io.Reader, parsers []string) (*Config, error) {
	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	validate := validator.New()
	validate.RegisterStructValidation(parsersValidator(parsers), Config{})
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

// Deck returns the deck config that matches the path.
func (c *Config) Deck(path string) (Deck, bool) {
	if path == "/" && c.SkipRoot {
		return Deck{}, false
	} else if path == "/" {
		return Deck{Path: "/", Name: c.RootName}, true
	}

	for _, deck := range c.Decks {
		if isParentDirectory(path, deck.Path) {
			return deck, true
		}
	}
	return Deck{}, false
}

func parsersValidator(parsers []string) validator.StructLevelFunc {
	return func(sl validator.StructLevel) {
		config := sl.Current().Interface().(Config)
		parserNames := parsers
		for vocabularyParser := range config.Vocabulary {
			parserNames = append(parserNames, vocabularyParser)
		}
		for _, deck := range config.Decks {
			if deck.Parser != "" && !slices.Contains(parserNames, deck.Parser) {
				sl.ReportError(deck.Parser, "parser", "Parser", "not found", "")
			}
		}
	}
}

func isParentDirectory(path, parent string) bool {
	pathParts := strings.Split(path, "/")
	parentParts := strings.Split(parent, "/")
	if len(parentParts) < len(pathParts) {
		return false
	}
	for i, part := range pathParts {
		if part != parentParts[i] {
			return false
		}
	}
	return true
}
