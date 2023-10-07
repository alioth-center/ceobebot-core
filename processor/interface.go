package processor

// MessageHandleFunction 消息处理函数
type MessageHandleFunction func(ctx Context)

// MustHandleFunction 必须被处理的消息处理函数
type MustHandleFunction MessageHandleFunction

// OptionalHandleFunction 可以不被处理的消息处理函数
type OptionalHandleFunction MessageHandleFunction
