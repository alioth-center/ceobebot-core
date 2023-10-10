package log

type Config struct {
	Level     string `yaml:"level"`
	Formatter string `yaml:"formatter"`
	FilePath  string `yaml:"log_path"`
}
