package service

import (
	"azuserver/config"
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
	AzutvTaskTypeYouTubeUser     AzutvTaskType = "youtube_user"
	AzutvTaskTypeBilibiliUser    AzutvTaskType = "bilibili_user"
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

	case AzutvTaskTypeYouTubeUser:
		userID := config.GetYouTubeDefaultUserID()
		if userID == "" {
			slog.Error("YouTube default user ID not configured")
			return
		}
		if err := SendYouTubeUserInfo(userID); err != nil {
			slog.Error("Failed to send YouTube user info", "error", err)
		}

	case AzutvTaskTypeBilibiliUser:
		uid := config.GetBilibiliDefaultUID()
		if uid == "" {
			slog.Error("Bilibili default UID not configured")
			return
		}
		if err := SendBilibiliUserInfo(uid); err != nil {
			slog.Error("Failed to send Bilibili user info", "error", err)
		}

	default:
		slog.Error(fmt.Sprintf("invalid task type %q", *task))
		return
	}
}

// RunServiceWithParams 运行需要参数的服务
func RunServiceWithParams(task string, params map[string]string) error {
	switch AzutvTaskType(task) {
	case AzutvTaskTypeYouTubeUser:
		userID, ok := params["userID"]
		if !ok || userID == "" {
			return errors.New("YouTube user service requires 'userID' parameter")
		}
		return SendYouTubeUserInfo(userID)

	case AzutvTaskTypeBilibiliUser:
		uid, ok := params["uid"]
		if !ok || uid == "" {
			return errors.New("Bilibili user service requires 'uid' parameter")
		}
		return SendBilibiliUserInfo(uid)

	default:
		return errors.Errorf("unsupported task type for parameterized service: %s", task)
	}
}
