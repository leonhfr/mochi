package filesystem

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type Interface interface {
	FileExists(path string) bool
	Read(path string) ([]byte, error)
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

func (f *Filesystem) fullPath(path string) string {
	return filepath.Join(f.workspace, path)
}
