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

type Zapper = ILogger[zap.Field]

func newZapLogger(cfg *LoggerConfig) Zapper {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.CallerKey = "caller"
	encoderCfg.EncodeCaller = zapcore.FullCallerEncoder

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
		Encoding:          cfg.Encoding,
		EncoderConfig:     encoderCfg,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	logger := zap.Must(config.Build(zap.AddCallerSkip(1)))

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
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
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
