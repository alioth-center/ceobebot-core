package log

// Debug 记录一个调试级别的日志
func Debug(message Fields) {
	logger.Debug(message)
}

// Info 记录一个信息级别的日志
func Info(message Fields) {
	logger.Info(message)
}

// Warn 记录一个警告级别的日志
func Warn(message Fields) {
	logger.Warn(message)
}

// Error 记录一个错误级别的日志
func Error(message Fields) {
	logger.Error(message)
}

// Fatal 记录一个致命错误级别的日志
func Fatal(message Fields) {
	logger.Fatal(message)
}

// Panic 记录一个崩溃级别的日志
func Panic(message Fields) {
	logger.Panic(message)
}

// Log 记录一个自定义级别的日志
func Log(level string, message Fields) {
	logger.Log(level, message)
}
