package chat

import (
	"encoding/json"
	"fmt"
	"github.com/ceobebot/qqchannel/infrastructure/log"
	"github.com/parnurzeal/gorequest"
	"strings"
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

func (c Client) ReplyConversation(text string, modelOpt GptModel) (reply string, footer string) {
	model, gptModel := "gpt-3.5-turbo", Gpt3Dot5Turbo
	if modelOpt != Gpt3Dot5Turbo {
		model = string(modelOpt)
		gptModel = modelOpt
	}

	payload := NewGptCompletionsRequest(
		GptConfigOptions{
			Model:            gptModel,
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

	if len(response.Choices) == 0 {
		reply, footer = fmt.Sprintf("ERROR: 没有回复\n%+v", openAiResponse), ""
	} else {
		promptPrice, answerPrice := GetModelThousandTokenPrice(Gpt3Dot5Turbo)
		promptCost := float64(response.Usage.PromptTokens) / float64(1000) * promptPrice
		answerCost := float64(response.Usage.CompletionTokens) / float64(1000) * answerPrice

		if gptModel == Gpt3Dot5Turbo {
			footer = fmt.Sprintf(
				"回复自%s模型\ntoken消耗：%d，扣费: %.6f CNY",
				model, response.Usage.TotalTokens,
				promptCost+answerCost,
			)
		} else {
			footer = fmt.Sprintf(
				"回复自%s模型\ntoken消耗：%d，扣费: N/A(只支持GPT3.5计算费用)",
				model, response.Usage.TotalTokens,
			)
		}

		reply = response.Choices[0].Message.Content
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
