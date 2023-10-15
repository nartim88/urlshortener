package config

type Config struct {
	RunAddr string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

var CFG Config

func NewCFG() *Config {
	cfg := Config{
		RunAddr: "localhost:8080",
		BaseURL: "http://localhost",
	}
	return &cfg
}

// InitConfigs инициализация конфигов приложения
func InitConfigs() {
	CFG = *NewCFG()

	parseFlags()
	parseEnv()
}
