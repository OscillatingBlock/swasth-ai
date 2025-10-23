package config

import (
	"errors"
	"github.com/spf13/viper"
	"log/slog"
)

type Config struct {
	Server     Server
	Db         Database
	JWT        JWT
	SMS        SMS
	Redis      RedisConfig
	Bun        BunConfig
	LoggerMode LoggerMode
}

type Server struct {
	Port        string
	Environment string
}
type Database struct {
	PostgresDSN string
	RedisAddr   string
}
type JWT struct {
	Secret    string
	ExpiresIn int
}
type SMS struct {
	Provider string
	APIKey   string
}

type BunConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type LoggerMode struct {
	Development bool
	Prod        bool
	Level       string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		slog.Error("unable to decode into struct", "err", err)
		return nil, err
	}

	return &c, nil
}
