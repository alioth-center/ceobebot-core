package chat

import (
	"fmt"
	"strings"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor/message"
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
	return "/ai chat [对话内容]"
}

func (g GptCommand) Triggered(content string) (triggered bool) {
	return strings.HasPrefix(content, "chat ")
}

func (g GptCommand) Handle(payload processor.Payload) (replyMessage message.Message) {
	start := time.Now()
	reply, additional := client.ReplyConversation(strings.TrimPrefix(payload.Content, "chat "))
	end := time.Since(start)
	return message.NewTextMessage().
		At(payload.Message.Author.ID).NewLine().
		Text(reply).NewLine().
		Text(additional).NewLine().
		Text(fmt.Sprintf("耗时: %s", end.String())).
		Reference(payload.Message.ID)
}
