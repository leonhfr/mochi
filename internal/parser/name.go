package parser

import (
	"path/filepath"
	"strings"
)

func getFilename(path string) string {
	return filepath.Base(path)
}

func getNameFromPath(path string) string {
	base := filepath.Base(path)
	for _, ext := range extensions {
		base = strings.TrimSuffix(base, ext)
	}
	return base
}
