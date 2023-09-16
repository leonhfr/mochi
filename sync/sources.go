package sync

import "github.com/leonhfr/mochi/filesystem"

var extensions = []string{".md"}

func Sources(config Config, fs filesystem.Interface) ([]string, error) {
	all, err := fs.Sources(extensions)
	if err != nil {
		return nil, err
	}

	var sources []string
	for _, source := range all {
		if !config.ignored(source) {
			sources = append(sources, source)
		}
	}
	return sources, nil
}
