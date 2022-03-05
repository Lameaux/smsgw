package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugaredLogger *zap.SugaredLogger //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	log, err := logConfig.Build()
	if err != nil {
		panic(err)
	}

	sugaredLogger = log.Sugar()
}

func Log() *zap.SugaredLogger {
	return sugaredLogger
}

func Infow(msg string, args ...interface{}) {
	Log().Infow(msg, args...)
}

func Errorw(msg string, args ...interface{}) {
	Log().Errorw(msg, args...)
}

func Error(args ...interface{}) {
	Log().Error(args...)
}

func Fatal(args ...interface{}) {
	Log().Fatal(args...)
}

func Fatalw(msg string, args ...interface{}) {
	Log().Fatalw(msg, args...)
}
