package processor

import "studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/log"

var (
	logger = log.NewLogger(log.Config{
		Level:     "info",
		Formatter: "json",
		FilePath:  "data/messages.log",
	})
)
