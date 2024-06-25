package driver

import (
	"context"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/alioth-center/ceobebot-core/core"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/utils/values"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

func InitializeZeroBot(ctx context.Context, coreConfig *core.Config, pluginConfigMap map[string]*core.PluginConfig) {
	// inject hard coded priority
	filtered := values.FilterArray(coreConfig.Plugins, func(cfg core.PluginConfig) bool { return cfg.Enable && len(cfg.Handlers) > 0 })
	for _, plugin := range filtered {
		pluginConfig, existPluginConfig := core.Components.Plugins().Get(plugin.Name)
		if !existPluginConfig {
			// skip plugin without config
			continue
		}

		if pluginConfig.Priority() != 0 {
			// replace config priority with hard coded priority
			plugin.Priority = pluginConfig.Priority()
			pluginConfigMap[plugin.Name] = &plugin
		}
	}

	// sort handlers by priority
	items := values.SortArray(filtered, func(a, b core.PluginConfig) bool { return a.Priority < b.Priority })

	// init plugins
	for _, item := range items {
		pluginBuffer, existPluginBuffer := core.Components.Plugins().Get(item.Name)
		if !existPluginBuffer {
			// skip plugin without register in manager
			continue
		}

		pluginBuffer.Init()
		core.Logger().Debug(logger.NewFields(ctx).WithMessage("plugin initialized").WithData(map[string]any{"plugin": item.Name}))
	}

	// register plugins
	for _, item := range items {
		core.Logger().Debug(logger.NewFields(ctx).WithMessage("registering plugin handlers").WithData(map[string]any{"plugin": item.Name, "handlers": len(item.Handlers)}))

		endpoints := map[string]func(*zero.Ctx){}
		for _, handler := range item.Handlers {
			impl, got := core.Components.Handlers().Get(handler.Name)
			if got && impl != nil {
				endpoints[handler.Name] = impl
			}
		}
		if len(endpoints) == 0 {
			// skip plugin without handlers
			continue
		}

		// init plugin control engine
		engine := control.Register(item.Name, &ctrl.Options[*zero.Ctx]{
			Brief:             item.Description,
			Help:              item.Help,
			Banner:            item.Banner,
			PublicDataFolder:  item.ResourceFolder,
			PrivateDataFolder: item.DataFolder,
			OnEnable:          enableCallback(item),
			OnDisable:         disableCallback(item),
		})

		// bind middlewares
		for _, middleware := range item.Middlewares.PreHandlers {
			impl, got := core.Components.Middlewares().Get(middleware)
			if got && impl != nil {
				engine.UsePreHandler(impl)
			}
		}
		for _, middleware := range item.Middlewares.MidHandlers {
			impl, got := core.Components.Middlewares().Get(middleware)
			if got && impl != nil {
				engine.UseMidHandler(impl)
			}
		}

		// register handlers
		bindHandler(ctx, engine, item, endpoints)
	}

	// lock components
	core.Components.Done()

	// serving bot
	serve(ctx, coreConfig)
}

func bindHandler(ctx context.Context, engine *control.Engine, plugin core.PluginConfig, endpoints map[string]func(*zero.Ctx)) {
	for _, handler := range plugin.Handlers {
		if endpoints[handler.Name] == nil {
			core.Logger().Info(logger.NewFields(ctx).WithMessage("handler not found").WithData(handler.Name))
			continue
		}

		extraRules, limiter := findRules(ctx, handler.Rules), findLimiter(ctx, handler.Limiter)
		if len(handler.Triggers.FullMatches) > 0 {
			engine.OnFullMatchGroup(handler.Triggers.FullMatches, extraRules...).
				SetBlock(handler.Blocked).Limit(limiter).Handle(endpoints[handler.Name])
		}
		if len(handler.Triggers.KeyWords) > 0 {
			engine.OnKeywordGroup(handler.Triggers.KeyWords, extraRules...).
				SetBlock(handler.Blocked).Limit(limiter).Handle(endpoints[handler.Name])
		}
		if len(handler.Triggers.Commands) > 0 {
			engine.OnCommandGroup(handler.Triggers.Commands, extraRules...).
				SetBlock(handler.Blocked).Limit(limiter).Handle(endpoints[handler.Name])
		}
		if len(handler.Triggers.Prefixes) > 0 {
			engine.OnPrefixGroup(handler.Triggers.Prefixes, extraRules...).
				SetBlock(handler.Blocked).Limit(limiter).Handle(endpoints[handler.Name])
		}
		if len(handler.Triggers.Suffixes) > 0 {
			engine.OnSuffixGroup(handler.Triggers.Suffixes, extraRules...).
				SetBlock(handler.Blocked).Limit(limiter).Handle(endpoints[handler.Name])
		}
		if handler.Triggers.Notice {
			engine.OnNotice(extraRules...).
				SetBlock(handler.Blocked).Limit(limiter).Handle(endpoints[handler.Name])
		}
		for _, regex := range handler.Triggers.Regexes {
			engine.OnRegex(regex, extraRules...).
				SetBlock(handler.Blocked).Limit(limiter).Handle(endpoints[handler.Name])
		}

		core.Logger().Debug(logger.NewFields(ctx).WithMessage("handler registered").WithData(map[string]any{"plugin": plugin.Name, "handler": handler.Name, "metadata": plugin.Handlers}))
	}
}

func findRules(ctx context.Context, rus []string) (result []zero.Rule) {
	result = []zero.Rule{}
	for _, rule := range rus {
		impl, got := core.Components.Rules().Get(rule)
		if !got || impl == nil {
			core.Logger().Info(logger.NewFields(ctx).WithMessage("rule not found").WithData(rule))
			continue
		}

		result = append(result, impl)
	}

	return result
}

func findLimiter(ctx context.Context, limiter string) func(*zero.Ctx) *rate.Limiter {
	impl, got := core.Components.Limiters().Get(limiter)
	if !got || impl == nil {
		core.Logger().Info(logger.NewFields(ctx).WithMessage("limiter not found, default use user limiter").WithData(limiter))
		return ctxext.LimitByUser
	}

	return impl
}

func enableCallback(plugin core.PluginConfig) func(_ *zero.Ctx) {
	return func(_ *zero.Ctx) {
		core.Logger().Info(logger.NewFields().WithMessage("plugin enabled").WithData(plugin.Name))
	}
}

func disableCallback(plugin core.PluginConfig) func(_ *zero.Ctx) {
	return func(_ *zero.Ctx) {
		core.Logger().Info(logger.NewFields().WithMessage("plugin disabled").WithData(plugin.Name))
	}
}

func serve(ctx context.Context, coreConfig *core.Config) {
	// start bot
	core.Logger().Infof(logger.NewFields(ctx), "startting bot, connecting to onebot adapter server: %s:%d", coreConfig.Websocket.Host, coreConfig.Websocket.Port)
	zero.Run(coreConfig.ZeroConfig)
	core.Logger().Info(logger.NewFields(ctx).WithMessage("bot started"))

	// register exit event
	exit.Register(func(_ string) string {
		// sleep 1 second to wait for bot exit
		time.Sleep(time.Second)
		return "bot exit"
	}, "bot exit")

	// wait for exit signal
	exit.BlockedUntilTerminate()
	core.Logger().Info(logger.NewFields(ctx).WithMessage("bot stopped"))
}
