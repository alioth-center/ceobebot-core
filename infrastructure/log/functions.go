package log

import "github.com/sirupsen/logrus"

// Debug 记录一个调试级别的日志
func Debug(message Fields) {
	logger.WithFields(logrus.Fields(message)).Debug()
}

// Info 记录一个信息级别的日志
func Info(message Fields) {
	logger.WithFields(logrus.Fields(message)).Info()
}

// Warn 记录一个警告级别的日志
func Warn(message Fields) {
	logger.WithFields(logrus.Fields(message)).Warn()
}

// Error 记录一个错误级别的日志
func Error(message Fields) {
	logger.WithFields(logrus.Fields(message)).Error()
}

// Fatal 记录一个致命错误级别的日志
func Fatal(message Fields) {
	logger.WithFields(logrus.Fields(message)).Fatal()
}

// Panic 记录一个崩溃级别的日志
func Panic(message Fields) {
	logger.WithFields(logrus.Fields(message)).Panic()
}

// Log 记录一个自定义级别的日志
func Log(level string, message Fields) {
	logger.WithFields(logrus.Fields(message)).Log(getLevel(level))
}
