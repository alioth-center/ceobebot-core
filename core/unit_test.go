package core

import (
	"context"
	"github.com/alioth-center/infrastructure/trace"
	"os"
	"path/filepath"
	"testing"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/utils/concurrency"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"gopkg.in/yaml.v3"
)

func reset() {
	handlers = concurrency.NewMap[string, func(*zero.Ctx)]()
	middlewares = concurrency.NewMap[string, func(*zero.Ctx) bool]()
	rules = concurrency.NewMap[string, func(*zero.Ctx) bool]()
	limiters = concurrency.NewMap[string, func(*zero.Ctx) *rate.Limiter]()
	plugins = concurrency.NewMap[string, *PluginOptions]()
	interfaces = concurrency.NewMap[string, any]()
	coreConfig = &Config{}
	pluginConfigMap = map[string]*PluginConfig{}
	coreLogger = nil
	tempLogger = nil
}

func TestMain(m *testing.M) {
	// set environment variables
	_ = os.Setenv("ci", "true")

	// run testing
	code := m.Run()

	// exit
	os.Exit(code)
}

func TestConfig(t *testing.T) {
	_ = os.MkdirAll("config", os.ModePerm)
	defer func() {
		_ = os.Remove("config/test.yaml")
		_ = os.Remove("config/invalid.yaml")
		_ = os.Remove("config")
	}()

	type TestConfig struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Enable      bool   `yaml:"enable"`
	}

	t.Run("LoadConfig", func(t *testing.T) {
		// write config to correct path
		f, _ := os.Create("config/test.yaml")
		cfgIst := TestConfig{
			Name:        "test-plugin",
			Description: "This is a test plugin",
			Enable:      true,
		}

		// write config to file
		bytes, _ := yaml.Marshal(cfgIst)
		_, _ = f.Write(bytes)

		// close file
		_ = f.Close()

		readIst := TestConfig{}
		RegisterPlugin("test-plugin", WithConfig(&readIst))

		loadConfig(&PluginConfig{
			Name:       "test-plugin",
			ConfigFile: "test.yaml",
		})

		if readIst.Name != cfgIst.Name || readIst.Description != cfgIst.Description || readIst.Enable != cfgIst.Enable {
			t.Errorf("Expected config to be loaded, but it was not")
		}
	})

	t.Run("LoadConfigFileNotFound", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to missing config file, but did not panic")
			}
		}()

		loadConfig(&PluginConfig{
			Name:       "test-plugin",
			ConfigFile: "nonexistent.yaml",
		})
	})

	t.Run("LoadConfigReceiverNotFound", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to missing plugin config receiver, but did not panic")
			}
		}()

		loadConfig(&PluginConfig{
			Name:       "unknown-plugin",
			ConfigFile: "test.yaml",
		})
	})

	t.Run("LoadInvalidConfig", func(t *testing.T) {
		// write invalid config to correct path
		f, _ := os.Create("config/invalid.yaml")
		_, _ = f.Write([]byte(`name: ["1", "2"]`))

		// close file
		_ = f.Close()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid config file content, but did not panic")
			}
		}()

		RegisterPlugin("invalid-plugin", WithConfig(&TestConfig{}))

		loadConfig(&PluginConfig{
			Name:       "invalid-plugin",
			ConfigFile: "invalid.yaml",
		})
	})
}

func TestBuiltin(t *testing.T) {
	t.Run("InitRegister", func(t *testing.T) {
		initRegister()
	})
}

func TestBus(t *testing.T) {
	t.Run("NotInitialized", func(t *testing.T) {
		if Components.Handlers() == nil {
			t.Errorf("Expected handlers to be initialized, but it was not")
		}

		if Components.Middlewares() == nil {
			t.Errorf("Expected middlewares to be initialized, but it was not")
		}

		if Components.Rules() == nil {
			t.Errorf("Expected rules to be initialized, but it was not")
		}

		if Components.Limiters() == nil {
			t.Errorf("Expected limiters to be initialized, but it was not")
		}

		if Components.Plugins() == nil {
			t.Errorf("Expected plugins to be initialized, but it was not")
		}

		if Components.Interfaces() == nil {
			t.Errorf("Expected interfaces to be initialized, but it was not")
		}
	})

	t.Run("Initialized", func(t *testing.T) {
		Components.Done()

		if Components.Handlers() != nil {
			t.Errorf("Expected handlers to be nil, but it was initialized")
		}

		if Components.Middlewares() != nil {
			t.Errorf("Expected middlewares to be nil, but it was initialized")
		}

		if Components.Rules() != nil {
			t.Errorf("Expected rules to be nil, but it was initialized")
		}

		if Components.Limiters() != nil {
			t.Errorf("Expected limiters to be nil, but it was initialized")
		}

		if Components.Plugins() != nil {
			t.Errorf("Expected plugins to be nil, but it was initialized")
		}

		if Components.Interfaces() != nil {
			t.Errorf("Expected interfaces to be nil, but it was initialized")
		}
	})
}

func TestRegister(t *testing.T) {
	nilHandlerImpl := func(*zero.Ctx) {}
	nilMiddlewareImpl := func(*zero.Ctx) bool { return true }
	nilRuleImpl := func(*zero.Ctx) bool { return true }
	nilLimiterImppl := func(*zero.Ctx) *rate.Limiter { return nil }

	t.Run("RegisterHandlerSuccess", func(t *testing.T) {
		RegisterHandler("nil-handler", nilHandlerImpl)
	})

	t.Run("RegisterHandlerFailed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to duplicate handler registration, but did not panic")
			}
		}()

		RegisterHandler("nil-handler", nilHandlerImpl)
	})

	t.Run("RegisterMiddlewareSuccess", func(t *testing.T) {
		RegisterMiddleware("nil-middleware", nilMiddlewareImpl)
	})

	t.Run("RegisterMiddlewareFailed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to duplicate middleware registration, but did not panic")
			}
		}()

		RegisterMiddleware("nil-middleware", nilMiddlewareImpl)
	})

	t.Run("RegisterRuleSuccess", func(t *testing.T) {
		RegisterTriggerRule("nil-rule", nilRuleImpl)
	})

	t.Run("RegisterRuleFailed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to duplicate rule registration, but did not panic")
			}
		}()

		RegisterTriggerRule("nil-rule", nilRuleImpl)
	})

	t.Run("RegisterLimiterSuccess", func(t *testing.T) {
		RegisterLimiter("nil-limiter", nilLimiterImppl)
	})

	t.Run("RegisterLimiterFailed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to duplicate limiter registration, but did not panic")
			}
		}()

		RegisterLimiter("nil-limiter", nilLimiterImppl)
	})

	t.Run("RegisterPluginSuccess", func(t *testing.T) {
		RegisterPlugin("nil-plugin")
	})

	t.Run("RegisterPluginFailed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to duplicate plugin registration, but did not panic")
			}
		}()

		RegisterPlugin("nil-plugin")
	})

	t.Run("RegisterInterfaceSuccess", func(t *testing.T) {
		RegisterInterface("nil-interface", nil)
	})

	t.Run("RegisterInterfaceFailed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to duplicate interface registration, but did not panic")
			}
		}()

		RegisterInterface("nil-interface", nil)
	})

	t.Run("ProcessPluginOpts", func(t *testing.T) {
		x, y := 0, 0
		nilInit := func() {
			x = 1
		}
		nilInitCtx := func(ctx context.Context) {
			if trace.GetTid(ctx) == "" {
				t.Errorf("Expected context to be set, but it was not")
			} else {
				y = 1
			}
		}
		opts := &PluginOptions{}
		WithConfig(1)(opts)
		if opts.Config() != 1 {
			t.Errorf("Expected config to be set, but it was not")
		}

		WithPriority(1)(opts)
		if opts.Priority() != 1 {
			t.Errorf("Expected priority to be set, but it was not")
		}

		WithInit(nilInit)(opts)
		WithInitCtx(nilInitCtx)(opts)
		opts.Init()
		opts.InitCtx(trace.NewContext())
		if x != 1 || y != 1 {
			t.Errorf("Expected init function to be set, but it was not")
		}
	})

	t.Run("GetIfraceNotExist", func(t *testing.T) {
		type emptyStruct struct {
			name string
		}

		x := GetIfrace[emptyStruct]("not-exist-ifrace")
		if x.name != "" {
			t.Errorf("Expected interface to be nil, but it was not")
		}
	})

	t.Run("GetIfraceNotAssignable", func(t *testing.T) {
		RegisterInterface("not-assignable-ifrace", 1)
		x := GetIfrace[float64]("not-assignable-ifrace")
		if x != 0 {
			t.Errorf("Expected interface to be nil, but it was not")
		}
	})

	t.Run("GetIfraceSuccess", func(t *testing.T) {
		RegisterInterface("assignable-ifrace", 1)
		x := GetIfrace[int]("assignable-ifrace")
		if x != 1 {
			t.Errorf("Expected interface to be set, but it was not")
		}
	})
}

func TestManager(t *testing.T) {
	t.Run("InitializeWithConfigFile", func(t *testing.T) {
		reset()

		// Setup environment for this test
		_ = os.Setenv("ci", "false")
		_ = os.MkdirAll("config", os.ModePerm)
		defer func() {
			_ = os.RemoveAll("config")
			_ = os.RemoveAll("logs")
		}()

		// Create a valid config file
		configPath := filepath.Join("config", "bot.yaml")
		configContent := `
bot:
  nickname: ["小刻"]
  trigger_prefix: ""
  supper_users: [1145141919]
  logger: "console"
websocket:
  host: "localhost"
  port: 8080
  access_token: "token"
plugins:
  - name: "test-plugin"
    description: "This is a test plugin"
    enable: true
    config_file: ""
`
		_ = os.WriteFile(configPath, []byte(configContent), os.ModePerm)
		ctx, cfg, mapping := Initialize()

		// Validate the results
		if cfg.Bot.Nickname[0] != "小刻" || cfg.Websocket.Host != "localhost" {
			t.Errorf("Expected config to be loaded correctly, but it was not")
		}

		if len(mapping) != 1 || mapping["test-plugin"].Name != "test-plugin" {
			t.Errorf("Expected plugin mapping to be loaded correctly, but it was not")
		}

		if coreLogger == nil {
			t.Errorf("Expected logger to be initialized, but it was not")
		}

		if ctx == nil {
			t.Errorf("Expected context to be initialized, but it was not")
		}
	})

	t.Run("InitializeWithoutConfigFile", func(t *testing.T) {
		reset()

		// Setup environment for this test
		_ = os.Setenv("ci", "false")
		_ = os.RemoveAll("config")
		defer func() {
			_ = os.RemoveAll("config")
			_ = os.RemoveAll("logs")
		}()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to missing config file, but did not panic")
			}
		}()

		Initialize()
	})

	t.Run("InitializeInvalidConfigFile", func(t *testing.T) {
		reset()

		// Setup environment for this test
		_ = os.Setenv("ci", "false")
		_ = os.MkdirAll("config", os.ModePerm)
		defer func() {
			_ = os.RemoveAll("config")
			_ = os.RemoveAll("logs")
		}()

		// Create an invalid config file
		configPath := filepath.Join("config", "bot.yaml")
		_ = os.WriteFile(configPath, []byte("invalid yaml content"), os.ModePerm)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid config file, but did not panic")
			}
		}()

		Initialize()
	})

	t.Run("InitializeWithCustomLogger", func(t *testing.T) {
		reset()

		// Setup environment for this test
		_ = os.Setenv("ci", "false")
		_ = os.MkdirAll("config", os.ModePerm)
		defer func() {
			_ = os.RemoveAll("config")
			_ = os.RemoveAll("logs")
		}()

		// Create a valid config file with custom logger
		configPath := filepath.Join("config", "bot.yaml")
		configContent := `
bot:
  nickname: ["小刻"]
  trigger_prefix: ""
  supper_users: [1145141919]
  logger: "custom"
websocket:
  host: "localhost"
  port: 8080
  access_token: "token"
plugins:
  - name: "test-plugin-2"
    description: "This is a test plugin"
    enable: true
    config_file: ""
`
		_ = os.WriteFile(configPath, []byte(configContent), os.ModePerm)

		customLogger := logger.NewLoggerWithConfig(logger.Config{
			Level:     "info",
			Formatter: "text",
		})
		SetLogger(customLogger)

		ctx, cfg, mapping := Initialize()

		// Validate the results
		if coreLogger != customLogger {
			t.Errorf("Expected custom logger to be set, but it was not")
		}

		if cfg.Bot.Nickname[0] != "小刻" || cfg.Websocket.Host != "localhost" {
			t.Errorf("Expected config to be loaded correctly, but it was not")
		}

		if len(mapping) != 1 {
			t.Log(mapping)
			t.Errorf("Expected plugin mapping to be loaded correctly, but it was not")
		}

		if ctx == nil {
			t.Errorf("Expected context to be initialized, but it was not")
		}
	})

	t.Run("SetLoggerWhenAlreadySet", func(t *testing.T) {
		reset()

		customLogger := logger.NewLoggerWithConfig(logger.Config{
			Level:     "info",
			Formatter: "text",
		})
		SetLogger(customLogger)

		// 尝试再次设置 logger
		anotherLogger := logger.NewLoggerWithConfig(logger.Config{
			Level:     "debug",
			Formatter: "json",
		})
		SetLogger(anotherLogger)

		if tempLogger != customLogger {
			t.Errorf("Expected tempLogger to remain the same, but it was changed")
		}
	})

	t.Run("InitializeWithDebugMode", func(t *testing.T) {
		reset()

		// Setup environment for this test
		_ = os.Setenv("ci", "false")
		_ = os.MkdirAll("config", os.ModePerm)
		defer func() {
			_ = os.RemoveAll("config")
			_ = os.RemoveAll("logs")
		}()

		// Create a valid config file with debug mode enabled
		configPath := filepath.Join("config", "bot.yaml")
		configContent := `
bot:
  nickname: ["小刻"]
  trigger_prefix: ""
  supper_users: [1145141919]
  logger: "console"
  debug: true
websocket:
  host: "localhost"
  port: 8080
  access_token: "token"
plugins:
  - name: "test-plugin"
    description: "This is a test plugin"
    enable: true
    config_file: ""
`
		_ = os.WriteFile(configPath, []byte(configContent), os.ModePerm)
		ctx, cfg, mapping := Initialize()

		// Validate the results
		if cfg.Bot.Nickname[0] != "小刻" || cfg.Websocket.Host != "localhost" {
			t.Errorf("Expected config to be loaded correctly, but it was not")
		}

		if len(mapping) != 1 || mapping["test-plugin"].Name != "test-plugin" {
			t.Errorf("Expected plugin mapping to be loaded correctly, but it was not")
		}

		if coreLogger == nil {
			t.Errorf("Expected logger to be initialized, but it was not")
		}

		if ctx == nil {
			t.Errorf("Expected context to be initialized, but it was not")
		}
	})

	t.Run("InitializeWithFileLogger", func(t *testing.T) {
		reset()

		// Setup environment for this test
		_ = os.Setenv("ci", "false")
		_ = os.MkdirAll("config", os.ModePerm)
		defer func() {
			_ = os.RemoveAll("config")
			_ = os.RemoveAll("logs")
		}()

		// Create a valid config file with file logger
		configPath := filepath.Join("config", "bot.yaml")
		configContent := `
bot:
  nickname: ["小刻"]
  trigger_prefix: ""
  supper_users: [1145141919]
  logger: "file"
websocket:
  host: "localhost"
  port: 8080
  access_token: "token"
plugins:
  - name: "test-plugin"
    description: "This is a test plugin"
    enable: true
    config_file: ""
`
		_ = os.WriteFile(configPath, []byte(configContent), os.ModePerm)
		ctx, cfg, mapping := Initialize()

		// Validate the results
		if cfg.Bot.Nickname[0] != "小刻" || cfg.Websocket.Host != "localhost" {
			t.Errorf("Expected config to be loaded correctly, but it was not")
		}

		if len(mapping) != 1 || mapping["test-plugin"].Name != "test-plugin" {
			t.Errorf("Expected plugin mapping to be loaded correctly, but it was not")
		}

		if coreLogger == nil {
			t.Errorf("Expected logger to be initialized, but it was not")
		}

		if ctx == nil {
			t.Errorf("Expected context to be initialized, but it was not")
		}
	})

	t.Run("LoadConfigsWithPluginConfigFile", func(t *testing.T) {
		reset()

		// Setup environment for this test
		_ = os.Setenv("ci", "false")
		_ = os.MkdirAll("config", os.ModePerm)
		defer func() {
			_ = os.RemoveAll("config")
			_ = os.RemoveAll("logs")
		}()

		// Create a valid config file with a plugin config file
		configPath := filepath.Join("config", "bot.yaml")
		configContent := `
bot:
  nickname: ["小刻"]
  trigger_prefix: ""
  supper_users: [1145141919]
  logger: "console"
websocket:
  host: "localhost"
  port: 8080
  access_token: "token"
plugins:
  - name: "test-plugin"
    description: "This is a test plugin"
    enable: true
    config_file: "plugin.yaml"
`
		_ = os.WriteFile(configPath, []byte(configContent), os.ModePerm)

		RegisterPlugin("test-plugin", WithConfig(&map[string]string{}))

		// Create the plugin config file
		pluginConfigPath := filepath.Join("config", "plugin.yaml")
		pluginConfigContent := `
key: value
`
		_ = os.WriteFile(pluginConfigPath, []byte(pluginConfigContent), os.ModePerm)

		ctx, cfg, mapping := Initialize()

		// Validate the results
		if cfg.Bot.Nickname[0] != "小刻" || cfg.Websocket.Host != "localhost" {
			t.Errorf("Expected config to be loaded correctly, but it was not")
		}

		if len(mapping) != 1 || mapping["test-plugin"].Name != "test-plugin" {
			t.Errorf("Expected plugin mapping to be loaded correctly, but it was not")
		}

		if coreLogger == nil {
			t.Errorf("Expected logger to be initialized, but it was not")
		}

		if ctx == nil {
			t.Errorf("Expected context to be initialized, but it was not")
		}
	})
}
