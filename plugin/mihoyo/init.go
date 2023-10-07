package mihoyo

import (
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
)

func init() {
	plugin.RegisterPlugin(Plugin{})
}

type Plugin struct{}

func (p Plugin) TriggerKey() string {
	return "/mihoyo"
}

func (p Plugin) Commands() []plugin.Command {
	return []plugin.Command{
		GenshinCommand{},
		StarRailCommand{},
	}
}
