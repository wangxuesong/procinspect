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

var logger, _ = zap.NewProduction(zap.WithCaller(false), zap.AddStacktrace(ErrorLevel))
