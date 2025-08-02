package service

import (
	"azuserver/config"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func SendVocaloidRanking() {
	messages, err := GetVocaloidRankingMessage()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	if err := SendMessageToDiscord(
		messages,
		config.GetDiscordChatWebhookUrl(),
		ServiceNameVocaloidnRanking,
	); err != nil {
		slog.Warn(errors.Wrapf(err, "failed to send Vocaloid Ranking to Discord").Error())
		return
	}
}

const ServiceNameVocaloidnRanking = "Vocaloid Ranking"

type VocaloidRankingEntry struct {
	Name   string `json:"name"`
	Artist string `json:"artistString"`
	ID     int    `json:"id"`
	Url    string
}

type PvEntry struct {
	Service string `json:"service"`
	Url     string `json:"url"`
}

type PvServiceEntrySong struct {
	Pvs []PvEntry `json:"pvs"`
}

type PvServiceEntry struct {
	Song PvServiceEntrySong `json:"song"`
}

// Get Youtube PV Link.
func fetchPvYoutubeLinkById(id string) (string, error) {
	var pvServiceEntry PvServiceEntry
	client := resty.New()
	resp, err := client.R().
		SetResult(&pvServiceEntry).
		Get(fmt.Sprintf("https://vocadb.net/api/songs/%s/with-rating", id))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get 'pvService' for song: %s", id)
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("error status code %q", resp.StatusCode())
	}

	var nicoUrl, youtubeUrl, bandcampUrl string

	for _, pv := range pvServiceEntry.Song.Pvs {

		switch pv.Service {
		case "Youtube":
			youtubeUrl = pv.Url
		case "NicoNicoDouga":
			nicoUrl = pv.Url
		case "Bandcamp":
			bandcampUrl = pv.Url
		default:
			continue
		}
	}

	if youtubeUrl != "" {
		return youtubeUrl, nil
	} else if nicoUrl != "" {
		return nicoUrl, nil
	} else if bandcampUrl != "" {
		return bandcampUrl, nil
	}

	return "", errors.New(fmt.Sprintf("Youtube url not found for song: %s", id))
}

func fetchPvYouTubeLinks(ctx context.Context, entries []VocaloidRankingEntry) error {
	eg, ctx := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(5)
	for idx, e := range entries {
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}
		eg.Go(func() error {
			defer sem.Release(1)
			url, _ := fetchPvYoutubeLinkById(strconv.Itoa(e.ID))
			entries[idx].Url = url
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func GetVocaloidRankingMessage() ([]string, error) {
	var entries []VocaloidRankingEntry
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			// TODO: weekly and overall.
			"durationHours": "24",
			"filterBy":      "CreateDate",
		}).
		SetResult(&entries).
		Get("https://vocadb.net/api/songs/top-rated")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send fetch Vocaloid ranking data")
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error status code %q", resp.StatusCode())
	}

	err = fetchPvYouTubeLinks(context.Background(), entries)
	if err != nil {
		return nil, err
	}

	var vocaloidRank strings.Builder
	var messages []string
	for idx, entry := range entries {
		if entry.Url == "" {
			vocaloidRank.WriteString(fmt.Sprintf("%d. %s - %s\n", idx+1, entry.Name, entry.Artist))
		} else {
			vocaloidRank.WriteString(fmt.Sprintf("%d. [%s](<%s>) - %s\n", idx+1, entry.Name, entry.Url, entry.Artist))
		}
		if (idx+1)%10 == 0 {
			messages = append(messages, vocaloidRank.String())
			vocaloidRank.Reset()
		}
	}
	left := vocaloidRank.String()
	if left != "" {
		messages = append(messages, left)
	}

	return messages, nil
}
