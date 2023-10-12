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
		Level:     "debug",
		Formatter: "json",
		FilePath:  "data/chat/history.log",
	})
	routerMap = map[ApiType]string{
		GptApi:   "chat/completions",
		ImageApi: "images/generations",
		ShortUrl: "https://tinyurl.com/api-create.php",
	}
)

func getApiUrl(api ApiType) string {
	return fmt.Sprintf("%s/%s/%s", chatConfig.BaseUrl, chatConfig.ApiVersion, routerMap[api])
}

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
		request.Post(getApiUrl(GptApi)).
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

func (c Client) DrawPicture(prompt string, size ImageSize) (u string, err error) {
	payload := ImageGenerationRequest{
		Prompt: prompt,
		Number: 1,
		Size:   size.toRequestSize(),
	}

	bytes, marshalErr := json.Marshal(&payload)
	if marshalErr != nil {
		return "", marshalErr
	}

	request := gorequest.New()
	responses, body, responseErrs :=
		request.Post(getApiUrl(ImageApi)).
			Set("Authorization", fmt.Sprintf("Bearer %s", chatConfig.AppToken)).
			Send(string(bytes)).
			End()

	if len(responseErrs) > 0 {
		var errStrings []string
		for _, errOne := range responseErrs {
			errStrings = append(errStrings, errOne.Error())
		}

		return "", fmt.Errorf(strings.Join(errStrings, " -> \n"))
	}

	var response ImageGenerationResponse
	unmarshalErr := json.Unmarshal([]byte(body), &response)
	if unmarshalErr != nil {
		return "", unmarshalErr
	}

	if len(response.Data) == 0 {
		return "", fmt.Errorf("no data: %v", responses)
	} else if len(response.Data) > 1 {
		return "", fmt.Errorf("too many data: %v", responses)
	}

	logger.Info(log.NewFieldsWithMessage("draw picture").With("prompt", prompt).With("url", response.Data[0].Url))
	return response.Data[0].Url, nil
}
