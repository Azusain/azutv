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

type AzutvTaskType string

const (
	AzutvTaskTypeOriconRanking   AzutvTaskType = "oricon_ranking"
	AzutvTaskTypeGithubTrending  AzutvTaskType = "github_trending"
	AzutvTaskTypeVocaloidRanking AzutvTaskType = "vocaloid_ranking"
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
		config.GetDiscordChatWebhookUrl(),
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
		slog.Warn(errors.Wrapf(err, "failed to send Oricon Ranking to Discord").Error())
		return
	}
}

func sendVocaloidRanking() {
	messages, err := service.GetVocaloidRankingMessage()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	if err := SendMessageToDiscord(
		messages,
		config.GetDiscordChatWebhookUrl(),
		service.ServiceNameVocaloidnRanking,
	); err != nil {
		slog.Warn(errors.Wrapf(err, "failed to send Vocaloid Ranking to Discord").Error())
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

	switch AzutvTaskType(*task) {
	case AzutvTaskTypeOriconRanking:
		sendOriconRanking()

	case AzutvTaskTypeGithubTrending:
		sendGithubTrending()

	case AzutvTaskTypeVocaloidRanking:
		sendVocaloidRanking()

	default:
		slog.Error(fmt.Sprintf("invalid task type %q", *task))
		return
	}

	slog.Info(fmt.Sprintf("task %s completed", *task))
}
