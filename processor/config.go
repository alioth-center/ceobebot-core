package processor

import "github.com/ceobebot/qqchannel/infrastructure/config"

var (
	systemConfig = SystemConfig{}
)

func init() {
	if loadConfigErr := config.LoadCustomConfig(&systemConfig, "bot_config"); loadConfigErr != nil {
		panic(loadConfigErr)
	}
}

type SystemConfig struct {
	AppID          uint64 `yaml:"app_id"`
	AppSecret      string `yaml:"app_secret"`
	AppToken       string `yaml:"app_token"`
	TimeoutSecond  int    `yaml:"timeout_second"`
	TestMode       bool   `yaml:"test_mode"`
	UndefinedReply string `yaml:"undefined_reply"`
}
