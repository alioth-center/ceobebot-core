package core

import (
	"github.com/alioth-center/infrastructure/utils/concurrency"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

type Bus interface {
	Handlers() concurrency.Map[string, func(*zero.Ctx)]
	Middlewares() concurrency.Map[string, func(*zero.Ctx) bool]
	Rules() concurrency.Map[string, func(*zero.Ctx) bool]
	Limiters() concurrency.Map[string, func(*zero.Ctx) *rate.Limiter]
	Plugins() concurrency.Map[string, *PluginOptions]
	Interfaces() concurrency.Map[string, any]
	Done()
}

type bus struct {
	initialized bool
}

func (b *bus) Handlers() concurrency.Map[string, func(*zero.Ctx)] {
	if !b.initialized {
		return handlers
	}

	return nil
}

func (b *bus) Middlewares() concurrency.Map[string, func(*zero.Ctx) bool] {
	if !b.initialized {
		return middlewares
	}

	return nil
}

func (b *bus) Rules() concurrency.Map[string, func(*zero.Ctx) bool] {
	if !b.initialized {
		return rules
	}

	return nil
}

func (b *bus) Limiters() concurrency.Map[string, func(*zero.Ctx) *rate.Limiter] {
	if !b.initialized {
		return limiters
	}

	return nil
}

func (b *bus) Plugins() concurrency.Map[string, *PluginOptions] {
	if !b.initialized {
		return plugins
	}

	return nil
}

func (b *bus) Interfaces() concurrency.Map[string, any] {
	if !b.initialized {
		return interfaces
	}

	return nil
}

func (b *bus) Done() {
	b.initialized = true
}

var Components Bus = &bus{}
