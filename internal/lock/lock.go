package lock

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sync"
)

const lockName = "mochi-lock.json"

type lockData map[string]lockDeck // indexed by deck id

type lockDeck struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// Lock represents a lockfile.
type Lock struct {
	data    lockData
	path    string
	updated bool
	mu      sync.RWMutex
	rw      ReaderWriter
}

// ReaderWriter represents the interface to interact with a lockfile.
type ReaderWriter interface {
	Exists(string) bool
	Read(string) (io.ReadCloser, error)
	Write(string) (io.WriteCloser, error)
}

// Parse parses the lockfile in the target directory.
func Parse(rw ReaderWriter, target string) (*Lock, error) {
	path := filepath.Join(target, lockName)
	lock := &Lock{path: path, rw: rw}
	if !rw.Exists(path) {
		return lock, nil
	}

	r, err := rw.Read(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if err := json.NewDecoder(r).Decode(&lock.data); err != nil {
		return nil, err
	}

	return lock, nil
}

// String implements fmt.Stringer.
func (l *Lock) String() string {
	return fmt.Sprint(l.data)
}

// Write writes the lockfile to the target directory.
func (l *Lock) Write() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.updated {
		return nil
	}

	w, err := l.rw.Write(l.path)
	if err != nil {
		return err
	}
	defer w.Close()

	return json.NewEncoder(w).Encode(l.data)
}
