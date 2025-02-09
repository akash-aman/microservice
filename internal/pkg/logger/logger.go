package logger

import "log"

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
	LogLevel    string     `mapstructure:"level" validate:"required"`
	Type        LoggerType `mapstructure:"type" validate:"required"`
	FileLogging bool       `mapstructure:"fileLogging" validate:"required"`
	AccessLog   string     `mapstructure:"accessLog" validate:"required"`
	ErrorLog    string     `mapstructure:"errorLog" validate:"required"`
}

/**
 * ILogger interface defines methods that both implementations must satisfy
 */
type ILogger interface {
	/**
	 * Structured logging methods (key-value pairs)
	 */
	Info(fields ...interface{})
	Error(fields ...interface{})
	Debug(fields ...interface{})
	Warn(fields ...interface{})
	Panic(fields ...interface{})
	Fatal(fields ...interface{})

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
func InitLogger(cfg *LoggerConfig) ILogger {
	if cfg == nil {
		log.Fatal("LoggerConfig is nil. Ensure it is properly initialized.")
	}

	switch cfg.Type {
	case LogrusLogger:
		return newLogrusLogger(cfg)
	default:
		return newZapLogger(cfg)
	}
}
