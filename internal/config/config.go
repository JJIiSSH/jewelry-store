package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all runtime configuration, loaded from environment variables
// (12-factor style) with sensible defaults for local development.
type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN assembles the PostgreSQL connection string from the DB parts, e.g.
// "postgres://user:password@host:port/name?sslmode=disable".
func (c DBConfig) DSN() string {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)

	return connectionString
}

// Load reads configuration via Viper and returns a populated Config.
// It fails fast when a required secret (DB password) is missing.
func Load() (*Config, error) {
	var cfg Config

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_USER", "IvanDev")
	viper.SetDefault("DB_NAME", "mydb")

	viper.AutomaticEnv()

	cfg.Server.Port = viper.GetString("SERVER_PORT")
	cfg.DB.Host = viper.GetString("DB_HOST")
	cfg.DB.Port = viper.GetString("DB_PORT")
	cfg.DB.User = viper.GetString("DB_USER")
	cfg.DB.Password = viper.GetString("DB_PASSWORD")
	cfg.DB.Name = viper.GetString("DB_NAME")
	cfg.DB.SSLMode = viper.GetString("DB_SSLMODE")

	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return &cfg, nil
}
