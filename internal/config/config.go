package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	BindAddress string `config:"bindaddress"`
	LogLevel    int    `config:"loglevel"`
	Debug       bool   `config:"debug"`

	DiscordToken string `config:"discord.token"`

	RedisAddress  string `config:"redis.address"`
	RedisUsername string `config:"redis.username"`
	RedisPassword string `config:"redis.password"`
	RedisDB       int    `config:"redis.db"`
}

var defaultConfig = Config{
	BindAddress: "0.0.0.0:80",
	LogLevel:    int(logrus.InfoLevel),
	Debug:       false,
}

func Load() (c Config, err error) {
	godotenv.Load()

	c = defaultConfig

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewOptionalBackend("config.json"),
		file.NewOptionalBackend("config.yaml"),
		flags.NewBackend(),
	)
	err = loader.Load(context.Background(), &c)
	return
}
