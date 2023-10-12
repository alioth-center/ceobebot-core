package log

type Config struct {
	Level        string `yaml:"level"`
	Formatter    string `yaml:"formatter"`
	FilePath     string `yaml:"log_path"`
	RelativePath bool   `yaml:"relative_path"`
	PackageName  string `yaml:"package_name"`
}
