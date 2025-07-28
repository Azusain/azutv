package config

import (
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
	DiscordWebhookUrl string `yaml:"discord_webhook"`
}

var (
	appConfig Config
)

func GetDiscordWebhookUrl() string {
	return appConfig.DiscordWebhookUrl
}

func LoadConfig() error {
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

	return nil
}
