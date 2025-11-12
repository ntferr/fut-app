package config

import (
	goenv "github.com/Netflix/go-env"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Host string `env:"APP_HOST,required"`
	Port string `env:"APP_PORT,required"`
	Name string `env:"APP_NAME,required"`
}

type PostgresConfig struct {
	Host     string `env:"PG_HOST,required"`
	Port     string `env:"PG_PORT,required"`
	User     string `env:"PG_USER,required"`
	Password string `env:"PG_PASSWORD,required"`
	Name     string `env:"PG_NAME,required"`
}

func NewConfig() (*Config, error) {
	var envs Config
	_, err := goenv.UnmarshalFromEnviron(&envs)
	return &envs, err
}
