package logger

type stdLogger struct {
	baseLogger
}

func (sl *stdLogger) Close() error {
	return nil
}
