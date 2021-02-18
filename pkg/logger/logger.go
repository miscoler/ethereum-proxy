package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
}

type loggerImpl zap.Logger

func New() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, errors.Wrap(err, "creating logger")
	}

	return (*loggerImpl)(logger), nil
}

func (l *loggerImpl) Info(msg string, fields ...zapcore.Field) {
	(*zap.Logger)(l).Info(msg, fields...)
}

func (l *loggerImpl) Error(msg string, fields ...zapcore.Field) {
	(*zap.Logger)(l).Error(msg, fields...)
}

func (l *loggerImpl) Warn(msg string, fields ...zapcore.Field) {
	(*zap.Logger)(l).Warn(msg, fields...)
}

func (l *loggerImpl) Fatal(msg string, fields ...zapcore.Field) {
	(*zap.Logger)(l).Fatal(msg, fields...)
}

func (l *loggerImpl) With(fields ...zapcore.Field) Logger {
	return (*loggerImpl)((*zap.Logger)(l).With(fields...))
}
