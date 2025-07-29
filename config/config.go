package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	YamlConfigPath = "./config.yaml"

	// Oricon.
	DomainOricon  = "www.oricon.co.jp"
	OriconRankUrl = "https://www.oricon.co.jp/rank/"
)

type Config struct {
	DiscordChatWebhookUrl string `yaml:"chat_webhook"`
	DiscordSysWebhookUrl  string `yaml:"system_webhook"`
}

var (
	appConfig Config
)

func GetDiscordChatWebhookUrl() string {
	return appConfig.DiscordChatWebhookUrl
}

func GetDiscordSystWebhookUrl() string {
	return appConfig.DiscordSysWebhookUrl
}

func LoadConfig() error {
	// from local file.
	if _, err := os.Stat(YamlConfigPath); os.IsExist(err) {
		file, err := os.Open(YamlConfigPath)
		if err != nil {
			return errors.Wrapf(err, "error loading config")
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&appConfig)
		if err != nil {
			return errors.Wrapf(err, "error decoding config")
		}
		slog.Info(fmt.Sprintf("loading configurations from local file: %q", YamlConfigPath))
		return nil
	}

	// or from shell env.
	// TODO: validate the args.
	appConfig.DiscordChatWebhookUrl = os.Getenv("DISCORD_CHAT_WEBHOOK_URL")
	appConfig.DiscordSysWebhookUrl = os.Getenv("DISCORD_SYS_WEBHOOK_URL")
	slog.Info("loading configurations from shell env")

	return nil
}
