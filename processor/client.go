package processor

import (
	"github.com/tencent-connect/botgo/openapi"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/log"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor/message"
)

type Client interface {
	ReplyMessage(ctx Context, reply message.Message)
	SendMessage(ctx Context, content message.Message)
}

type client struct {
	api openapi.OpenAPI
}

func (c *client) ReplyMessage(ctx Context, reply message.Message) {
	payload := ctx.GetPayload()
	msg := reply.Build()
	msg.MsgID = payload.Message.ID
	result, err := c.api.PostMessage(ctx.GetContext(), payload.Message.ChannelID, msg)
	if err != nil {
		result, _ = c.api.PostMessage(ctx.GetContext(), payload.Message.ChannelID, message.NewTextMessage().Text("回复被夹掉了，换个话题吧").Build())
	}
	logger.Info(log.NewFieldsWithMessage("message sent").With("content", result.Content))
}

func (c *client) SendMessage(ctx Context, content message.Message) {
	result, err := c.api.PostMessage(ctx.GetContext(), ctx.GetPayload().Message.ChannelID, content.Build())
	if err != nil {
		result, _ = c.api.PostMessage(ctx.GetContext(), ctx.GetPayload().Message.ChannelID, message.NewTextMessage().Text("回复被夹掉了，换个话题吧").Build())
	}
	logger.Info(log.NewFieldsWithMessage("message sent").With("content", result.Content))
}

func NewClient(api openapi.OpenAPI) Client {
	return &client{api: api}
}
