package plugin

import "studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"

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
	Handle(payload processor.Payload) (replyMessage string)
}

type TextReplyCommand Command

type ImageReplyCommand Command
