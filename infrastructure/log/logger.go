package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/config"
)

var (
	logConfig Config
	logger    *Logger
)

func init() {
	// 从配置文件中加载日志配置
	if loadConfigErr := config.LoadCustomConfigWithKeys(&logConfig, "infrastructure", "logger"); loadConfigErr != nil {
		panic(loadConfigErr)
	}

	// 初始化日志记录器
	logger = NewLogger(logConfig)
}

// getLevel 获取日志等级
func getLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

// getFormatter 获取日志格式化器
func getFormatter(formatter string) logrus.Formatter {
	switch formatter {
	case "json":
		return &jsonFormatter{original: &logrus.JSONFormatter{}}
	case "text":
		return &logrus.TextFormatter{}
	default:
		return &jsonFormatter{original: &logrus.JSONFormatter{}}
	}
}

type Logger struct {
	logger *logrus.Logger
}

// Debug 记录一个调试级别的日志
func (l *Logger) Debug(message Fields) {
	l.logger.WithFields(logrus.Fields(message)).Debug()
}

// Info 记录一个信息级别的日志
func (l *Logger) Info(message Fields) {
	l.logger.WithFields(logrus.Fields(message)).Info()
}

// Warn 记录一个警告级别的日志
func (l *Logger) Warn(message Fields) {
	l.logger.WithFields(logrus.Fields(message)).Warn()
}

// Error 记录一个错误级别的日志
func (l *Logger) Error(message Fields) {
	l.logger.WithFields(logrus.Fields(message)).Error()
}

// Fatal 记录一个致命错误级别的日志
func (l *Logger) Fatal(message Fields) {
	l.logger.WithFields(logrus.Fields(message)).Fatal()
}

// Panic 记录一个崩溃级别的日志
func (l *Logger) Panic(message Fields) {
	l.logger.WithFields(logrus.Fields(message)).Panic()
}

// Log 记录一个自定义级别的日志
func (l *Logger) Log(level string, message Fields) {
	l.logger.WithFields(logrus.Fields(message)).Log(getLevel(level))
}

func NewLogger(conf Config) *Logger {
	// 初始化日志输出文件
	var output *os.File
	outputFile, initOutputFileErr := os.OpenFile(conf.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o755)
	if initOutputFileErr != nil {
		panic(initOutputFileErr)
	} else {
		output = outputFile
	}

	// 初始化日志记录器
	l := &Logger{logger: logrus.New()}
	l.logger.SetLevel(getLevel(conf.Level))
	l.logger.SetFormatter(getFormatter(conf.Formatter))
	l.logger.SetOutput(output)

	return l
}
