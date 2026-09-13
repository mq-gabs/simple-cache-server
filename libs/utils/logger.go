package utils

import (
	"context"

	"go.uber.org/zap"
)

var loggerKey = "logger"

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) {
	context.WithValue(ctx, loggerKey, logger)
}

func FromCTX(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger)
	if !ok || logger == nil {
		return zap.NewNop().Sugar()
	}

	return logger
}
