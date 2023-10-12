package message

import (
	"github.com/tencent-connect/botgo/dto"
	"strings"
)

type Type string

const (
	TextMessageType     Type = "text"
	ImageMessageType    Type = "image"
	MarkdownMessageType Type = "markdown"
)

type Message interface {
	Build() *dto.MessageToCreate
	Type() Type
}

func NewTextMessage() *TextMessage {
	return &TextMessage{
		content:   strings.Builder{},
		image:     "",
		reference: nil,
	}
}

func NewImageMessage() *ImageMessage {
	return &ImageMessage{
		image:     "",
		reference: nil,
	}
}
