package sync

type testLogger struct{}

func (l testLogger) Error(_ string)            {}
func (l testLogger) Errorf(_ string, _ ...any) {}
func (l testLogger) Info(_ string)             {}
func (l testLogger) Infof(_ string, _ ...any)  {}

var _ Logger = testLogger{}
