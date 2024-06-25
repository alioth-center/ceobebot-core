package core

import (
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func initRegister() {
	RegisterLimiter("user", ctxext.LimitByUser)
	RegisterLimiter("group", ctxext.LimitByGroup)

	RegisterTriggerRule("bot_owner", zero.SuperUserPermission)
	RegisterTriggerRule("group_owner", zero.OwnerPermission)
	RegisterTriggerRule("group_admin", zero.AdminPermission)
	RegisterTriggerRule("only_to_me", zero.OnlyToMe)
	RegisterTriggerRule("only_private", zero.OnlyPrivate)
	RegisterTriggerRule("only_public", zero.OnlyPublic)
	RegisterTriggerRule("only_group", zero.OnlyGroup)
	RegisterTriggerRule("only_guild", zero.OnlyGuild)
}
