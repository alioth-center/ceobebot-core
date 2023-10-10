package mihoyo

import (
	"strings"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor/message"
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

func (s StarRailCommand) Handle(payload processor.Payload) (replyMessage message.Message) {
	return message.NewTextMessage().
		At(payload.Message.Author.ID).
		Emojis("[太阳]").
		Text("星铁相关功能仍在开发适配中，敬请期待！")
}
