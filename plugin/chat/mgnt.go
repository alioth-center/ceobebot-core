package chat

import (
	"github.com/ceobebot/qqchannel/plugin"
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
	"strings"
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
