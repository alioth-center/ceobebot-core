package processor

import "studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/config"

var (
	systemConfig = SystemConfig{}
)

func init() {
	if loadConfigErr := config.LoadCustomConfig(&systemConfig, "bot_config"); loadConfigErr != nil {
		panic(loadConfigErr)
	}
}

type SystemConfig struct {
	AppID         uint64 `yaml:"app_id"`
	AppSecret     string `yaml:"app_secret"`
	AppToken      string `yaml:"app_token"`
	TimeoutSecond int    `yaml:"timeout_second"`
	TestMode      bool   `yaml:"test_mode"`
}
