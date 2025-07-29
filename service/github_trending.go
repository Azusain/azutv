package service

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const ServiceNameGithubTrending = "Github Trending"

type TrendingEntry struct {
	Title    string
	Link     string
	Stars    string
	Language string
}

func GetGithubTrendingMessage() (string, error) {
	entries := []TrendingEntry{}
	c := colly.NewCollector()

	c.OnHTML("article.Box-row", func(e *colly.HTMLElement) {
		repoPath := strings.TrimSpace(e.ChildAttr("h2 a", "href"))
		repoPath = strings.TrimPrefix(repoPath, "/")
		repoURL := "https://github.com/" + repoPath
		title := strings.ReplaceAll(repoPath, " ", "")
		stars := strings.TrimSpace(e.ChildText("a[href$='/stargazers']"))
		language := strings.TrimSpace(e.ChildText("span[itemprop='programmingLanguage']"))

		entries = append(entries, TrendingEntry{
			Title:    title,
			Link:     repoURL,
			Stars:    stars,
			Language: language,
		})
	})

	err := c.Visit("https://github.com/trending")
	if err != nil {
		return "", errors.Wrapf(err, "failed to visit Github trending")
	}

	var message strings.Builder
	for idx, entry := range entries {
		if entry.Language != "" {
			entry.Language = fmt.Sprintf("**%s** - ", entry.Language)
		}
		message.WriteString(fmt.Sprintf("## \\#%d  [%s](%s)\n%s⭐ %s \n\n",
			idx+1, entry.Title, entry.Link, entry.Language, entry.Stars))
	}

	return message.String(), nil
}
