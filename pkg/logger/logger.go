package logger

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	LevelDebug Level = zapcore.DebugLevel
	LevelInfo  Level = zapcore.InfoLevel
	LevelWarn  Level = zapcore.WarnLevel
	LevelError Level = zapcore.ErrorLevel
	LevelFatal Level = zapcore.FatalLevel
	LevelPanic Level = zapcore.PanicLevel
)

type Fields map[string]interface{}

type Logger struct {
	l     *zap.Logger
	sugar *zap.SugaredLogger
}

func NewLogger(writer io.Writer, prefix string, flag int) *Logger {
	writeSyncer := zapcore.AddSync(writer)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	logger := zap.New(core, zap.AddCaller())
	sugar := logger.Sugar()

	return &Logger{
		l:     logger,
		sugar: sugar,
	}
}

func (l *Logger) clone() *Logger {
	return &Logger{
		l:     l.l,
		sugar: l.sugar,
	}
}

func (l *Logger) WithTrace() *Logger {
	return l
}

func (l *Logger) WithFields(f Fields) *Logger {
	fields := make([]zap.Field, 0, len(f))
	for k, v := range f {
		fields = append(fields, zap.Any(k, v))
	}
	nl := l.clone()
	nl.l = nl.l.With(fields...)
	nl.sugar = nl.l.Sugar()
	return nl
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	nl := l.clone()
	nl.l = nl.l.With(zap.Any("context", ctx))
	nl.sugar = nl.l.Sugar()
	return nl
}

func (l *Logger) WithCaller(skip int) *Logger {
	nl := l.clone()
	nl.l = nl.l.WithOptions(zap.AddCallerSkip(skip))
	nl.sugar = nl.l.Sugar()
	return nl
}

func (l *Logger) WithCallersFrames() *Logger {
	// zap already includes caller information, so this is a no-op
	return l
}

func (l *Logger) Output(level Level, message string) {
	switch level {
	case LevelDebug:
		l.l.Debug(message)
	case LevelInfo:
		l.l.Info(message)
	case LevelWarn:
		l.l.Warn(message)
	case LevelError:
		l.l.Error(message)
	case LevelFatal:
		l.l.Fatal(message)
	case LevelPanic:
		l.l.Panic(message)
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.sugar.Info(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.sugar.Infof(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.sugar.Fatal(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.sugar.Fatalf(format, v...)
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).sugar.Errorf(format, v...)
}

// 新增的方法，用于支持trace功能
func (l *Logger) WithTraceFromGin(ctx *gin.Context) *Logger {
	traceID, _ := ctx.Get("X-Trace-ID")
	spanID, _ := ctx.Get("X-Span-ID")
	return l.WithFields(Fields{
		"trace_id": traceID,
		"span_id":  spanID,
	})
}

// 如果需要访问底层的zap.Logger
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.l
}

// 关闭日志，确保所有日志都被写入
func (l *Logger) Sync() error {
	return l.l.Sync()
}
