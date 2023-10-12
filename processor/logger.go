package processor

import "github.com/ceobebot/qqchannel/infrastructure/log"

var (
	logger = log.NewLogger(log.Config{
		Level:     "info",
		Formatter: "json",
		FilePath:  "data/messages.log",
	})
)
