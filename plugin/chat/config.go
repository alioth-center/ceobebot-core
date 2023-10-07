package chat

import "studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/config"

var (
	chatConfig Config
)

func init() {
	if loadConfigErr := config.LoadExternalConfigWithKeys(&chatConfig, "data/chat/config.yaml", "chat"); loadConfigErr != nil {
		panic(loadConfigErr)
	}
}

type Config struct {
	BaseUrl          string  `yaml:"base_url"`
	ApiVersion       string  `yaml:"api_version"`
	AppToken         string  `yaml:"app_token"`
	Prompt           string  `yaml:"prompt"`
	Temperature      float64 `yaml:"temperature"`
	PresencePenalty  float64 `yaml:"presence_penalty"`
	FrequencyPenalty float64 `yaml:"frequency_penalty"`
}
