package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/config"
)

var (
	logConfig Config
	output    *os.File
	logger    = logrus.New()
)

func init() {
	// 从配置文件中加载日志配置
	if loadConfigErr := config.LoadCustomConfigWithKeys(&logConfig, "infrastructure", "logger"); loadConfigErr != nil {
		panic(loadConfigErr)
	}

	// 初始化日志输出文件
	outputFile, initOutputFileErr := os.OpenFile(logConfig.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o755)
	if initOutputFileErr != nil {
		panic(initOutputFileErr)
	} else {
		output = outputFile
	}

	// 初始化日志记录器
	logger.SetLevel(getLevel(logConfig.Level))
	logger.SetFormatter(getFormatter(logConfig.Formatter))
	logger.SetOutput(output)
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
		return &logrus.JSONFormatter{}
	case "text":
		return &logrus.TextFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}
