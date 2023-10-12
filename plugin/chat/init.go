package chat

import (
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/config"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/sqlite"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin"
)

var (
	database sqlite.Database
)

func init() {
	db, initDbErr := sqlite.NewSqliteDb("chat/database.db", &Permission{})
	if initDbErr != nil {
		panic(initDbErr)
	} else {
		database = db
	}

	if loadConfigErr := config.LoadExternalConfigWithKeys(&chatConfig, "data/chat/config.yaml", "chat"); loadConfigErr != nil {
		panic(loadConfigErr)
	}

	masterConfig := map[uint64][]uint64{}
	if loadConfigErr := config.LoadExternalConfigWithKeys(&masterConfig, "data/chat/config.yaml", "master", "chanel"); loadConfigErr != nil {
		panic(loadConfigErr)
	} else {
		loadMasters(masterConfig)
	}

	plugin.RegisterPlugin(Plugin{})
}

type Plugin struct{}

func (p Plugin) TriggerKey() string {
	return "/ai"
}

func (p Plugin) Commands() []plugin.Command {
	return []plugin.Command{
		GptCommand{},
		DrawCommand{},
	}
}
