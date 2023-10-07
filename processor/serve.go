package processor

import (
	"context"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"time"
)

func Serve() {
	botToken := token.BotToken(systemConfig.AppID, systemConfig.AppToken)

	var client openapi.OpenAPI

	if systemConfig.TestMode {
		client = botgo.NewSandboxOpenAPI(botToken).WithTimeout(time.Duration(systemConfig.TimeoutSecond) * time.Second)
	} else {
		client = botgo.NewOpenAPI(botToken).WithTimeout(time.Duration(systemConfig.TimeoutSecond) * time.Second)
	}

	ws, initWsErr := client.WS(context.Background(), nil, "")
	if initWsErr != nil {
		panic(initWsErr)
	}

	intent := websocket.RegisterHandlers(
		AtMessageEventHandler(client),
	)

	startBotErr := botgo.NewSessionManager().Start(ws, botToken, &intent)
	if startBotErr != nil {
		panic(startBotErr)
	}
}
