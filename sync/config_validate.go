package sync

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/parser"
)

func validateConfig(config Config, templates []api.Template) error {
	var errs []string
	syncPaths := make(map[string]int)

	for _, s := range config.Sync {
		validate(&errs, len(s.Path) > 0, "want path to be defined")
		oneof := (len(s.Parser) == 0 && len(s.Template) > 0) || (len(s.Parser) > 0 && len(s.Template) == 0)
		validate(&errs, oneof, "want only one of template or parser")
		syncPaths[s.Path]++
	}

	for path, n := range syncPaths {
		validate(&errs, n == 1, "want no duplicates of sync paths, found \"%s\" %d times", path, n)
	}

	for _, p := range config.Ignore {
		_, err := filepath.Match("", p)
		validate(&errs, err == nil, "malformed ignore pattern \"%s\"", p)
	}

	for name, template := range config.Templates {
		validateTemplateConfig(&errs, name, template, config.parsers)
		validateAPITemplateConfig(&errs, name, template, templates, config.parsers)
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func validateTemplateConfig(errs *[]string, name string, template Template, parsers []parser.Parser) {
	validParsers := make([]string, 0, len(parsers))
	validFields := make(map[string][]string)
	for _, parser := range parsers {
		validParsers = append(validParsers, parser.String())
		validFields[parser.String()] = parser.Fields()
	}

	parserDefined := len(template.Parser) > 0
	validate(errs, parserDefined, "want parser to be defined on template \"%s\"", name)
	if !parserDefined {
		return
	}

	validParser := slices.Contains[[]string](validParsers, template.Parser)
	validate(errs, validParser, "want parser to be one of %v on template \"%s\", got \"%s\"", parsers, name, template.Parser)
	if !validParser {
		return
	}

	fields := validFields[template.Parser]
	for _, field := range template.Fields {
		validField := slices.Contains[[]string](fields, field)
		validate(errs, validField, "want fields to be one of %v on template \"%s\", got \"%s\"", fields, name, field)
	}
}

func validateAPITemplateConfig(errs *[]string, name string, template Template, apiTemplates []api.Template, parsers []parser.Parser) {
	templateDefined := len(template.TemplateID) > 0
	validate(errs, templateDefined, "want template id to be defined on template \"%s\"", name)
	if !templateDefined {
		return
	}

	index := slices.IndexFunc[[]api.Template](apiTemplates, func(t api.Template) bool {
		return template.TemplateID == t.ID
	})
	validTemplateID := index >= 0
	validate(errs, validTemplateID, "want template id to be valid on template \"%s\"", name)
	if !validTemplateID {
		return
	}

	validFieldIds := make([]string, 0, len(apiTemplates[index].Fields))
	for _, field := range apiTemplates[index].Fields {
		validFieldIds = append(validFieldIds, field.ID)
	}

	for _, parser := range parsers {
		if template.Parser != parser.String() {
			continue
		}

		for id := range template.Fields {
			validate(errs, slices.Contains[[]string](validFieldIds, id), "want template field id to be one of %v on template \"%s\", got \"%s\"", validFieldIds, name, id)
		}
	}
}

func validate(errs *[]string, ok bool, msg string, args ...any) {
	if !ok {
		*errs = append(*errs, fmt.Sprintf(msg, args...))
	}
}
