package sync

type Logger interface {
	Info(message string)
	Infof(format string, args ...any)
}
