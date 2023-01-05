package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
}

// ChangeLogLevel accepts log level case insensitive.
// options - "panic", "fatal", "error", "warn", "warning", "info", "debug", "trace"
func ChangeLogLevel(level string) error {
	logrusLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		return err
	}

	logrus.SetLevel(logrusLevel)
	return nil
}
