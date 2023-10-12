package mihoyo

import (
	"github.com/ceobebot/qqchannel/plugin"
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
	"strings"
)

type GenshinCommand struct {
	plugin.TextReplyMessageBaseImplementation
}

func (g GenshinCommand) Name() string {
	return "genshin"
}

func (g GenshinCommand) Description() string {
	return "原神相关命令"
}

func (g GenshinCommand) Example() string {
	return "/mihoyo gs [命令]\n/mihoyo ys [命令]"
}

func (g GenshinCommand) Triggered(content string) (triggered bool) {
	return strings.HasPrefix(content, "gs") || strings.HasPrefix(content, "ys")
}

func (g GenshinCommand) Handle(payload processor.Payload) (replyMessage message.Message) {
	return message.NewTextMessage().
		At(payload.Message.Author.ID).
		Emojis("[太阳]").
		Text("原神相关功能仍在开发适配中，敬请期待！")
}
