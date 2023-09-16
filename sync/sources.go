package sync

import (
	"path/filepath"

	"golang.org/x/exp/slices"

	"github.com/leonhfr/mochi/filesystem"
)

var extensions = []string{".md"}

func Sources(changed []string, config Config, fs filesystem.Interface) ([]string, error) {
	sources, err := fs.Sources(extensions)
	if err != nil {
		return nil, err
	}

	if len(changed) > 0 {
		sources = filterUnchanged(sources, changed)
	}
	sources = filterIgnored(sources, config)

	return sources, nil
}

func filterUnchanged(sources, changed []string) []string {
	dirMap := make(map[string]struct{})
	for _, path := range changed {
		dirMap[filepath.Dir(path)] = struct{}{}
	}

	var filtered []string
	for _, source := range sources {
		if _, ok := dirMap[filepath.Dir(source)]; ok {
			filtered = append(filtered, source)
		}
	}
	return filtered
}

func filterIgnored(sources []string, config Config) []string {
	var filtered []string
	for _, source := range sources {
		if !config.ignored(source) {
			filtered = append(filtered, source)
		}
	}
	return filtered
}

func uniqueDirs(sources []string) []string {
	dirMap := make(map[string]struct{})
	for _, path := range sources {
		for len(path) > 1 {
			path = filepath.Dir(path)
			dirMap[path] = struct{}{}
		}
	}
	dirs := make([]string, 0, len(dirMap))
	for dir := range dirMap {
		dirs = append(dirs, dir)
	}
	slices.Sort[[]string](dirs)
	return dirs
}
