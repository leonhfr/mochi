package cli

import (
	"fmt"
	"io"
)

type logger struct {
	w io.Writer
}

func newLogger(w io.Writer) *logger {
	return &logger{w}
}

func (l *logger) Debugf(format string, args ...any) {
	fmt.Fprintf(l.w, fmt.Sprintf("%s\n", format), args...)
}

func (l *logger) Errorf(format string, args ...any) {
	fmt.Fprintf(l.w, fmt.Sprintf("%s\n", format), args...)
}

func (l *logger) Infof(format string, args ...any) {
	fmt.Fprintf(l.w, fmt.Sprintf("%s\n", format), args...)
}
