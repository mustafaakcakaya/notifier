package main

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	BotToken string `yaml:"botToken"`
	ChatID   int64  `yaml:"chatID"`
}

func GetSecretConfigs() Config {
	viper.BindEnv("BOT_TOKEN")
	viper.BindEnv("CHAT_ID")

	botToken := viper.GetString("BOT_TOKEN")
	chatID := viper.GetInt64("CHAT_ID")

	if botToken == "" || chatID == 0 {
		log.Fatalf("Config is not filled: %s", "empty-config-parameters")
	}

	return Config{
		BotToken: botToken,
		ChatID:   chatID,
	}
}
