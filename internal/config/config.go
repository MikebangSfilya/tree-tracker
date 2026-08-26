package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	APIHost     string `env:"API_HOST"`
	APIPort     int    `env:"API_PORT" env-default:"8080"`
	DataBaseURL string `env:"DATABASE_URL"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
