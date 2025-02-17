package logger

import (
	"context"
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

func (l *logrusLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Info(fields...)
}

func (l *logrusLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Error(fields...)
}

func (l *logrusLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Debug(fields...)
}

func (l *logrusLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Warn(fields...)
}

func (l *logrusLogger) Panic(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Panic(fields...)
}

func (l *logrusLogger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Fatal(fields...)
}

func (l *logrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l *logrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
func (l *logrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}
func (l *logrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}
func (l *logrusLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}
func (l *logrusLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) Sync() { /* Logrus doesn't require sync */ }
