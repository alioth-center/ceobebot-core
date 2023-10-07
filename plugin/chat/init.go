package chat

import (
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
)

func init() {
	plugin.RegisterPlugin(Plugin{})
}

type Plugin struct{}

func (p Plugin) TriggerKey() string {
	return "/chat"
}

func (p Plugin) Commands() []plugin.Command {
	return []plugin.Command{
		GptCommand{},
	}
}
