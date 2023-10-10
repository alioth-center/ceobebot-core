package processor

import (
	"github.com/tencent-connect/botgo/openapi"
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
	_, _ = c.api.PostMessage(ctx.GetContext(), payload.Message.ChannelID, msg)
}

func (c *client) SendMessage(ctx Context, content message.Message) {
	_, _ = c.api.PostMessage(ctx.GetContext(), ctx.GetPayload().Message.ChannelID, content.Build())
}

func NewClient(api openapi.OpenAPI) Client {
	return &client{api: api}
}
