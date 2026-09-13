package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(string)
	Warning(string)
	Error(string)
	Fatal(string)
}

var defaultConfig zap.Config = zap.Config{
	Level:    zap.NewAtomicLevel(),
	Encoding: "json",
	OutputPaths: []string{
		"stdout",
		// "scas.log",
	},
	ErrorOutputPaths: []string{"stderr"},
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:    "ts",
		LevelKey:   "lvl",
		MessageKey: "msg",
		NameKey:    "scp",

		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalLevelEncoder,
	},
}

func New() *zap.SugaredLogger {
	logger, err := defaultConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
