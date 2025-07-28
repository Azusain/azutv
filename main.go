package main

import (
	"azuserver/config"
	"azuserver/service"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/gtuk/discordwebhook"
	"github.com/pkg/errors"
)

func SendMessageToDiscord(message string, channelUrl string) error {
	dcMessage := discordwebhook.Message{
		Content: &message,
	}
	return discordwebhook.SendMessage(channelUrl, dcMessage)
}

func main() {
	if err := config.LoadConfig(); err != nil {
		slog.Error(err.Error())
		return
	}

	// setup scheduler.
	tokyoLoc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.Error("failed to load location")
		return
	}
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(tokyoLoc))
	if err != nil {
		slog.Error(errors.Wrapf(err, "failed to init scheduler").Error())
		return
	}

	// Oricon Ranking.
	sendOriconRanking := func() {
		oriconRankMessage, err := service.GetOriconRankingDataMessage()
		if err != nil {
			slog.Warn(err.Error())
			return
		}
		if err := SendMessageToDiscord(oriconRankMessage, config.GetDiscordChatWebhookUrl()); err != nil {
			slog.Warn(errors.Wrapf(err, "failed to send Oricon Raning to Discord").Error())
			return
		}
	}
	if _, err := scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(10, 0, 0))),
		gocron.NewTask(sendOriconRanking),
	); err != nil {
		slog.Error(err.Error())
		return
	}

	// Heartz.
	sendHeartz := func() {
		SendMessageToDiscord("Δzu Server is up and running...", config.GetDiscordSystWebhookUrl())
	}
	if _, err := scheduler.NewJob(
		gocron.DurationJob(time.Hour*1),
		gocron.NewTask(sendHeartz),
	); err != nil {
		slog.Error(err.Error())
		return
	}

	// Github Trending.
	sendGithubTrending := func() {
		githubTrendingMessage, err := service.GetGithubTrendingMessage()
		if err != nil {
			slog.Warn(errors.Wrapf(err, "failed to get Github Trending").Error())
			return
		}
		if err := SendMessageToDiscord(githubTrendingMessage, config.GetDiscordChatWebhookUrl()); err != nil {
			slog.Warn(errors.Wrapf(err, "failed to send Github Trending to Discord").Error())
			return
		}
	}
	if _, err := scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(10, 0, 0))),
		gocron.NewTask(sendGithubTrending),
	); err != nil {
		slog.Error(err.Error())
		return
	}

	// start scheduler and block forever.
	scheduler.Start()
	slog.Info("Δzu Server is running...")
	select {}
}
