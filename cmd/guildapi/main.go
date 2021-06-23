package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/guildapi/internal/config"
	"github.com/zekrotja/guildapi/internal/discord"
	"github.com/zekrotja/guildapi/internal/web"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatal("Failed loading config: ", err)
	}

	logrus.SetLevel(logrus.Level(cfg.LogLevel))
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: cfg.Debug,
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	dc, err := discord.New(cfg.DiscordToken, redisClient)
	if err != nil {
		logrus.Fatal("Failed iniitalizing discord session: ", err)
	}

	if err = dc.Open(); err != nil {
		logrus.Fatal("Failed opening discord session: ", err)
	}
	defer dc.Close()

	ws := web.New(dc, cfg.AllowOrigins)
	if err = ws.Open(cfg.BindAddress); err != nil {
		logrus.Fatal("Failed opening web server: ", err)
	}
}
