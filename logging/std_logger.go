package logging

type stdLogger struct {
	baseLogger
}

func (sl *stdLogger) Close() error {
	return nil
}
