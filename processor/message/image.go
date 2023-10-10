package message

import "github.com/tencent-connect/botgo/dto"

type ImageMessage struct {
	image     string
	reference *dto.MessageReference
}

func (m *ImageMessage) Image(url string) *ImageMessage {
	m.image = url
	return m
}

func (m *ImageMessage) Reference(messageID string) *ImageMessage {
	m.reference = &dto.MessageReference{
		MessageID:             messageID,
		IgnoreGetMessageError: true,
	}
	return m
}

func (m *ImageMessage) Build() *dto.MessageToCreate {
	return &dto.MessageToCreate{
		Image:            m.image,
		MessageReference: m.reference,
	}
}

func NewImageMessage() *ImageMessage {
	return &ImageMessage{
		image:     "",
		reference: nil,
	}
}
