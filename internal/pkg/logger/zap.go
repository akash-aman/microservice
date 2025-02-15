package logger

/**
 * https://www.youtube.com/watch?v=I2mWnh66Bkg
 *
 * Need fixes according to above.
 */
import (
	"context"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func newZapLogger(cfg *LoggerConfig) ILogger[zap.Field] {
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
		Encoding:          cfg.Encoding,
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

func (l *zapLogger) AddTraceAttribute(ctx context.Context, msg string, fields ...zap.Field) {

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		l.logger.Warn("No active span found in context")
		return
	}

	var otelAttributes []attribute.KeyValue
	for _, field := range fields {
		switch field.Type {
		case zapcore.StringType:
			otelAttributes = append(otelAttributes, attribute.String(field.Key, field.String))
		case zapcore.Int64Type:
			otelAttributes = append(otelAttributes, attribute.Int64(field.Key, field.Integer))
		case zapcore.Float64Type:
			otelAttributes = append(otelAttributes, attribute.Float64(field.Key, float64(field.Integer)))
		case zapcore.BoolType:
			otelAttributes = append(otelAttributes, attribute.Bool(field.Key, field.Integer == 1))
		default:
			otelAttributes = append(otelAttributes, attribute.String(field.Key, field.String))
		}
	}

	span.SetAttributes(otelAttributes...)
}

// Structured logging methods
func (l *zapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

func (l *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

func (l *zapLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Sugar().Panicf(format, args...)
}

func (l *zapLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

func (l *zapLogger) Sync() { l.logger.Sync() }
