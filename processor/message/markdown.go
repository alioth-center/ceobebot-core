package message

import (
	"github.com/tencent-connect/botgo/dto"
	"strings"
)

type MarkdownMessage struct {
	content   strings.Builder
	reference *dto.MessageReference
}

func (m *MarkdownMessage) Text(text string) *MarkdownMessage {
	m.content.WriteString(text)
	return m
}

func (m *MarkdownMessage) Reference(messageID string) *MarkdownMessage {
	m.reference = &dto.MessageReference{
		MessageID:             messageID,
		IgnoreGetMessageError: true,
	}
	return m
}

func (m *MarkdownMessage) Build() *dto.MessageToCreate {
	return &dto.MessageToCreate{
		Content:          m.content.String(),
		MessageReference: m.reference,
	}
}
