package plugin

import (
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
)

type MessageCommandType string

const (
	TextReplyMessageCommandType  MessageCommandType = "TextReplyMessageCommandType"
	ImageReplyMessageCommandType MessageCommandType = "ImageReplyMessageCommandType"
)

// MessagePluginInfo 消息插件信息
type MessagePluginInfo struct {
	Name       string
	TriggerKey string
	BlockChain bool
}

// MessageCommandInfo 消息命令信息
type MessageCommandInfo struct {
	Name        string
	Description string
	Type        MessageCommandType
	Example     []string
}

type MessagePlugin interface {
	TriggerKey() string
	Commands() []MessageCommand
}

type MessageCommand interface {
	Name() string
	Description() string
	Example() string
	Triggered(content string) (triggered bool)
	Type() MessageCommandType
	Handle(payload processor.Payload) (replyMessage message.Message)
}

type TextReplyCommand MessageCommand

type ImageReplyCommand MessageCommand
