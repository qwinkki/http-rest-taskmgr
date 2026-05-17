package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `env:"LOG_LEVEL" required:"true"`
	Folder string `env:"LOG_FOLDER" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("process env: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err := fmt.Errorf("new config: %w", err)
		panic(err)
	}
	return config
}
