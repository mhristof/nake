package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// SetOutput Set log output.
func SetOutput(w io.Writer) {
	logger.Out = w
}

// WithField Print with fields.
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// Info print.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Panic print.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Fatal print.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// WithFields print.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

// WithError print.
func WithError(err error) *logrus.Entry {
	return logger.WithField("error", err)
}

// SetLevel Set the debug level.
func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

// Fields The fields.
type Fields = logrus.Fields

// DebugLevel Debug level.
var DebugLevel = logrus.DebugLevel
