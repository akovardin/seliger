package config

import (
	"os"

	"go.uber.org/config"
	"go.uber.org/fx"
)

type (
	Config struct {
		fx.Out

		Database Database

		Provider config.Provider `name:"config_provider"`
	}
)

type Database struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

func New(file string) (Config, error) {
	provider, err := config.NewYAML(
		config.Expand(os.LookupEnv),
		config.File(file),
		config.Permissive(),
	)

	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Provider: provider,
	}

	err = provider.Get("").Populate(&cfg)
	if err != nil {
		return Config{}, err
	}

	var database Database
	if err := provider.Get("db").Populate(&database); err != nil {
		return Config{}, err
	}

	cfg.Database = database

	return cfg, nil
}
