package service

import (
	"azuserver/config"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const ServiceNameOriconRanking = "Oricon Ranking"

type OriconRankingDataEntry struct {
	Title  string
	Artist string
	Link   string
}

type OriconRankingData struct {
	Rule    string
	Entries []OriconRankingDataEntry
}
type OriconRankingDataArray []OriconRankingData

func FetchRankingDataFromOricon() (OriconRankingDataArray, error) {
	const RankDataSelector = "#content-main > div.content-main-inner > div.content-rank-main > div > article > section:nth-child(2)"
	var retErr error
	c := colly.NewCollector(
		colly.AllowedDomains(config.DomainOricon),
	)

	getRankDataByRule := func(index int, rankData *OriconRankingData, e *colly.HTMLElement) {
		prefix := fmt.Sprintf("div:nth-child(%d) > ", index)

		// rule.
		e.DOM.Find(prefix + "h3").Each(func(i int, s *goquery.Selection) {
			rankData.Rule = s.Text()
		})

		// entries.
		e.DOM.Find(prefix + "div > div").Each(func(i int, s *goquery.Selection) {
			s.Find("dl").Each(func(i int, dl *goquery.Selection) {
				var title, artist, href string
				dl.Find("a").Each(func(_ int, a *goquery.Selection) {
					href, _ = a.Attr("href")
				})

				dl.Find("h4").Each(func(i int, h4 *goquery.Selection) {
					title = h4.Text()
				})
				dl.Find("p").Each(func(i int, p *goquery.Selection) {
					if p.HasClass("name") {
						artist = p.Text()
					}
				})
				rankData.Entries = append(rankData.Entries, OriconRankingDataEntry{
					Title:  title,
					Artist: artist,
					Link:   href,
				})
			})
		})
	}

	dailySingleRankData := OriconRankingData{}
	dailyAlbumRankData := OriconRankingData{}
	weeklySingleRankData := OriconRankingData{}
	weeklyAlbumRankData := OriconRankingData{}

	c.OnHTML(RankDataSelector, func(e *colly.HTMLElement) {
		getRankDataByRule(2, &dailySingleRankData, e)
		getRankDataByRule(6, &dailyAlbumRankData, e)
		getRankDataByRule(4, &weeklySingleRankData, e)
		getRankDataByRule(7, &weeklyAlbumRankData, e)
	})

	c.OnError(func(r *colly.Response, err error) {
		retErr = err
	})

	c.Visit(config.OriconRankUrl)

	return []OriconRankingData{
		dailySingleRankData,
		dailyAlbumRankData,
		weeklySingleRankData,
		weeklyAlbumRankData,
	}, retErr
}

func (oriconRankData OriconRankingDataArray) Dump() string {
	var oriconRank strings.Builder
	for _, data := range oriconRankData {
		oriconRank.WriteString(fmt.Sprintf("### %s\n", data.Rule))
		for _, entry := range data.Entries {
			if entry.Link == "" {
				oriconRank.WriteString(fmt.Sprintf("%s - %s\n", entry.Title, entry.Artist))
				continue
			}
			oriconRank.WriteString(fmt.Sprintf("[%s](https://%s/%s) - %s\n", entry.Title, config.DomainOricon, entry.Link, entry.Artist))
		}
		oriconRank.WriteRune('\n')
	}
	return oriconRank.String()
}

func GetOriconRankingDataMessage() (string, error) {
	rankData, err := FetchRankingDataFromOricon()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get ranking data from Oricon")
	}
	return rankData.Dump(), nil
}
