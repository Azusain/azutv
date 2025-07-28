package main

import (
	"azuserver/config"
	"azuserver/service"
	"log/slog"

	"github.com/gtuk/discordwebhook"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		slog.Error(err.Error())
		return
	}

	oriconRankMessage, err := service.GetOriconRankingDataMessage()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	message := discordwebhook.Message{
		Content: &oriconRankMessage,
	}
	if err := discordwebhook.SendMessage(config.GetDiscordWebhookUrl(), message); err != nil {
		slog.Warn(err.Error())
		return
	}
}
