package sync

import (
	"context"
	"fmt"
	"path/filepath"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
)

const configName = "mochi"

var configExtensions = [2]string{"yaml", "yml"}

type Config struct {
	Sync      []Sync              `yaml:"sync"`
	Ignore    []string            `yaml:"ignore"`
	Templates map[string]Template `yaml:"templates"`

	parsers []parser.Parser
}

type Sync struct {
	Path     string `yaml:"path"`
	Name     string `yaml:"name"`
	Parser   string `yaml:"parser"`
	Template string `yaml:"template"`
	Archive  bool   `yaml:"archive"`
	Walk     bool   `yaml:"walk"`
}

type Template struct {
	Parser     string            `yaml:"parser"`
	TemplateID string            `yaml:"templateId"`
	Fields     map[string]string `yaml:"fields"`
}

func ReadConfig(ctx context.Context, parsers []parser.Parser, client Client, fs filesystem.Interface) (Config, error) {
	config := Config{parsers: parsers}
	path := configPath(fs)
	if path == "" {
		return config, nil
	}

	config, err := parseConfig(config, path, fs)
	if err != nil {
		return config, err
	}

	config = cleanConfig(config)

	templates, err := client.ListTemplates(ctx)
	if err != nil {
		return config, err
	}

	return config, validateConfig(config, templates)
}

func (c *Config) ignored(path string) bool {
	for _, pattern := range c.Ignore {
		ok, err := filepath.Match(pattern, path)
		if ok || err != nil {
			return true
		}
	}
	return false
}

func configPath(fs filesystem.Interface) string {
	for _, ext := range configExtensions {
		path := fmt.Sprintf("%s.%s", configName, ext)
		if fs.FileExists(path) {
			return path
		}
	}
	return ""
}

func parseConfig(config Config, path string, fs filesystem.Interface) (Config, error) {
	source, err := fs.Read(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func cleanConfig(config Config) Config {
	for i, s := range config.Sync {
		config.Sync[i].Path = filepath.Clean(filepath.Join("/", s.Path))
	}

	slices.SortFunc[[]Sync](config.Sync, func(a, b Sync) int {
		return len(b.Path) - len(a.Path)
	})

	for i, p := range config.Ignore {
		config.Ignore[i] = filepath.Clean(filepath.Join("/", p))
	}

	return config
}
