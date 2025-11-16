package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Mode    string `env:"MODE" envDefault:"dev"`
	Address string `env:"REVIEWER_ADDRESS" envDefault:":8080"`
	DB      DBConfig
}

type DBConfig struct {
	Name     string `env:"DB_NAME" envDefault:"pr_reviewer"`
	Host     string `env:"DB_HOST" envDefault:"postgres"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"password"`
}

func MustLoadConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		slog.Error("failed parse config", "error", err.Error())
		panic(err)
	}

	return &cfg
}

func (cfg *Config) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
}
