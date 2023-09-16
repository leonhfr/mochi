package sync

import (
	"fmt"
	"path/filepath"

	"github.com/leonhfr/mochi/parser"
)

type deckJob struct {
	sources     []string
	id          string
	archive     bool
	hasTemplate bool
	parser      parser.Parser
	template    Template
}

type jobMap map[string]*deckJob

func newJobMap(parsers []parser.Parser, sources []string, lock *Lock, config Config) (jobMap, error) {
	jobs := make(jobMap)
	for _, source := range sources {
		sync, ok := config.matchSync(source)
		if !ok {
			continue
		}

		path := filepath.Dir(source)
		if _, ok := jobs[path]; ok {
			jobs[path].sources = append(jobs[path].sources, source)
		} else {
			job, err := newJob(parsers, path, source, sync, lock, config)
			if err != nil {
				return jobs, err
			}
			jobs[path] = job
		}
	}
	return jobs, nil
}

func newJob(parsers []parser.Parser, path, source string, sync Sync, lock *Lock, config Config) (*deckJob, error) {
	lock.mu.RLock()
	defer lock.mu.RUnlock()

	deck, ok := lock.Decks[path]
	if !ok {
		return nil, fmt.Errorf("deck id of path %s not found", path)
	}

	template, hasTemplate := config.Templates[sync.Template]
	parser, err := getParser(parsers, template.Parser, sync.Parser)
	if err != nil {
		return nil, err
	}

	return &deckJob{
		sources:     []string{source},
		id:          deck[indexDeckID],
		archive:     sync.Archive,
		hasTemplate: hasTemplate,
		parser:      parser,
		template:    template,
	}, nil
}

func getParser(parsers []parser.Parser, names ...string) (parser.Parser, error) {
	for _, name := range names {
		for _, parser := range parsers {
			if parser.String() == name {
				return parser, nil
			}
		}
	}
	return nil, fmt.Errorf("parsers %v not found, want one of %v", names, parsers)
}
