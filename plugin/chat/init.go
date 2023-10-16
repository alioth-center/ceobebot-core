package chat

import (
	"github.com/ceobebot/qqchannel/infrastructure/config"
	"github.com/ceobebot/qqchannel/infrastructure/sqlite"
	"github.com/ceobebot/qqchannel/plugin"
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

	for _, perm := range chatConfig.DefaultPermission {
		if p, exist := format(perm); !exist {
			panic("invalid default permission: " + perm)
		} else {
			defaultPermission = defaultPermission.AddPermission(p)
		}
	}

	plugin.RegisterPlugin(Plugin{})
}

type Plugin struct{}

func (p Plugin) TriggerKey() string {
	return "/ai"
}

func (p Plugin) Commands() []plugin.MessageCommand {
	return []plugin.MessageCommand{
		GptCommand{},
		DrawCommand{},
		ManagementCommand{},
	}
}
