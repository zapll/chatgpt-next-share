package main

import "github.com/daodao97/goadmin/pkg/logger"

var defaultLogger = logger.Default()

func init() {
	logger.SetCaller(false)
}

func Error(message string, content ...interface{}) {
	defaultLogger.Log(logger.LevelError, message, content...)
}

func Info(message string, content ...interface{}) {
	defaultLogger.Log(logger.LevelInfo, message, content...)
}

func Debug(message string, content ...interface{}) {
	defaultLogger.Log(logger.LevelDebug, message, content...)
}

func Warn(message string, content ...interface{}) {
	defaultLogger.Log(logger.LevelWarn, message, content...)
}
