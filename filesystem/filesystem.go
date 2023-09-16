package filesystem

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type Interface interface {
	FileExists(path string) bool
	Read(path string) ([]byte, error)
	Write(path, content string) error
	Sources(extensions []string) ([]string, error)
}

type Filesystem struct {
	workspace string
}

func New(workspace string) Interface {
	return &Filesystem{workspace: workspace}
}

// FileExists returns true if the file exists and is not a directory.
func (f *Filesystem) FileExists(path string) bool {
	info, err := os.Stat(f.fullPath(path))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Read returns the content of the file.
// If the file does not exist, no error is returned.
func (f *Filesystem) Read(path string) ([]byte, error) {
	bytes, err := os.ReadFile(f.fullPath(path))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (f *Filesystem) Write(path, content string) error {
	file, err := os.Create(f.fullPath(path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// Sources expects the extension with a dot: [".md"].
func (f *Filesystem) Sources(extensions []string) ([]string, error) {
	var sources []string
	err := filepath.WalkDir(f.workspace, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, f.workspace)
		base := filepath.Base(path)

		isDot := base != "." && base[0] == '.'
		isDir := d.IsDir()

		if isDir && isDot {
			return fs.SkipDir
		}

		if isDir || isDot {
			return nil
		}

		if ext := filepath.Ext(path); slices.Contains[[]string](extensions, ext) {
			sources = append(sources, path)
		}

		return nil
	})
	return sources, err
}

func (f *Filesystem) fullPath(path string) string {
	return filepath.Join(f.workspace, path)
}
