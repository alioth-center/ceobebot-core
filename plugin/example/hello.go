package example

import (
	"fmt"
	"github.com/ceobebot/qqchannel/plugin"
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
)

func init() {
	plugin.RegisterPlugin(HelloPlugin{})
}

type HelloPlugin struct{}

func (h HelloPlugin) Info() plugin.MessagePluginInfo {
	return plugin.MessagePluginInfo{
		Name:       "问候插件",
		TriggerKey: "hello",
		BlockChain: false,
	}
}

func (h HelloPlugin) TriggerKey() string {
	return "/hello"
}

func (h HelloPlugin) Commands() []plugin.MessageCommand {
	return []plugin.MessageCommand{
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
