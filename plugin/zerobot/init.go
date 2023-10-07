package zerobot

import (
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
)

func init() {
	plugin.RegisterPlugin(Plugin{})
}

type Plugin struct{}

func (p Plugin) TriggerKey() string {
	return "/zero"
}

func (p Plugin) Commands() []plugin.Command {
	return []plugin.Command{
		MenuCommand{},
	}
}
