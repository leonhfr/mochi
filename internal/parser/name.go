package parser

import (
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getNameFromPath(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), Extension)
	return cases.Title(language.Und).String(base)
}
