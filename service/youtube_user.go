package service

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const ServiceNameYouTubeUser = "YouTube User Info"

// YouTubeUserInfo 存储用户基本信息
type YouTubeUserInfo struct {
	UserID        string
	ChannelName   string
	SubscriberCount string
	VideoCount    string
	ViewCount     string
	Description   string
	AvatarURL     string
	ChannelURL    string
}

// YouTubeVideoInfo 存储视频信息
type YouTubeVideoInfo struct {
	VideoID     string
	Title       string
	ViewCount   string
	LikeCount   string
	UploadDate  string
	Duration    string
	Description string
	ThumbnailURL string
	VideoURL    string
}

// GetYouTubeUserInfo 根据用户ID或频道ID获取YouTube用户信息
func GetYouTubeUserInfo(userID string) (*YouTubeUserInfo, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	userInfo := &YouTubeUserInfo{
		UserID: userID,
	}

	// 处理不同格式的用户ID
	var channelURL string
	if strings.HasPrefix(userID, "UC") && len(userID) == 24 {
		// 频道ID格式 (UCxxxxxxxxxxxxxxxxxxxx)
		channelURL = fmt.Sprintf("https://www.youtube.com/channel/%s", userID)
	} else if strings.HasPrefix(userID, "@") {
		// 新格式用户名 (@username)
		channelURL = fmt.Sprintf("https://www.youtube.com/%s", userID)
	} else {
		// 传统用户名或其他格式
		channelURL = fmt.Sprintf("https://www.youtube.com/c/%s", userID)
	}

	userInfo.ChannelURL = channelURL

	// 获取频道基本信息
	c.OnHTML("meta[name='description']", func(e *colly.HTMLElement) {
		userInfo.Description = e.Attr("content")
	})

	// 获取频道名称
	c.OnHTML("meta[property='og:title']", func(e *colly.HTMLElement) {
		title := e.Attr("content")
		if title != "" {
			userInfo.ChannelName = title
		}
	})

	// 获取头像URL
	c.OnHTML("meta[property='og:image']", func(e *colly.HTMLElement) {
		userInfo.AvatarURL = e.Attr("content")
	})

	// 获取订阅数、视频数等信息（通过页面脚本获取）
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		
		// 尝试提取订阅数
		if strings.Contains(scriptContent, "subscriberCountText") {
			re := regexp.MustCompile(`"subscriberCountText":\{"simpleText":"([^"]+)"`)
			matches := re.FindStringSubmatch(scriptContent)
			if len(matches) > 1 {
				userInfo.SubscriberCount = matches[1]
			}
		}

		// 尝试提取视频数
		if strings.Contains(scriptContent, "videosCountText") {
			re := regexp.MustCompile(`"videosCountText":\{"runs":\[\{"text":"([^"]+)"`)
			matches := re.FindStringSubmatch(scriptContent)
			if len(matches) > 1 {
				userInfo.VideoCount = matches[1]
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("YouTube scraping error", "url", r.Request.URL, "error", err)
	})

	err := c.Visit(channelURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to visit YouTube channel: %s", channelURL)
	}

	return userInfo, nil
}

// GetYouTubeUserVideos 获取用户最新的视频列表
func GetYouTubeUserVideos(userID string, limit int) ([]YouTubeVideoInfo, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	var videos []YouTubeVideoInfo
	
	// 构建视频列表页面URL
	var videosURL string
	if strings.HasPrefix(userID, "UC") && len(userID) == 24 {
		videosURL = fmt.Sprintf("https://www.youtube.com/channel/%s/videos", userID)
	} else if strings.HasPrefix(userID, "@") {
		videosURL = fmt.Sprintf("https://www.youtube.com/%s/videos", userID)
	} else {
		videosURL = fmt.Sprintf("https://www.youtube.com/c/%s/videos", userID)
	}

	// 通过script标签获取视频信息
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		
		// 使用正则表达式提取视频信息
		if strings.Contains(scriptContent, "ytInitialData") {
			// 提取视频ID
			videoIDRegex := regexp.MustCompile(`"videoId":"([^"]+)"`)
			videoIDs := videoIDRegex.FindAllStringSubmatch(scriptContent, limit)
			
			// 提取视频标题
			titleRegex := regexp.MustCompile(`"title":\{"runs":\[\{"text":"([^"]+)"`)
			titles := titleRegex.FindAllStringSubmatch(scriptContent, limit)
			
			// 提取视频观看次数
			viewRegex := regexp.MustCompile(`"shortViewCountText":\{"simpleText":"([^"]+)"`)
			views := viewRegex.FindAllStringSubmatch(scriptContent, limit)
			
			// 提取发布时间
			publishRegex := regexp.MustCompile(`"publishedTimeText":\{"simpleText":"([^"]+)"`)
			publishTimes := publishRegex.FindAllStringSubmatch(scriptContent, limit)

			// 组合视频信息
			maxLen := len(videoIDs)
			if len(titles) < maxLen {
				maxLen = len(titles)
			}
			if limit > 0 && maxLen > limit {
				maxLen = limit
			}

			for i := 0; i < maxLen; i++ {
				video := YouTubeVideoInfo{}
				
				if i < len(videoIDs) && len(videoIDs[i]) > 1 {
					video.VideoID = videoIDs[i][1]
					video.VideoURL = fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.VideoID)
					video.ThumbnailURL = fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", video.VideoID)
				}
				
				if i < len(titles) && len(titles[i]) > 1 {
					video.Title = titles[i][1]
				}
				
				if i < len(views) && len(views[i]) > 1 {
					video.ViewCount = views[i][1]
				}
				
				if i < len(publishTimes) && len(publishTimes[i]) > 1 {
					video.UploadDate = publishTimes[i][1]
				}

				if video.VideoID != "" {
					videos = append(videos, video)
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("YouTube videos scraping error", "url", r.Request.URL, "error", err)
	})

	err := c.Visit(videosURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to visit YouTube videos page: %s", videosURL)
	}

	return videos, nil
}

// GetYouTubeVideoDetails 获取单个视频的详细信息（包括点赞数等）
func GetYouTubeVideoDetails(videoID string) (*YouTubeVideoInfo, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	video := &YouTubeVideoInfo{
		VideoID:      videoID,
		VideoURL:     fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		ThumbnailURL: fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID),
	}

	// 获取视频标题
	c.OnHTML("meta[property='og:title']", func(e *colly.HTMLElement) {
		video.Title = e.Attr("content")
	})

	// 获取视频描述
	c.OnHTML("meta[property='og:description']", func(e *colly.HTMLElement) {
		video.Description = e.Attr("content")
	})

	// 通过script获取详细信息
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		
		// 提取观看次数
		if strings.Contains(scriptContent, "viewCount") {
			re := regexp.MustCompile(`"viewCount":"([^"]+)"`)
			matches := re.FindStringSubmatch(scriptContent)
			if len(matches) > 1 {
				video.ViewCount = matches[1]
			}
		}

		// 提取点赞数
		if strings.Contains(scriptContent, "defaultText") {
			// YouTube的点赞数格式可能会变化，这里是一个通用的匹配
			re := regexp.MustCompile(`"defaultText":{"accessibility":{"accessibilityData":{"label":"([0-9,]+) likes"`)
			matches := re.FindStringSubmatch(scriptContent)
			if len(matches) > 1 {
				video.LikeCount = matches[1]
			}
		}

		// 提取视频时长
		if strings.Contains(scriptContent, "lengthSeconds") {
			re := regexp.MustCompile(`"lengthSeconds":"([^"]+)"`)
			matches := re.FindStringSubmatch(scriptContent)
			if len(matches) > 1 {
				if seconds, err := strconv.Atoi(matches[1]); err == nil {
					video.Duration = formatDuration(seconds)
				}
			}
		}

		// 提取发布日期
		if strings.Contains(scriptContent, "publishDate") {
			re := regexp.MustCompile(`"publishDate":"([^"]+)"`)
			matches := re.FindStringSubmatch(scriptContent)
			if len(matches) > 1 {
				video.UploadDate = matches[1]
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("YouTube video details scraping error", "url", r.Request.URL, "error", err)
	})

	err := c.Visit(video.VideoURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to visit YouTube video: %s", video.VideoURL)
	}

	return video, nil
}

// formatDuration 将秒数转换为时分秒格式
func formatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

// FormatYouTubeUserMessage 格式化YouTube用户信息为消息
func FormatYouTubeUserMessage(userInfo *YouTubeUserInfo, videos []YouTubeVideoInfo) []string {
	var messages []string
	var message strings.Builder

	// 用户基本信息
	message.WriteString(fmt.Sprintf("# YouTube 用户信息\n"))
	message.WriteString(fmt.Sprintf("**频道名称**: %s\n", userInfo.ChannelName))
	message.WriteString(fmt.Sprintf("**用户ID**: %s\n", userInfo.UserID))
	message.WriteString(fmt.Sprintf("**订阅数**: %s\n", userInfo.SubscriberCount))
	message.WriteString(fmt.Sprintf("**视频总数**: %s\n", userInfo.VideoCount))
	message.WriteString(fmt.Sprintf("**频道链接**: %s\n", userInfo.ChannelURL))
	if userInfo.Description != "" {
		message.WriteString(fmt.Sprintf("**简介**: %s\n", userInfo.Description))
	}
	message.WriteString("\n")

	messages = append(messages, message.String())
	message.Reset()

	// 最新视频信息
	if len(videos) > 0 {
		message.WriteString("## 最新视频\n")
		for idx, video := range videos {
			if idx >= 10 { // 限制显示数量
				break
			}
			message.WriteString(fmt.Sprintf("### %d. [%s](%s)\n", idx+1, video.Title, video.VideoURL))
			if video.ViewCount != "" {
				message.WriteString(fmt.Sprintf("**观看次数**: %s\n", video.ViewCount))
			}
			if video.LikeCount != "" {
				message.WriteString(fmt.Sprintf("**点赞数**: %s\n", video.LikeCount))
			}
			if video.UploadDate != "" {
				message.WriteString(fmt.Sprintf("**发布时间**: %s\n", video.UploadDate))
			}
			if video.Duration != "" {
				message.WriteString(fmt.Sprintf("**时长**: %s\n", video.Duration))
			}
			message.WriteString("\n")

			// 每5个视频分一条消息
			if (idx+1)%5 == 0 {
				messages = append(messages, message.String())
				message.Reset()
			}
		}

		if left := message.String(); left != "" {
			messages = append(messages, left)
		}
	}

	return messages
}

// SendYouTubeUserInfo 获取并发送YouTube用户信息到Discord
func SendYouTubeUserInfo(userID string) error {
	// 获取用户信息
	userInfo, err := GetYouTubeUserInfo(userID)
	if err != nil {
		return errors.Wrapf(err, "failed to get YouTube user info for %s", userID)
	}

	// 获取最新视频
	videos, err := GetYouTubeUserVideos(userID, 10)
	if err != nil {
		slog.Warn("Failed to get YouTube videos", "userID", userID, "error", err)
		videos = []YouTubeVideoInfo{} // 继续处理，但没有视频信息
	}

	// 获取视频详细信息（包括点赞数）
	for i := range videos {
		if videoDetails, err := GetYouTubeVideoDetails(videos[i].VideoID); err == nil {
			videos[i].LikeCount = videoDetails.LikeCount
			videos[i].Duration = videoDetails.Duration
			if videos[i].UploadDate == "" {
				videos[i].UploadDate = videoDetails.UploadDate
			}
		}
		// 添加延迟避免被限制
		time.Sleep(100 * time.Millisecond)
	}

	// 格式化消息
	messages := FormatYouTubeUserMessage(userInfo, videos)

	// 发送到Discord（这里需要配置Discord webhook）
	// 注意：实际使用时需要确保config包中有相应的配置
	// return SendMessageToDiscord(messages, config.GetDiscordChatWebhookUrl(), ServiceNameYouTubeUser)
	
	// 临时：打印到日志
	for _, msg := range messages {
		slog.Info("YouTube User Info", "message", msg)
	}

	return nil
}
