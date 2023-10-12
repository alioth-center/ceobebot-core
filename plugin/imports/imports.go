package imports

// 将需要的插件在这里引入
import (
	_ "github.com/ceobebot/qqchannel/plugin/chat"    // gpt聊天插件
	_ "github.com/ceobebot/qqchannel/plugin/example" // 示例插件
	_ "github.com/ceobebot/qqchannel/plugin/mihoyo"  // 米哈游相关插件
	_ "github.com/ceobebot/qqchannel/plugin/zerobot" // zbp机器人插件
)
