package plugin

import (
	"github.com/ceobebot/qqchannel/processor"
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
		ctx.GetApi().ReplyMessage(ctx, command.Handle(ctx.GetPayload()))
	}
}
