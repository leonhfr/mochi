package parser

import (
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getFilename(path string) string {
	return filepath.Base(path)
}

func getNameFromPath(path string) string {
	base := filepath.Base(path)
	for _, ext := range extensions {
		base = strings.TrimSuffix(base, ext)
	}
	return cases.Title(language.Und).String(base)
}
