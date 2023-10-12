package chat

import (
	"strings"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor/message"
)

type ManagementCommand struct {
	plugin.TextReplyMessageBaseImplementation
}

func (m ManagementCommand) Name() string {
	return "ChatManagement"
}

func (m ManagementCommand) Description() string {
	return "管理Chat插件"
}

func (m ManagementCommand) Example() string {
	return "/chat mgnt [子命令]"
}

func (m ManagementCommand) Triggered(content string) (triggered bool) {
	return strings.HasPrefix(content, "mgnt ")
}

func (m ManagementCommand) Handle(payload processor.Payload) (replyMessage message.Message) {
	return nil
}
