package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/parser"
)

func Test_validateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		templates []api.Template
		want      string
	}{
		{
			"missing path",
			Config{Sync: []Sync{{Parser: "note"}}},
			nil,
			"want path to be defined",
		},
		{
			"both template and parser defined",
			Config{Sync: []Sync{{Path: "/", Parser: "note", Template: "xxx"}}},
			nil,
			"want only one of template or parser",
		},
		{
			"neither template or parser defined",
			Config{Sync: []Sync{{Path: "/"}}},
			nil,
			"want only one of template or parser",
		},
		{
			"duplicate sync paths",
			Config{Sync: []Sync{
				{Path: "/german", Parser: "note"},
				{Path: "/german/vocabulary", Parser: "note"},
				{Path: "/german", Parser: "note"},
			}},
			nil,
			"want no duplicates of sync paths, found \"/german\" 2 times",
		},
		{
			"invalid fields",
			Config{
				Templates: map[string]Template{
					"german": {
						Parser:     "vocabulary",
						TemplateID: "xxxxxxxx",
						Fields: map[string]string{
							"aaaaaaaa": "INVALID",
						},
					},
				},
				parsers: []parser.Parser{parser.NewNote(), parser.NewVocabulary()},
			},
			[]api.Template{
				{
					ID: "xxxxxxxx",
					Fields: map[string]api.FieldTemplate{
						"aaaaaaaa": {ID: "aaaaaaaa"},
					},
				},
			},
			"want fields to be one of [word examples notes] on template \"german\", got \"INVALID\"",
		},
		{
			"missing parser",
			Config{Templates: map[string]Template{
				"german": {TemplateID: "xxxxxxxx", Fields: map[string]string{}},
			}},
			[]api.Template{{ID: "xxxxxxxx"}},
			"want parser to be defined on template \"german\"",
		},
		{
			"invalid template id",
			Config{
				Templates: map[string]Template{
					"german": {Parser: "vocabulary", TemplateID: "xxxxxxxx", Fields: map[string]string{}},
				},
				parsers: []parser.Parser{parser.NewNote(), parser.NewVocabulary()},
			},
			[]api.Template{{ID: "yyyyyyyy"}},
			"want template id to be valid on template \"german\"",
		},
		{
			"invalid template field id",
			Config{
				Templates: map[string]Template{
					"german": {
						Parser:     "vocabulary",
						TemplateID: "xxxxxxxx",
						Fields: map[string]string{
							"INVALID": "word",
						},
					},
				},
				parsers: []parser.Parser{parser.NewNote(), parser.NewVocabulary()},
			},
			[]api.Template{
				{
					ID: "xxxxxxxx",
					Fields: map[string]api.FieldTemplate{
						"aaaaaaaa": {ID: "aaaaaaaa"},
					},
				},
			},
			"want template field id to be one of [aaaaaaaa] on template \"german\", got \"INVALID\"",
		},
		{
			"valid",
			Config{
				Sync: []Sync{
					{
						Path:     "german/vocabulary",
						Template: "german",
						Archive:  true,
					},
					{
						Path:   "german/grammar",
						Parser: "headings",
					},
					{
						Path:   "/",
						Name:   "Notes (root)",
						Parser: "note",
					},
				},
				Templates: map[string]Template{
					"german": {
						Parser:     "vocabulary",
						TemplateID: "xxxxxxxx",
						Fields: map[string]string{
							"aaaaaaaa": "word",
							"bbbbbbbb": "examples",
							"cccccccc": "notes",
						},
					},
				},
				parsers: []parser.Parser{parser.NewNote(), parser.NewVocabulary()},
			},
			[]api.Template{
				{
					ID: "xxxxxxxx",
					Fields: map[string]api.FieldTemplate{
						"aaaaaaaa": {ID: "aaaaaaaa"},
						"bbbbbbbb": {ID: "bbbbbbbb"},
						"cccccccc": {ID: "cccccccc"},
					},
				},
			},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config, tt.templates)

			if tt.want == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want)
			}
		})
	}
}
