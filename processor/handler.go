package processor

import (
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"strings"
)

func init() {
	defaultMatcher.SetDefaultHandlerFunctions([]MustHandleFunction{}, []OptionalHandleFunction{
		NotFoundCommandAfterIncorrectCommandCheckReplyHandler(),
	})
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
		_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), payload.Message.ChannelID, &dto.MessageToCreate{
			Content: FormatNotMatchedCommandMessage(payload.Command, payload.Content),
			MessageReference: &dto.MessageReference{
				MessageID:             payload.Message.ID,
				IgnoreGetMessageError: true,
			},
		})
	}
}

func IncorrectCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		ctx.Next()
		payload := ctx.GetPayload()
		commandInfos := defaultMatcher.GetHelpInfo(payload.Command, payload.Content)
		_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), payload.Message.ChannelID, &dto.MessageToCreate{
			Content: FormatIncorrectCommandArgumentMessage(commandInfos),
			MessageReference: &dto.MessageReference{
				MessageID:             payload.Message.ID,
				IgnoreGetMessageError: true,
			},
		})
	}
}

func ListCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		payload := ctx.GetPayload()
		var infos []HandlerInfo
		for _, info := range defaultMatcher.ListCommands(payload.Content) {
			infos = append(infos, info)
		}

		_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), payload.Message.ChannelID, &dto.MessageToCreate{
			Content: FormatListHelpInfoMessage(payload.Content, infos),
			MessageReference: &dto.MessageReference{
				MessageID:             payload.Message.ID,
				IgnoreGetMessageError: true,
			},
		})
	}
}

func ListAllCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		payload := ctx.GetPayload()
		_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), payload.Message.ChannelID, &dto.MessageToCreate{
			Content: FormatListAllServiceInfoMessage(),
			MessageReference: &dto.MessageReference{
				MessageID:             payload.Message.ID,
				IgnoreGetMessageError: true,
			},
		})
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
				_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), payload.Message.ChannelID, &dto.MessageToCreate{
					Content: FormatIncorrectCommandArgumentMessage(commandInfos),
					MessageReference: &dto.MessageReference{
						MessageID:             payload.Message.ID,
						IgnoreGetMessageError: true,
					},
				})
				return
			}
		}

		_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), payload.Message.ChannelID, &dto.MessageToCreate{
			Content: FormatNotMatchedCommandMessage(payload.Command, payload.Content),
			MessageReference: &dto.MessageReference{
				MessageID:             payload.Message.ID,
				IgnoreGetMessageError: true,
			},
		})

	}
}

func HelpCommandReplyHandler() OptionalHandleFunction {
	return func(ctx Context) {
		payload := ctx.GetPayload()
		command := message.ParseCommand(payload.Content)
		commandInfos := defaultMatcher.GetHelpInfo(command.Cmd, command.Content)
		_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), payload.Message.ChannelID, &dto.MessageToCreate{
			Content: FormatCommandHelpArgumentMessage(command.Cmd, command.Content, commandInfos),
			MessageReference: &dto.MessageReference{
				MessageID:             payload.Message.ID,
				IgnoreGetMessageError: true,
			},
		})
	}
}
