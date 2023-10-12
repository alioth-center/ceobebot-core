package chat

import (
	"github.com/ceobebot/qqchannel/plugin"
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
	"strconv"
	"strings"
)

type DrawCommand struct {
	plugin.ImageReplyCommand
}

func (d DrawCommand) Name() string { return "Dalle" }

func (d DrawCommand) Description() string { return "使用Dalle画画" }

func (d DrawCommand) Example() string {
	return "/ai draw [(可选)画画大小，小/中/大] [画画描述]"
}

func (d DrawCommand) Triggered(content string) (triggered bool) {
	return strings.HasPrefix(content, "draw ")
}

func (d DrawCommand) Handle(payload processor.Payload) (replyMessage message.Message) {
	var size ImageSize
	args := strings.Split(payload.Content, " ")
	var prompt = ""
	if len(args) < 2 {
		return message.NewTextMessage().Text("请输入画画描述").Reference(payload.Message.ID)
	} else if len(args) >= 3 {
		switch args[1] {
		case "小", "small", "s", "1":
			size = SizeSmall
			prompt = strings.Join(args[2:], " ")
		case "中", "medium", "m", "2":
			size = SizeMedium
			prompt = strings.Join(args[2:], " ")
		case "大", "large", "l", "3":
			size = SizeLarge
			prompt = strings.Join(args[2:], " ")
		default:
			size = SizeSmall
			prompt = strings.Join(args[1:], " ")
		}
	} else {
		size = SizeSmall
		prompt = args[1]
	}

	userID, getUserErr := strconv.ParseUint(payload.Message.Author.ID, 10, 64)
	if getUserErr != nil {
		return message.NewTextMessage().Text("获取用户ID失败").Reference(payload.Message.ID)
	}

	var permission string
	switch size {
	case SizeSmall:
		permission = PermissionImageSmall
	case SizeMedium:
		permission = PermissionImageMedium
	case SizeLarge:
		permission = PermissionImageLarge
	default:
		return message.NewTextMessage().Text("未知的使用权限").Reference(payload.Message.ID)
	}

	if hasPermission, checkErr := checkPermission(userID, permission); checkErr != nil {
		return message.NewTextMessage().Text("检查权限失败，请联系管理员查看：" + checkErr.Error()).Reference(payload.Message.ID)
	} else if !hasPermission {
		return message.NewTextMessage().Text("您没有使用中画画的权限").Reference(payload.Message.ID)
	}

	url, drawErr := client.DrawPicture(prompt, size)
	if drawErr != nil {
		return message.NewTextMessage().Text("画画失败，请联系管理员查看：" + drawErr.Error()).Reference(payload.Message.ID)
	}

	return message.NewImageMessage().Image(url).Reference(payload.Message.ID)
}
