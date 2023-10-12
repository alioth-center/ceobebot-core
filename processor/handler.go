package processor

import (
	"github.com/ceobebot/qqchannel/processor/message"
	txMessage "github.com/tencent-connect/botgo/dto/message"
	"strings"
)

func init() {
	// 设置默认未找到命令回复
	var notFoundHandler OptionalHandleFunction
	switch strings.ToLower(systemConfig.UndefinedReply) {
	case "not_matched_if_no_help":
		notFoundHandler = NotFoundCommandAfterIncorrectCommandCheckReplyHandler()
	case "ignore":
		notFoundHandler = nil
	case "not_matched":
		notFoundHandler = NotMatchedCommandReplyHandler()
	case "help":
		notFoundHandler = IncorrectCommandReplyHandler()
	default:
		notFoundHandler = nil
	}
	defaultMatcher.SetDefaultHandlerFunctions([]MustHandleFunction{}, []OptionalHandleFunction{notFoundHandler})

	// 注册默认指令
	defaultMatcher.RegisterCommand("/list", "List Commands", "列出某服务的所有指令",
		"/list /list", func(s string) bool { return strings.HasPrefix(s, "/") }, []MustHandleFunction{},
		[]OptionalHandleFunction{ListCommandReplyHandler()})
	defaultMatcher.RegisterCommand("/list", "List All Commands", " 列出所有服务",
		"/list", func(s string) bool { return s == "" }, []MustHandleFunction{},
		[]OptionalHandleFunction{ListAllCommandReplyHandler()})
	defaultMatcher.RegisterCommand("/help", "Help", "查看某服务的指令帮助",
		"/help /help", func(s string) bool { return strings.HasPrefix(s, "/") }, []MustHandleFunction{},
		[]OptionalHandleFunction{HelpCommandReplyHandler()})
}

func NotMatchedCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		ctx.Next()
		payload := ctx.GetPayload()
		msg := message.NewTextMessage().
			Reference(payload.Message.ID).
			Text(FormatNotMatchedCommandMessage(payload.Command, payload.Content))
		ctx.GetApi().ReplyMessage(ctx, msg)
	}
}

func IncorrectCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		ctx.Next()
		payload := ctx.GetPayload()
		commandInfos := defaultMatcher.GetHelpInfo(payload.Command, payload.Content)
		msg := message.NewTextMessage().
			Reference(payload.Message.ID).
			Text(FormatIncorrectCommandArgumentMessage(commandInfos))
		ctx.GetApi().ReplyMessage(ctx, msg)
	}
}

func ListCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		payload := ctx.GetPayload()
		var infos []HandlerInfo
		for _, info := range defaultMatcher.ListCommands(payload.Content) {
			infos = append(infos, info)
		}
		msg := message.NewTextMessage().
			Reference(payload.Message.ID).
			Text(FormatListHelpInfoMessage(payload.Content, infos))
		ctx.GetApi().ReplyMessage(ctx, msg)
	}
}

func ListAllCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		payload := ctx.GetPayload()
		msg := message.NewTextMessage().
			Reference(payload.Message.ID).
			Text(FormatListAllServiceInfoMessage())
		ctx.GetApi().ReplyMessage(ctx, msg)
	}
}

func NotFoundCommandAfterIncorrectCommandCheckReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		ctx.Next()
		payload := ctx.GetPayload()

		if ctx.IsAborted() {
			commandInfos := defaultMatcher.GetHelpInfo(payload.Command, payload.Content)
			// 如果被中断且找到了命令，说明是命令参数错误，需要提示
			if commandInfos != nil && len(commandInfos) != 0 {
				msg := message.NewTextMessage().
					Reference(payload.Message.ID).
					Text(FormatIncorrectCommandArgumentMessage(commandInfos))
				ctx.GetApi().ReplyMessage(ctx, msg)
				return
			}
		}
		msg := message.NewTextMessage().
			Reference(payload.Message.ID).
			Text(FormatNotMatchedCommandMessage(payload.Command, payload.Content))
		ctx.GetApi().ReplyMessage(ctx, msg)
	}
}

func HelpCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		payload := ctx.GetPayload()
		command := txMessage.ParseCommand(payload.Content)
		commandInfos := defaultMatcher.GetHelpInfo(command.Cmd, command.Content)
		msg := message.NewTextMessage().
			Reference(payload.Message.ID).
			Text(FormatCommandHelpArgumentMessage(command.Cmd, command.Content, commandInfos))
		ctx.GetApi().ReplyMessage(ctx, msg)
	}
}
