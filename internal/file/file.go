package file

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// System represents an interface with the filesystem.
type System struct{}

// NewSystem returns a new System.
func NewSystem() *System { return &System{} }

// Walk walks recursively in workspace and calls the callback
// for all files that match one of the extension and is not hidden
// (does not start with a dot ".").
//
// The function expects the extensions with a dot: [".md"].
func (System) Walk(workspace string, extensions []string, cb func(string)) error {
	return filepath.WalkDir(workspace, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, workspace)
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
			cb(path)
		}

		return nil
	})
}

// Read returns an io.ReadCloser from the file at path.
//
// If file not exists, it returns fs.ErrNotExist.
func (System) Read(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return file, fs.ErrNotExist
	}
	return file, err
}

// Write returns an io.WriteCloser to the file at path.
func (System) Write(path string) (io.WriteCloser, error) {
	file, err := os.Create(path)
	return file, err
}
