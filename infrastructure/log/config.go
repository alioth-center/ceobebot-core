package log

type Config struct {
	Level     string `yaml:"level"`
	Formatter string `yaml:"formatter"`
	FilePath  string `yaml:"file_path"`
}
