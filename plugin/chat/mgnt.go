package chat

import (
	"fmt"
	"github.com/ceobebot/qqchannel/plugin"
	"github.com/ceobebot/qqchannel/processor"
	"github.com/ceobebot/qqchannel/processor/message"
	"strings"
)

type ManagementCommand struct {
	plugin.TextReplyMessageBaseImplementation
}

func (m ManagementCommand) Name() string {
	return "ChatManagement"
}

func (m ManagementCommand) Description() string {
	return "管理AI插件的权限"
}

func (m ManagementCommand) Example() string {
	return "/ai mgnt [子命令]"
}

func (m ManagementCommand) Triggered(content string) (triggered bool) {
	return strings.HasPrefix(content, "mgnt ")
}

func (m ManagementCommand) Handle(payload processor.Payload) (replyMessage message.Message) {
	command := strings.Split(payload.Content, " ")[1]
	switch command {
	case "ap", "添加权限", "add", "添加":
		return m.addPermission(payload)
	case "rp", "删除权限", "remove", "删除":
		return m.removePermission(payload)
	case "lp", "列举权限", "list", "列举":
		return m.listPermission(payload)
	case "dp", "默认权限", "default", "默认":
		permissions := strings.Join(defaultPermission.AllPermissions(), ", ")
		return message.NewTextMessage().Text(fmt.Sprintf("默认权限有: %s", permissions))
	default:
		return message.NewTextMessage().Text("未知的子命令")
	}
}

func (m ManagementCommand) addPermission(payload processor.Payload) (replyMessage message.Message) {
	userList := message.GetAtMembersFromRawContent(payload.RawContent, 1)
	if len(userList) == 0 {
		return message.NewTextMessage().Text("没有检查到足够的参数，请 @ 需要添加权限的用户")
	}

	if addErr := addPermission(payload.Message.Author.ID, payload.Message.Author.Username, strings.Split(payload.Content, " ")[2:]...); addErr != nil {
		return message.NewTextMessage().Text(fmt.Sprintf("添加权限失败: %s", addErr.Error()))
	} else {
		return message.NewTextMessage().Text(fmt.Sprintf("添加权限成功: 共 %d 条记录", len(userList)))
	}
}

func (m ManagementCommand) removePermission(payload processor.Payload) (replyMessage message.Message) {
	userList := message.GetAtMembersFromRawContent(payload.RawContent, 1)
	if len(userList) == 0 {
		return message.NewTextMessage().Text("没有检查到足够的参数，请 @ 需要删除权限的用户")
	} else if len(userList) > 1 {
		return message.NewTextMessage().Text("只能删除一个用户的权限")
	}

	if removeErr := removePermission(payload.Message.Author.ID, strings.Split(payload.Content, " ")[2:]...); removeErr != nil {
		return message.NewTextMessage().Text(fmt.Sprintf("删除权限失败: %s", removeErr.Error()))
	} else {
		return message.NewTextMessage().Text(fmt.Sprintf("删除权限成功: 共 %d 条记录", len(userList)))
	}
}

func (m ManagementCommand) listPermission(payload processor.Payload) (replyMessage message.Message) {
	userList := message.GetAtMembersFromRawContent(payload.RawContent, 1)
	if len(userList) == 0 {
		return message.NewTextMessage().Text("没有检查到足够的参数，请 @ 需要列举权限的用户")
	} else if len(userList) > 1 {
		return message.NewTextMessage().Text("只能列举一个用户的权限")
	}

	if permission, listErr := getPermission(payload.Message.Author.ID); listErr != nil {
		return message.NewTextMessage().Text(fmt.Sprintf("列举权限失败: %s", listErr.Error()))
	} else {
		permissions := strings.Join(PermissionData(permission.Permissions).AllPermissions(), ", ")
		return message.NewTextMessage().Text(fmt.Sprintf("用户 %s 的权限有: %s", permission.UserName, permissions))
	}
}
