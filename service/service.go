package service

import (
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

func RunService(task *string) {
	switch AzutvTaskType(*task) {
	case AzutvTaskTypeOriconRanking:
		SendOriconRanking()

	case AzutvTaskTypeGithubTrending:
		SendGithubTrending()

	case AzutvTaskTypeVocaloidRanking:
		SendVocaloidRanking()

	default:
		slog.Error(fmt.Sprintf("invalid task type %q", *task))
		return
	}
}
