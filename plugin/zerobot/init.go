package zerobot

import (
	"github.com/ceobebot/qqchannel/plugin"
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
