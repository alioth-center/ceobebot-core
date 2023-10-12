package message

import (
	"github.com/tencent-connect/botgo/dto"
	"strings"
)

type TextMessage struct {
	content   strings.Builder
	image     string
	reference *dto.MessageReference
}

func (m *TextMessage) At(userID string) *TextMessage {
	m.content.WriteString("<@!")
	m.content.WriteString(userID)
	m.content.WriteString(">")
	return m
}

func (m *TextMessage) AtAll() *TextMessage {
	m.content.WriteString("@everyone")
	return m
}

func (m *TextMessage) Text(text string) *TextMessage {
	m.content.WriteString(text)
	return m
}

func (m *TextMessage) NewLine() *TextMessage {
	m.content.WriteString("\n")
	return m
}

func (m *TextMessage) Emojis(emojis ...string) *TextMessage {
	for _, emoji := range emojis {
		if strings.HasPrefix(emoji, "[") {
			emoji = emoji[1:]
		}
		if strings.HasSuffix(emoji, "]") {
			emoji = emoji[:len(emoji)-1]
		}

		emojiID, existEmoji := GetEmojiFormKeyword(emoji)
		if existEmoji {
			// 如果有 emoji，使用 emoji 表情，格式为 <emoji:emojiID>
			m.content.WriteString("<emoji:")
			m.content.WriteString(emojiID)
			m.content.WriteString(">")
		} else {
			// 如果没有 emoji，使用中括号包裹，格式为 [emoji]
			m.content.WriteString("[")
			m.content.WriteString(emoji)
			m.content.WriteString("]")
		}
	}
	return m
}

func (m *TextMessage) Image(url string) *TextMessage {
	m.image = url
	return m
}

func (m *TextMessage) Reference(messageID string) *TextMessage {
	m.reference = &dto.MessageReference{
		MessageID:             messageID,
		IgnoreGetMessageError: true,
	}
	return m
}

func (m *TextMessage) Build() *dto.MessageToCreate {
	return &dto.MessageToCreate{
		Content:          m.content.String(),
		Image:            m.image,
		MessageReference: m.reference,
	}
}

func (m *TextMessage) Type() Type {
	return TextMessageType
}
