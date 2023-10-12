package zerobot

import (
	"fmt"
	"github.com/ceobebot/qqchannel/infrastructure/config"
	"github.com/ceobebot/qqchannel/plugin"
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
	"math/rand"
	"time"
)

var (
	menu menuEnum
)

type menuEnum struct {
	Breakfast []string `yaml:"breakfast"`
	Lunch     []string `yaml:"lunch"`
	Dinner    []string `yaml:"dinner"`
}

func init() {
	if loadConfigErr := config.LoadExternalConfig(&menu, "data/zero/menu.yaml"); loadConfigErr != nil {
		panic(loadConfigErr)
	}
}

func getBreakfast() string {
	return menu.Breakfast[rand.Intn(len(menu.Breakfast))]
}

func getLunch() string {
	return menu.Lunch[rand.Intn(len(menu.Lunch))]
}

func getDinner() string {
	return menu.Dinner[rand.Intn(len(menu.Dinner))]
}

type MenuCommand struct {
	plugin.TextReplyMessageBaseImplementation
}

func (m MenuCommand) Name() string {
	return "menu"
}

func (m MenuCommand) Description() string {
	return "食谱菜单"
}

func (m MenuCommand) Example() string {
	return "/zero 吃什么\n/zero ${早上/中午/晚上}吃什么"
}

func (m MenuCommand) Triggered(content string) (triggered bool) {
	if content == "吃什么" || content == "早上吃什么" || content == "中午吃什么" || content == "晚上吃什么" {
		return true
	}

	return false
}

func (m MenuCommand) Handle(payload processor.Payload) (reply message.Message) {
	replyMessage := ""
	getMenuTemplate := "%s今天%s的食谱为：%s"
	switch payload.Content {
	case "吃什么":
		switch time.Now().Hour() {
		case 6, 7, 8, 9:
			replyMessage = fmt.Sprintf(getMenuTemplate, payload.Message.Author.Username, "早餐", getBreakfast())
		case 10, 11, 12, 13, 14:
			replyMessage = fmt.Sprintf(getMenuTemplate, payload.Message.Author.Username, "午餐", getLunch())
		case 15, 16, 17, 18, 19, 20, 21, 22, 23:
			replyMessage = fmt.Sprintf(getMenuTemplate, payload.Message.Author.Username, "晚餐", getDinner())
		default:
			replyMessage = "饭店这个时候都关了，吃西北风吧"
		}
	case "早上吃什么":
		replyMessage = fmt.Sprintf(getMenuTemplate, payload.Message.Author.Username, "早餐", getBreakfast())
	case "中午吃什么":
		replyMessage = fmt.Sprintf(getMenuTemplate, payload.Message.Author.Username, "午餐", getLunch())
	case "晚上吃什么":
		replyMessage = fmt.Sprintf(getMenuTemplate, payload.Message.Author.Username, "晚餐", getDinner())
	default:
		replyMessage = "饭店这个时候都关了，吃西北风吧"
	}

	return message.NewTextMessage().At(payload.Message.Author.ID).Text(replyMessage)
}
