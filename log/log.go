package log

import "go.uber.org/zap"

// Slogger for logging
var Slogger *zap.SugaredLogger

// GetLogger function for initializing the log
func GetLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	Slogger = logger.Sugar()
	return logger
}
