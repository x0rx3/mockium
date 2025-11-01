package ports

type Config struct {
	DefaultLogger LoggerConfig            `yaml:"default_logger"`
	Servers       map[string]ServerConfig `yaml:"servers"`
}

type LoggerConfig struct {
	Type     string `yaml:"type"`
	Level    string `yaml:"level"`
	FilePath string `yaml:"file_path"`
}

type ServerConfig struct {
	Port   string `yaml:"port"`
	IP     string `yaml:"ip"`
	Logger string `yaml:"logger"`
}
