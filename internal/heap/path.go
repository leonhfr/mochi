package heap

import (
	"path/filepath"
	"strings"
)

// Path represents a file path.
type Path string

var _ Item = Path("")

// Base implements the PriorityItem interface.
func (p Path) Base() string {
	return filepath.Dir(string(p))
}

// Priority implements the PriorityItem interface.
func (p Path) Priority() int {
	base := p.Base()
	if base == "/" {
		return 0
	}
	return strings.Count(base, "/")
}

// ConvertPaths concerts a slice of paths back to string.
func ConvertPaths(items []Path) []string {
	paths := make([]string, len(items))
	for i, item := range items {
		paths[i] = string(item)
	}
	return paths
}
