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

type MessagePlugin interface {
	TriggerKey() string
	Commands() []Command
}

type Command interface {
	Name() string
	Description() string
	Example() string
	Triggered(content string) (triggered bool)
	Type() MessageCommandType
	Handle(payload processor.Payload) (replyMessage message.Message)
}

type TextReplyCommand Command

type ImageReplyCommand Command
