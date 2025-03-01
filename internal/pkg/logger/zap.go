package logger

/**
 * https://www.youtube.com/watch?v=I2mWnh66Bkg
 * https://betterstack.com/community/guides/logging/go/zap/
 * https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelzap
 * https://uptrace.dev/get/opentelemetry-go#exporting-logs
 * Need fixes according to above.
 */
import (
	"context"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *otelzap.Logger
}

type Zapper = ILogger[zap.Field]

func newZapLogger(cfg *LoggerConfig, provider *sdklog.LoggerProvider) Zapper {
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

	otelLogger := otelzap.New(
		logger,
		otelzap.WithLoggerProvider(provider),
	)
	return &zapLogger{logger: otelLogger}
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
func (l *zapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Ctx(ctx).Info(msg, fields...)
}

func (l *zapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Ctx(ctx).Error(msg, fields...)
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Ctx(ctx).Debug(msg, fields...)
}

func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Ctx(ctx).Warn(msg, fields...)
}

func (l *zapLogger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Ctx(ctx).Panic(msg, fields...)
}

func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Ctx(ctx).Fatal(msg, fields...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Ctx(ctx).Sugar().Infof(format, args...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Ctx(ctx).Sugar().Errorf(format, args...)
}

func (l *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Ctx(ctx).Sugar().Debugf(format, args...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Ctx(ctx).Sugar().Warnf(format, args...)
}

func (l *zapLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Ctx(ctx).Sugar().Panicf(format, args...)
}

func (l *zapLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Ctx(ctx).Sugar().Fatalf(format, args...)
}

func (l *zapLogger) Sync() { l.logger.Sync() }
