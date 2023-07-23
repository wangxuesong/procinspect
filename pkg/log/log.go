package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Level = zapcore.Level
	Field = zap.Field

	Logger struct {
		l     *zap.Logger
		level Level

		config *zap.Config
	}
)

const (
	InfoLevel  Level = zap.InfoLevel
	WarnLevel  Level = zap.WarnLevel
	ErrorLevel Level = zap.ErrorLevel
	DebugLevel Level = zap.DebugLevel
)

var (
	Skip    = zap.Skip
	Binary  = zap.Binary
	Bool    = zap.Bool
	Int     = zap.Int
	Int64   = zap.Int64
	Float64 = zap.Float64
	String  = zap.String
	Any     = zap.Any

	Info  = logger.Info
	Warn  = logger.Warn
	Error = logger.Error
	Debug = logger.Debug
	Sync  = logger.Sync
)

var logger, _ = DefaultLogger()

func DefaultLogger() (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	cfg.Level = zap.NewAtomicLevelAt(InfoLevel)
	log, err := cfg.Build(zap.WithCaller(false), zap.AddStacktrace(ErrorLevel))
	if err != nil {
		return nil, err
	}
	return &Logger{l: log, level: InfoLevel, config: &cfg}, nil
}

func CurrentLogger() *Logger {
	return logger
}

func (log *Logger) Info(msg string, fields ...Field) {
	log.l.Info(msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...Field) {
	log.l.Warn(msg, fields...)
}

func (log *Logger) Error(msg string, fields ...Field) {
	log.l.Error(msg, fields...)
}

func (log *Logger) Debug(msg string, fields ...Field) {
	log.l.Debug(msg, fields...)
}

func (log *Logger) Sync() {
	log.l.Sync()
}

func (log *Logger) SetLevel(level Level) {
	log.level = level
	log.config.Level.SetLevel(level)
}

func SetLevel(level Level) {
	logger.SetLevel(level)
}
