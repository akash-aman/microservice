package logger

/**
 * https://www.youtube.com/watch?v=I2mWnh66Bkg
 * 
 * Need fixes according to above.
 */
import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func newZapLogger(cfg *LoggerConfig) ILogger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	outputPaths := []string{"stderr"}
	errorOutputPaths := []string{"stderr"}
	if cfg.FileLogging {
		outputPaths = append(outputPaths, cfg.AccessLog)
		errorOutputPaths = append(errorOutputPaths, cfg.ErrorLog)
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(getZapLevel(cfg.LogLevel)),
		Development:       os.Getenv("APP_ENV") != "production",
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	logger := zap.Must(config.Build())

	return &zapLogger{logger: logger}
}

func getZapLevel(level string) zapcore.Level {
	levels := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
	}
	if l, ok := levels[level]; ok {
		return l
	}
	return zapcore.InfoLevel
}

// Structured logging methods
func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info("", zap.Any("value", args))
}

func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Error("", zap.Any("value", args))
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Debug("", zap.Any("value", args))
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.logger.Warn("", zap.Any("value", args))
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.logger.Panic("", zap.Any("value", args))
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal("", zap.Any("value", args))
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}
func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.logger.Sugar().Panicf(format, args...)
}
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

func (l *zapLogger) Sync() { l.logger.Sync() }
