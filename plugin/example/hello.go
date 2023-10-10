package example

import (
	"fmt"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor/message"
)

func init() {
	plugin.RegisterPlugin(HelloPlugin{})
}

type HelloPlugin struct{}

func (h HelloPlugin) TriggerKey() string {
	return "/hello"
}

func (h HelloPlugin) Commands() []plugin.Command {
	return []plugin.Command{
		HelloCommand{},
	}
}

type HelloCommand struct {
	plugin.TextReplyMessageBaseImplementation
}

func (h HelloCommand) Name() string {
	return "Hello"
}

func (h HelloCommand) Description() string {
	return "你好，世界！"
}

func (h HelloCommand) Example() string {
	return "/hello"
}

func (h HelloCommand) Triggered(_ string) (triggered bool) {
	return true
}

func (h HelloCommand) Handle(payload processor.Payload) (replyMessage message.Message) {
	return message.NewTextMessage().
		At(payload.Message.Author.ID).
		Text(fmt.Sprintf("Hello, %s!", payload.Message.Author.Username))
}
