package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Logger
}

func newLogrusLogger(cfg *LoggerConfig) ILogger[interface{}] {
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	logger.SetLevel(getLogrusLevel(cfg.LogLevel))

	return &logrusLogger{logger: logger}
}

func getLogrusLevel(level string) logrus.Level {
	levels := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
	}
	if l, ok := levels[level]; ok {
		return l
	}
	return logrus.InfoLevel
}

func (l *logrusLogger) Info(msg string, fields ...interface{}) {
	l.logger.Info(fields...)
}

func (l *logrusLogger) Error(msg string, fields ...interface{}) {
	l.logger.Error(fields...)
}

func (l *logrusLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Debug(fields...)
}

func (l *logrusLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Warn(fields...)
}

func (l *logrusLogger) Panic(msg string, fields ...interface{}) {
	l.logger.Panic(fields...)
}

func (l *logrusLogger) Fatal(msg string, fields ...interface{}) {
	l.logger.Fatal(fields...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}
func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) Sync() { /* Logrus doesn't require sync */ }
