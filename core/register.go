package core

import (
	"reflect"

	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/values"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

var (
	handlers    = concurrency.NewMap[string, func(*zero.Ctx)]()
	middlewares = concurrency.NewMap[string, func(*zero.Ctx) bool]()
	rules       = concurrency.NewMap[string, func(*zero.Ctx) bool]()
	limiters    = concurrency.NewMap[string, func(*zero.Ctx) *rate.Limiter]()
	plugins     = concurrency.NewMap[string, *PluginOptions]()
	interfaces  = concurrency.NewMap[string, any]()
)

func RegisterHandler(name string, handler func(*zero.Ctx)) {
	_, exist := handlers.Get(name)
	if exist {
		// cannot rewrite handler, it will replace built-in handlers
		panic("handler already exist")
	}

	handlers.Set(name, handler)
}

func RegisterTriggerRule(name string, rule func(*zero.Ctx) bool) {
	_, exist := rules.Get(name)
	if exist {
		// cannot rewrite rule, it will replace built-in rules
		panic("rule already exist")
	}

	rules.Set(name, rule)
}

func RegisterLimiter(name string, limiter func(*zero.Ctx) *rate.Limiter) {
	_, exist := limiters.Get(name)
	if exist {
		// cannot rewrite limiter, it will replace built-in limiters
		panic("limiter already exist")
	}

	limiters.Set(name, limiter)
}

func RegisterMiddleware(name string, middleware func(*zero.Ctx) bool) {
	_, exist := middlewares.Get(name)
	if exist {
		// cannot rewrite middleware, it will replace built-in middlewares
		panic("middleware already exist")
	}

	middlewares.Set(name, middleware)
}

func RegisterPlugin(name string, options ...PluginOpts) {
	_, exist := plugins.Get(name)
	if exist {
		// cannot rewrite config, it will replace built-in plugins
		panic("config already exist")
	}

	// attach options
	opt := &PluginOptions{}
	for _, o := range options {
		o(opt)
	}

	plugins.Set(name, opt)
}

func RegisterInterface(name string, ifrace any) {
	_, exist := interfaces.Get(name)
	if exist {
		// cannot rewrite interface, it will replace built-in interfaces
		panic("interface already exist")
	}

	interfaces.Set(name, ifrace)
}

// GetIfrace get interface by name, if not exist, return nil interface
func GetIfrace[T any](name string) T {
	nilIfrace := values.Nil[T]()

	// get interface, if not exist, return nil interface
	ifrace, exist := interfaces.Get(name)
	if !exist {
		return nilIfrace
	}

	// check if interface is assignable to T, if not, return nil interface
	if !reflect.TypeOf(ifrace).AssignableTo(reflect.TypeOf(nilIfrace)) {
		return nilIfrace
	}

	// convert interface to T, if failed, return nil interface
	converted, convertSuccess := ifrace.(T)
	if !convertSuccess {
		return nilIfrace
	}

	return converted
}

type PluginOpts func(opt *PluginOptions)

// WithConfig set plugin config, must be a pointer which can be unmarshalled from yaml
func WithConfig(config any) PluginOpts {
	return func(opt *PluginOptions) {
		opt.config = config
	}
}

// WithInit set plugin init function, will be called when plugin loaded
func WithInit(init func()) PluginOpts {
	return func(opt *PluginOptions) {
		opt.init = init
	}
}

// WithPriority set plugin priority, it will replace priority in config file
func WithPriority(priority int) PluginOpts {
	return func(opt *PluginOptions) {
		opt.priority = priority
	}
}

type PluginOptions struct {
	config   any
	priority int
	init     func()
}

func (opts PluginOptions) Config() any {
	return opts.config
}

func (opts PluginOptions) Priority() int {
	return opts.priority
}

func (opts PluginOptions) Init() {
	if opts.init != nil {
		opts.init()
	}
}
