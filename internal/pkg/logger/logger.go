package logger

import (
	"log"

	"go.uber.org/zap"
)

/**
 * LoggerType represents the type of logger to use
 */
type LoggerType string

const (
	ZapLogger    LoggerType = "zap"
	LogrusLogger LoggerType = "logrus"
)

/**
 * LoggerConfig holds configuration for the logger
 */
type LoggerConfig struct {
	LogLevel    string `mapstructure:"level" validate:"required"`
	FileLogging bool   `mapstructure:"fileLogging" validate:"required"`
	AccessLog   string `mapstructure:"accessLog" validate:"required"`
	ErrorLog    string `mapstructure:"errorLog" validate:"required"`
	Encoding    string `mapstructure:"encoding" validate:"required"`
}

/**
 * ILogger interface defines methods that both implementations must satisfy
 */
type ILogger[T any] interface {
	/**
	 * Structured logging methods (key-value pairs)
	 */
	Info(msg string, fields ...T)
	Error(msg string, fields ...T)
	Debug(msg string, fields ...T)
	Warn(msg string, fields ...T)
	Panic(msg string, fields ...T)
	Fatal(msg string, fields ...T)

	/**
	 * Format logging methods (printf style)
	 */
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Sync()
}

/**
 * Logger factory function
 */
func InitLogger[T any](cfg *LoggerConfig) ILogger[T] {
	if cfg == nil {
		log.Fatal("LoggerConfig is nil. Ensure it is properly initialized.")
	}

	var sample T

	switch any(sample).(type) {
	case zap.Field:
		return any(newZapLogger(cfg)).(ILogger[T])
	case interface{}:
		return any(newLogrusLogger(cfg)).(ILogger[T])
	}
	return nil
}
