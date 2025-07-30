package main

import (
	"azuserver/config"
	"azuserver/service"
	"flag"
	"fmt"
	"log/slog"

	"github.com/gtuk/discordwebhook"
	"github.com/pkg/errors"
)

func SendMessageToDiscord(message string, channelUrl string, username string) error {
	dcMessage := discordwebhook.Message{
		Username: &username,
		Content:  &message,
	}
	return discordwebhook.SendMessage(channelUrl, dcMessage)
}

func sendGithubTrending() {
	githubTrendingMessage, err := service.GetGithubTrendingMessage()
	if err != nil {
		slog.Warn(errors.Wrapf(err, "failed to get Github Trending").Error())
		return
	}
	if err := SendMessageToDiscord(githubTrendingMessage, config.GetDiscordChatWebhookUrl(), service.ServiceNameGithubTrending); err != nil {
		slog.Warn(errors.Wrapf(err, "failed to send Github Trending to Discord").Error())
		return
	}
}

// Oricon Ranking.
func sendOriconRanking() {
	oriconRankMessage, err := service.GetOriconRankingDataMessage()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	if err := SendMessageToDiscord(oriconRankMessage, config.GetDiscordSystWebhookUrl(), service.ServiceNameOriconRanking); err != nil {
		slog.Warn(errors.Wrapf(err, "failed to send Oricon Raning to Discord").Error())
		return
	}
}

func main() {
	// TODO: Load plugin configuration only if the plugin is enabled.
	if err := config.LoadConfig(); err != nil {
		slog.Error(err.Error())
		return
	}

	task := flag.String("task", "", "oricon_rank, github_trending")
	flag.Parse()

	switch *task {
	case "oricon_rank":
		sendOriconRanking()

	case "github_trending":
		sendGithubTrending()

	default:
		slog.Error(fmt.Sprintf("invalid task type %q", *task))
		return
	}

	slog.Info(fmt.Sprintf("task %s completed", *task))
}
