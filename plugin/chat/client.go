package chat

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"strings"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/log"
	"time"
)

var (
	client = Client{}
	logger = log.NewLogger(log.Config{
		Level:     "info",
		Formatter: "json",
		FilePath:  "data/chat/history.log",
	})
)

type Client struct{}

func (c Client) ReplyConversation(text string) (reply string, footer string) {
	model := "gpt-3.5-turbo"
	payload := NewGptCompletionsRequest(
		GptConfigOptions{
			Model:            Gpt3Dot5Turbo,
			Temperature:      chatConfig.Temperature,
			PresencePenalty:  chatConfig.PresencePenalty,
			FrequencyPenalty: chatConfig.FrequencyPenalty,
		},
		NewGptMessageWithPrompt(chatConfig.Prompt, text, RoleUser),
	)

	bytes, err := json.Marshal(&payload)
	if err != nil {
		return err.Error(), ""
	}

	request := gorequest.New()
	openAiResponse, body, responseErrs :=
		request.Post(fmt.Sprintf("%s/%s/chat/completions", chatConfig.BaseUrl, chatConfig.ApiVersion)).
			Set("Authorization", fmt.Sprintf("Bearer %s", chatConfig.AppToken)).
			Send(string(bytes)).
			End()

	if len(responseErrs) > 0 {
		var errStrings []string
		for _, errOne := range responseErrs {
			errStrings = append(errStrings, errOne.Error())
		}

		return strings.Join(errStrings, " -> \n"), ""
	}

	var response GptCompletionsResponse
	unmarshalErr := json.Unmarshal([]byte(body), &response)
	if unmarshalErr != nil {
		return fmt.Sprintf("unmarshal: %s", unmarshalErr.Error()), ""
	}

	timeUnix := time.Unix(response.Created, 0)
	timeUnix = timeUnix.In(time.Local)
	timeString := timeUnix.Format("2006年1月2日15:04")

	fmt.Println("problem", text, "cost", response.Usage.TotalTokens, "tokens")

	if len(response.Choices) == 0 {
		reply, footer = fmt.Sprintf("ERROR: 没有回复\n%+v", openAiResponse), fmt.Sprintf(
			"在%s回复自%s模型\n问题token消耗：%d，回复token消耗：%d，总token消耗：%d\n扣费: %d CNY",
			time.Now().In(time.Local).Format("2006年1月2日15:04"), model, 0, 0, 0, 0,
		)
	} else {
		promptPrice, answerPrice := GetModelThousandTokenPrice(Gpt3Dot5Turbo)
		promptCost := float64(response.Usage.PromptTokens) / float64(1000) * promptPrice
		answerCost := float64(response.Usage.CompletionTokens) / float64(1000) * answerPrice
		reply, footer = response.Choices[0].Message.Content, fmt.Sprintf(
			"在%s回复自%s模型\n问题token消耗：%d，回复token消耗：%d，总token消耗：%d\n扣费: %.6f CNY",
			timeString, model, response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens,
			promptCost+answerCost,
		)
	}

	logger.Info(log.NewFieldsWithMessage("complete conversation").With("question", text).With("answer", reply).With("options", footer))
	return reply, footer
}
