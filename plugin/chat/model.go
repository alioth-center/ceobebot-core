package chat

type GptModel string

const (
	Gpt3Dot5Turbo        GptModel = "gpt-3.5-turbo"
	Gpt3Dot5Turbo16K     GptModel = "gpt-3.5-turbo-16k"
	Gpt3Dot5Turbo0613    GptModel = "gpt-3.5-turbo-0613"
	Gpt3Dot5Turbo16K0613 GptModel = "gpt-3.5-turbo-16k-0613"
	Gpt4                 GptModel = "gpt-4"
	Gpt40314             GptModel = "gpt-4-0314"
	Gpt40613             GptModel = "gpt-4-0613"
	Gpt4Poe              GptModel = "gpt-4-poe"
)

var (
	AllSupportedModels = []GptModel{
		Gpt3Dot5Turbo, Gpt3Dot5Turbo16K, Gpt3Dot5Turbo0613, Gpt3Dot5Turbo16K0613, Gpt4, Gpt40314, Gpt40613, Gpt4Poe,
	}

	SupportedGptModels = map[string]GptModel{
		"gpt-3.5-turbo":          Gpt3Dot5Turbo,
		"gpt-3.5-turbo-16k":      Gpt3Dot5Turbo16K,
		"gpt-3.5-turbo-0613":     Gpt3Dot5Turbo0613,
		"gpt-3.5-turbo-16k-0613": Gpt3Dot5Turbo16K0613,
		"gpt-4":                  Gpt4,
		"gpt-4-poe":              Gpt4Poe,
		"gpt-4-0314":             Gpt40314,
		"gpt-4-0613":             Gpt40613,
	}
)

type GptMessageRole string

const (
	RoleUser      GptMessageRole = "user"
	RoleSystem    GptMessageRole = "system"
	RoleAssistant GptMessageRole = "assistant"
)

var (
	SupportedGptMessageRoles = map[string]GptMessageRole{
		"user":      RoleUser,
		"system":    RoleSystem,
		"assistant": RoleAssistant,
	}
)

type GptMessage struct {
	Role    GptMessageRole `json:"role"`
	Content string         `json:"content"`
}

type GptMessageChain struct {
	messages []GptMessage
}

func (chain *GptMessageChain) AddMessage(message GptMessage) *GptMessageChain {
	chain.messages = append(chain.messages, message)
	return chain
}

func (chain *GptMessageChain) AddMessages(messages ...GptMessage) *GptMessageChain {
	chain.messages = append(chain.messages, messages...)
	return chain
}

func (chain *GptMessageChain) MergeChainAfter(another *GptMessageChain) *GptMessageChain {
	chain.messages = append(chain.messages, another.messages...)
	return chain
}

func (chain *GptMessageChain) MergeChainBefore(another *GptMessageChain) *GptMessageChain {
	chain.messages = append(another.messages, chain.messages...)
	return chain
}

func (chain *GptMessageChain) GetMessages() []GptMessage {
	return chain.messages
}

func NewGptMessageChain(messages ...GptMessage) *GptMessageChain {
	return &GptMessageChain{
		messages: messages,
	}
}

type GptCompletionsRequest struct {
	Model            GptModel     `json:"model"`
	Messages         []GptMessage `json:"messages"`
	Temperature      float64      `json:"temperature"`
	PresencePenalty  float64      `json:"presence_penalty"`
	FrequencyPenalty float64      `json:"frequency_penalty"`
	TopP             float64      `json:"top_p"`
	Stream           bool         `json:"stream"`
}

type GptConfigOptions struct {
	Model            GptModel `json:"model"`
	Temperature      float64  `json:"temperature"`
	PresencePenalty  float64  `json:"presence_penalty"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
}

func NewGptCompletionsRequest(options GptConfigOptions, messages *GptMessageChain) GptCompletionsRequest {
	request := GptCompletionsRequest{
		Model:            options.Model,
		Messages:         messages.GetMessages(),
		Temperature:      options.Temperature,
		PresencePenalty:  options.PresencePenalty,
		FrequencyPenalty: options.FrequencyPenalty,
		TopP:             1,
		Stream:           false,
	}

	if _, exist := SupportedGptModels[string(request.Model)]; !exist {
		request.Model = Gpt3Dot5Turbo
	}

	if request.Temperature > 1.0 || request.Temperature < 0.0 {
		request.Temperature = 0.5
	}

	if request.PresencePenalty > 1.0 || request.PresencePenalty < 0.0 {
		request.PresencePenalty = 0.0
	}

	if request.FrequencyPenalty > 1.0 || request.FrequencyPenalty < 0.0 {
		request.FrequencyPenalty = 0.0
	}

	return request
}

func NewGptMessage(message string, role GptMessageRole) *GptMessageChain {
	gptMsg := GptMessage{
		Role:    role,
		Content: message,
	}

	if _, exist := SupportedGptMessageRoles[string(gptMsg.Role)]; !exist {
		gptMsg.Role = RoleUser
	}

	return NewGptMessageChain(gptMsg)
}

func NewGptPrompt(prompt string) *GptMessageChain {
	return NewGptMessageChain(GptMessage{
		Role:    RoleSystem,
		Content: prompt,
	})
}

func NewGptMessageWithPrompt(prompt string, message string, role GptMessageRole) *GptMessageChain {
	return NewGptPrompt(prompt).MergeChainAfter(NewGptMessage(message, role))
}

type GptReplyChoices struct {
	Index        int        `json:"index"`
	Message      GptMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type GptUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GptCompletionsResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []GptReplyChoices `json:"choices"`
	Usage   GptUsage          `json:"usage"`
}

func GetModelThousandTokenPrice(_ GptModel) (prompt float64, answer float64) {
	return 0.03, 0.04
}
