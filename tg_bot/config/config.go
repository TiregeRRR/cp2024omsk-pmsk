package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MinioEndpoint        string
	MinioAccessKey       string
	MinioSecretAccessKey string

	PostgresUsername string
	PostgresPassword string
	PostgresAddress  string
	PostgresDatabase string

	WhisperAddr  string
	ReporterAddr string
	LlamaAddr    string
}

func New() (Config, error) {
	var conf Config

	if err := envconfig.Process("bot", &conf); err != nil {
		return Config{}, fmt.Errorf("proccesing env failed: %w", err)
	}

	return conf, nil
}
