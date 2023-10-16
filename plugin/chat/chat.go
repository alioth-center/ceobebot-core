package chat

import (
	"fmt"
	"github.com/ceobebot/qqchannel/plugin"
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
	"strings"
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
	commands := strings.Split(payload.Content, " ")

	gptModel, hasModelOpt := Gpt3Dot5Turbo, false
	if model, exist := SupportedGptModels[commands[1]]; exist {
		gptModel = model
		hasModelOpt = true
	} else {
		gptModel = Gpt3Dot5Turbo
		hasModelOpt = false
	}

	question := ""
	if hasModelOpt {
		question = strings.Join(commands[2:], " ")
	} else {
		question = strings.Join(commands[1:], " ")
	}

	reply, additional := client.ReplyConversation(question, gptModel)
	end := time.Since(start)
	return message.NewTextMessage().
		At(payload.Message.Author.ID).NewLine().
		Text(reply).NewLine().
		Text(additional).NewLine().
		Text(fmt.Sprintf("耗时: %s", end.String())).
		Reference(payload.Message.ID)
}
