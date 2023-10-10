package message

import "github.com/tencent-connect/botgo/dto"

type Message interface {
	Build() *dto.MessageToCreate
}
