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

func SendMessageToDiscord(messages []string, channelUrl string, username string) error {
	for _, m := range messages {
		dcMessage := discordwebhook.Message{
			Username: &username,
			Content:  &m,
		}
		if err := discordwebhook.SendMessage(channelUrl, dcMessage); err != nil {
			return errors.Wrapf(err, "failed to send message to Discord")
		}
	}
	return nil
}

func sendGithubTrending() {
	githubTrendingMessages, err := service.GetGithubTrendingMessage()
	if err != nil {
		slog.Warn(errors.Wrapf(err, "failed to get Github Trending").Error())
		return
	}
	if err := SendMessageToDiscord(
		githubTrendingMessages,
		config.GetDiscordSystWebhookUrl(),
		service.ServiceNameGithubTrending,
	); err != nil {
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
	if err := SendMessageToDiscord(
		[]string{oriconRankMessage},
		config.GetDiscordChatWebhookUrl(),
		service.ServiceNameOriconRanking,
	); err != nil {
		slog.Warn(errors.Wrapf(err, "failed to send Oricon Raning to Discord").Error())
		return
	}
}

// TODO: test bin options.
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
