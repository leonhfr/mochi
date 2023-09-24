package sync

type Logger interface {
	Error(message string)
	Errorf(format string, args ...any)
	Info(message string)
	Infof(format string, args ...any)
}
