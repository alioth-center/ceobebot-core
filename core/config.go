package core

import (
	"os"
	"path/filepath"

	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/alioth-center/infrastructure/utils/shortcut"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot        BotConfig       `yaml:"bot" json:"bot"`
	Websocket  WebsocketConfig `yaml:"websocket" json:"websocket"`
	Plugins    []PluginConfig  `yaml:"plugins" json:"plugins,omitempty"`
	ZeroConfig *zero.Config    `yaml:"-" json:"-"`
}

type BotConfig struct {
	Nickname      []string `yaml:"nickname" json:"nickname,omitempty"`
	TriggerPrefix string   `yaml:"trigger_prefix" json:"trigger_prefix,omitempty"`
	SupperUsers   []int64  `yaml:"supper_users" json:"supper_users,omitempty"`
	Logger        string   `yaml:"logger" json:"logger,omitempty"`
	Debug         bool     `yaml:"debug" json:"debug,omitempty"`
}

type WebsocketConfig struct {
	Host        string `yaml:"host" json:"host,omitempty"`
	Port        int    `yaml:"port" json:"port,omitempty"`
	AccessToken string `yaml:"access_token" json:"access_token,omitempty"`
}

type (
	PluginConfig struct {
		Name           string           `yaml:"name" json:"name,omitempty"`
		Description    string           `yaml:"description" json:"description,omitempty"`
		Help           string           `yaml:"help" json:"help,omitempty"`
		Enable         bool             `yaml:"enable" json:"enable,omitempty"`
		Banner         string           `yaml:"banner" json:"banner,omitempty"`
		ConfigFile     string           `yaml:"config_file" json:"config_file,omitempty"`
		ResourceFolder string           `yaml:"resource_folder" json:"resource_folder,omitempty"`
		DataFolder     string           `yaml:"data_folder" json:"data_folder,omitempty"`
		Priority       int              `yaml:"priority" json:"priority,omitempty"`
		Middlewares    MiddlewareConfig `yaml:"middlewares" json:"middlewares"`
		Handlers       []HandlerConfig  `yaml:"handlers" json:"handlers,omitempty"`
	}

	HandlerConfig struct {
		Name     string        `yaml:"name" json:"name,omitempty"`
		Blocked  bool          `yaml:"blocked" json:"blocked,omitempty"`
		Limiter  string        `yaml:"limiter" json:"limiter,omitempty"`
		Rules    []string      `yaml:"rules" json:"rules,omitempty"`
		Triggers TriggerConfig `yaml:"triggers" json:"triggers"`
	}

	MiddlewareConfig struct {
		PreHandlers []string `yaml:"pre_handlers" json:"pre_handlers,omitempty"`
		MidHandlers []string `yaml:"mid_handlers" json:"mid_handlers,omitempty"`
	}

	TriggerConfig struct {
		FullMatches []string `yaml:"full_matches" json:"full_matches,omitempty"`
		KeyWords    []string `yaml:"key_words" json:"key_words,omitempty"`
		Commands    []string `yaml:"commands" json:"commands,omitempty"`
		Prefixes    []string `yaml:"prefixes" json:"prefixes,omitempty"`
		Suffixes    []string `yaml:"suffixes" json:"suffixes,omitempty"`
		Regexes     []string `yaml:"regexes" json:"regexes,omitempty"`
		Notice      bool     `yaml:"notice" json:"notice,omitempty"`
	}
)

func loadConfig(metadata *PluginConfig) {
	// get config file path
	filename := shortcut.Ternary(filepath.Ext(metadata.ConfigFile) == "", metadata.ConfigFile+".yaml", metadata.ConfigFile)
	path := filepath.Join("./", "config", filename)

	// check file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	// check receiver exist
	pluginConfig, existConfig := plugins.Get(metadata.Name)
	if !existConfig || pluginConfig == nil || pluginConfig.config == nil {
		panic("config receiver not found: " + metadata.Name)
	}

	// Load config file
	file, openErr := os.OpenFile(path, os.O_RDONLY, 0o644)
	if openErr != nil {
		panic("failed to open config file: " + path)
	}
	defer file.Close()

	// unmarshal config file
	decodeErr := yaml.NewDecoder(file).Decode(pluginConfig.config)
	if decodeErr != nil {
		panic("failed to decode config file: " + path)
	}
}
