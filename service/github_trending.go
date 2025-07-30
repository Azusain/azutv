package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const ServiceNameGithubTrending = "Github Trending"

type GithubTrendingEntry struct {
	Title       string
	Link        string
	Stars       string
	Language    string
	Description string
}

func GetGithubTrendingMessage() (string, error) {
	entries := []GithubTrendingEntry{}
	c := colly.NewCollector()

	c.OnHTML("article.Box-row", func(e *colly.HTMLElement) {
		repoPath := strings.TrimSpace(e.ChildAttr("h2 a", "href"))
		repoPath = strings.TrimPrefix(repoPath, "/")
		repoURL := "https://github.com/" + repoPath
		title := strings.ReplaceAll(repoPath, " ", "")
		stars := strings.TrimSpace(e.ChildText("a[href$='/stargazers']"))
		language := strings.TrimSpace(e.ChildText("span[itemprop='programmingLanguage']"))
		var description string
		e.DOM.Find("p").Each(func(_ int, s *goquery.Selection) {
			description = strings.TrimSpace(s.Text())
		})

		entries = append(entries, GithubTrendingEntry{
			Title:       title,
			Link:        repoURL,
			Stars:       stars,
			Language:    language,
			Description: description,
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
		message.WriteString(fmt.Sprintf("## \\#%d  [%s](%s)\n%s⭐ %s\n",
			idx+1, entry.Title, entry.Link, entry.Language, entry.Stars))
		message.WriteString(fmt.Sprintf("%s\n\n", entry.Description))
		if idx >= 5 {
			break
		}
	}

	return message.String(), nil
}
