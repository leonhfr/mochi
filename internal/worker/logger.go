package worker

// Logger is the interface to log output.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}
