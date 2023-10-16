package processor

import (
	"context"
	"github.com/ceobebot/qqchannel/infrastructure/log"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"strings"
)

func AtMessageEventHandler(api openapi.OpenAPI) event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		logger.Info(log.NewFieldsWithMessage("message received").With("content", input))
		cmd := message.ParseCommand(input)
		mustHandlers, optionalHandlers := DefaultMatcher().MatchHandlers(cmd.Cmd, cmd.Content)

		ctx := NewContext(context.Background(), api, Payload{
			Message:    data,
			Event:      event,
			Command:    cmd.Cmd,
			Content:    cmd.Content,
			RawContent: data.Content,
		}, mustHandlers, optionalHandlers)

		ctx.Next()
		return nil
	}
}
