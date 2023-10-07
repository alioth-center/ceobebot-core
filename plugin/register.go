package plugin

import (
	"github.com/tencent-connect/botgo/dto"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
)

func RegisterPlugin(plugins ...MessagePlugin) {
	for _, plugin := range plugins {
		for _, command := range plugin.Commands() {
			processor.DefaultMatcher().RegisterCommand(plugin.TriggerKey(), command.Name(), command.Description(), command.Example(),
				command.Triggered, []processor.MustHandleFunction{}, []processor.OptionalHandleFunction{ProcessMessageHandler(command)})
		}
	}
}

func ProcessMessageHandler(command TextReplyCommand) processor.OptionalHandleFunction {
	return func(ctx processor.Context) {
		_, _ = ctx.GetApi().PostMessage(ctx.GetContext(), ctx.GetPayload().Message.ChannelID, &dto.MessageToCreate{
			Content: command.Handle(ctx.GetPayload()),
			MessageReference: &dto.MessageReference{
				MessageID:             ctx.GetPayload().Message.ID,
				IgnoreGetMessageError: true,
			},
		})
	}
}
