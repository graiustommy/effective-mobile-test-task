package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PostgresConfig `envPrefix:"db_"`
}

type PostgresConfig struct {
	Host     string `env:"host" envDefault:"localhost"`
	Port     int    `env:"port" envDefault:"5432"`
	User     string `env:"user" envDefault:"postgres"`
	Password string `env:"password" envDefault:""`
	DBName   string `env:"name" envDefault:"db"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load(".env")
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
