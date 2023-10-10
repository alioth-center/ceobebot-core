package message

import (
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/config"
)

type Emoji = string

type EmojiConfig struct {
	Keyword string `yaml:"keyword"`
	Emoji   Emoji  `yaml:"emoji_id"`
}

var (
	emojiMap = map[string]Emoji{}
)

func init() {
	var emojiConfigs []EmojiConfig
	loadConfigErr := config.LoadExternalConfigWithKeys(&emojiConfigs, "data/emoji.yaml", "emoji")
	if loadConfigErr != nil {
		panic(loadConfigErr)
	}

	for _, emojiConfig := range emojiConfigs {
		emojiMap[emojiConfig.Keyword] = emojiConfig.Emoji
	}
}

func GetEmojiFormKeyword(keyword string) (emoji Emoji, exist bool) {
	emoji, exist = emojiMap[keyword]
	return
}
