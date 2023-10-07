package chat

import (
	"fmt"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
	"time"
)

type GptCommand struct {
	plugin.TextReplyMessageBaseImplementation
}

func (g GptCommand) Name() string {
	return "ChatGPT"
}

func (g GptCommand) Description() string {
	return "与ChatGPT3.5对话"
}

func (g GptCommand) Example() string {
	return "/chat ${你的信息}"
}

func (g GptCommand) Triggered(content string) (triggered bool) {
	return content != ""
}

func (g GptCommand) Handle(payload processor.Payload) (replyMessage string) {
	start := time.Now()
	reply, additional := client.ReplyConversation(payload.Content)
	end := time.Since(start)
	return fmt.Sprintf("%s\n----\n%s\n耗时:%s", reply, additional, end.String())
}
