package config

import "os"

type Logger struct {
	LogFile string
}

func (l *Logger) FillFromEnvs() Configurable {
	l.LogFile = os.Getenv("LOG_FILE")

	return l
}

func (l *Logger) Validate() (Configurable, error) {
	return l, nil
}
