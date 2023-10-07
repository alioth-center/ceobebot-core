package mihoyo

import (
	"strings"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
)

type StarRailCommand struct {
	plugin.TextReplyMessageBaseImplementation
}

func (s StarRailCommand) Name() string {
	return "star-rail"
}

func (s StarRailCommand) Description() string {
	return "星穹铁道相关指令"
}

func (s StarRailCommand) Example() string {
	return "/mihoyo sr [命令]"
}

func (s StarRailCommand) Triggered(content string) (triggered bool) {
	return strings.HasPrefix(content, "sr")
}

func (s StarRailCommand) Handle(payload processor.Payload) (replyMessage string) {
	return "适配中，敬请期待"
}
