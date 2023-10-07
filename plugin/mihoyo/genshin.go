package mihoyo

import (
	"strings"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
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
	return "/mihoyo gs [命令]"
}

func (g GenshinCommand) Triggered(content string) (triggered bool) {
	return strings.HasPrefix(content, "gs")
}

func (g GenshinCommand) Handle(payload processor.Payload) (replyMessage string) {
	return "适配中，敬请期待"
}
