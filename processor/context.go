package processor

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

type Payload struct {
	Message *dto.WSATMessageData
	Event   *dto.WSPayload
	Command string
	Content string
}

type Context interface {
	GetContext() context.Context
	SetContext(ctx context.Context)
	GetApi() Client
	GetPayload() Payload
	Next()
	Abort()
	IsAborted() bool
}

type defaultContext struct {
	ctx              context.Context
	api              Client
	mustHandlers     []MustHandleFunction
	optionalHandlers []OptionalHandleFunction
	idx              int
	mustIdx          bool
	payload          Payload
}

func (c *defaultContext) GetContext() context.Context {
	return c.ctx
}

func (c *defaultContext) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *defaultContext) GetApi() Client {
	return c.api
}

func (c *defaultContext) GetPayload() Payload {
	return c.payload
}

func (c *defaultContext) Next() {
	// 执行可选内容
	c.idx++
	for c.idx < len(c.optionalHandlers) {
		c.optionalHandlers[c.idx](c)
		c.idx++
	}

	// 执行必须内容
	if c.mustIdx {
		for _, handler := range c.mustHandlers {
			handler(c)
		}
		c.mustIdx = false
	}
}

func (c *defaultContext) Abort() {
	c.idx = len(c.optionalHandlers) + 1
}

func (c *defaultContext) IsAborted() bool {
	return c.idx > len(c.optionalHandlers)
}

func NewContext(ctx context.Context, api openapi.OpenAPI, payload Payload, functions []MustHandleFunction, handlers []OptionalHandleFunction) Context {
	if functions == nil {
		functions = []MustHandleFunction{}
	}

	if handlers == nil {
		handlers = []OptionalHandleFunction{}
	}

	return &defaultContext{
		ctx:              ctx,
		api:              NewClient(api),
		mustHandlers:     functions,
		optionalHandlers: handlers,
		idx:              -1,
		mustIdx:          true,
		payload:          payload,
	}
}
