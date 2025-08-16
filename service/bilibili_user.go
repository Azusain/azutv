package service

import (
	"azuserver/config"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const ServiceNameBilibiliUser = "Bilibili User Info"

// BilibiliUserInfo 存储用户基本信息
type BilibiliUserInfo struct {
	UserID      string
	Username    string
	FollowerCount int64    // 粉丝数
	FollowingCount int64   // 关注数
	LikeCount   int64     // 获赞数
	PlayCount   int64     // 播放数
	VideoCount  int       // 视频数
	Description string
	AvatarURL   string
	SpaceURL    string
	Level       int
	VipType     int       // 0:无 1:月度 2:年度
}

// BilibiliVideoInfo 存储视频信息  
type BilibiliVideoInfo struct {
	BvID        string
	AvID        string
	Title       string
	ViewCount   int64
	LikeCount   int64
	CoinCount   int64     // 投币数
	FavoriteCount int64   // 收藏数
	ShareCount  int64     // 分享数
	ReplyCount  int64     // 评论数
	UploadDate  string
	Duration    string
	Description string
	CoverURL    string
	VideoURL    string
	Author      string
}

// GetBilibiliUserInfo 根据用户UID获取Bilibili用户信息
func GetBilibiliUserInfo(uid string) (*BilibiliUserInfo, error) {
	userInfo := &BilibiliUserInfo{
		UserID:   uid,
		SpaceURL: fmt.Sprintf("https://space.bilibili.com/%s", uid),
	}

	// 首先尝试通过页面获取基本信息（用户名、头像等）
	if err := getBilibiliUserBasicInfo(userInfo); err != nil {
		slog.Warn("Failed to get basic user info from page", "uid", uid, "error", err)
	}

	// 通过API获取统计信息
	if err := getBilibiliUserStatsFromAPI(userInfo); err != nil {
		slog.Warn("Failed to get user stats from API", "uid", uid, "error", err)
	}

	return userInfo, nil
}

// getBilibiliUserBasicInfo 从用户页面获取基本信息
func getBilibiliUserBasicInfo(userInfo *BilibiliUserInfo) error {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// 从页面标题获取用户名
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text
		if strings.Contains(title, "的个人空间") {
			parts := strings.Split(title, "的个人空间")
			if len(parts) > 0 && userInfo.Username == "" {
				userInfo.Username = parts[0]
			}
		}
	})

	// 从描述获取信息
	c.OnHTML("meta[name='description']", func(e *colly.HTMLElement) {
		desc := e.Attr("content")
		if desc != "" && userInfo.Description == "" {
			if strings.Contains(desc, "个人空间，提供") {
				parts := strings.Split(desc, "个人空间，提供")
				if len(parts) > 1 {
					detail := parts[1]
					if strings.Contains(detail, "内容，关注") {
						contentParts := strings.Split(detail, "内容，关注")
						if len(contentParts) > 0 {
							userInfo.Description = strings.TrimSpace(contentParts[0])
						}
					}
				}
			}
		}
	})

	// 从脚本获取基本信息
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		if strings.Contains(scriptContent, "__INITIAL_STATE__") || strings.Contains(scriptContent, "_render_data_") {
			// 提取用户名
			if strings.Contains(scriptContent, "\"name\":") && userInfo.Username == "" {
				re := regexp.MustCompile(`"name":"([^"]+)"`)
				matches := re.FindAllStringSubmatch(scriptContent, -1)
				for _, match := range matches {
					if len(match) > 1 && match[1] != "" {
						userInfo.Username = match[1]
						break
					}
				}
			}

			// 提取头像URL
			if strings.Contains(scriptContent, "\"face\":") {
				re := regexp.MustCompile(`"face":"([^"]+)"`)
				matches := re.FindStringSubmatch(scriptContent)
				if len(matches) > 1 {
					userInfo.AvatarURL = matches[1]
				}
			}

			// 提取用户等级
			if strings.Contains(scriptContent, "\"level\":") {
				re := regexp.MustCompile(`"level":(\d+)`)
				matches := re.FindStringSubmatch(scriptContent)
				if len(matches) > 1 {
					if level, err := strconv.Atoi(matches[1]); err == nil {
						userInfo.Level = level
					}
				}
			}

			// 提取简介
			if strings.Contains(scriptContent, "\"sign\":") {
				re := regexp.MustCompile(`"sign":"([^"]+)"`)
				matches := re.FindStringSubmatch(scriptContent)
				if len(matches) > 1 {
					userInfo.Description = matches[1]
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Bilibili basic info scraping error", "url", r.Request.URL, "error", err)
	})

	return c.Visit(userInfo.SpaceURL)
}

// getBilibiliUserStatsFromAPI 通过API获取用户统计信息
func getBilibiliUserStatsFromAPI(userInfo *BilibiliUserInfo) error {
	// 创建一个新的 collector 来调用 API
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// 设置请求头，模拟浏览器请求
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", userInfo.SpaceURL)
		r.Headers.Set("Accept", "application/json, text/plain, */*")
	})

	// 处理API响应
	c.OnResponse(func(r *colly.Response) {
		var apiResp struct {
			Code int `json:"code"`
			Data struct {
				Mid       int64  `json:"mid"`
				Name      string `json:"name"`
				Sex       string `json:"sex"`
				Face      string `json:"face"`
				Sign      string `json:"sign"`
				Level     int    `json:"level"`
				Follower  int64  `json:"follower"`
				Following int64  `json:"following"`
			} `json:"data"`
		}

		if err := json.Unmarshal(r.Body, &apiResp); err != nil {
			slog.Error("Failed to parse API response", "error", err)
			return
		}

		if apiResp.Code == 0 {
			// 更新用户信息
			if userInfo.Username == "" {
				userInfo.Username = apiResp.Data.Name
			}
			if userInfo.AvatarURL == "" {
				userInfo.AvatarURL = apiResp.Data.Face
			}
			if userInfo.Description == "" {
				userInfo.Description = apiResp.Data.Sign
			}
			if userInfo.Level == 0 {
				userInfo.Level = apiResp.Data.Level
			}
			userInfo.FollowerCount = apiResp.Data.Follower
			userInfo.FollowingCount = apiResp.Data.Following

			slog.Info("Successfully got user stats from API",
				"uid", userInfo.UserID,
				"name", userInfo.Username,
				"followers", userInfo.FollowerCount,
				"following", userInfo.FollowingCount)
		} else {
			slog.Warn("API returned error", "code", apiResp.Code)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("API request error", "url", r.Request.URL, "error", err)
	})

	// 构建API URL - 使用 Bilibili 的用户信息API
	apiURL := fmt.Sprintf("https://api.bilibili.com/x/space/wbi/acc/info?mid=%s", userInfo.UserID)
	return c.Visit(apiURL)
}

// GetBilibiliUserVideos 获取用户最新的视频列表
func GetBilibiliUserVideos(uid string, limit int) ([]BilibiliVideoInfo, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	var videos []BilibiliVideoInfo
	videosURL := fmt.Sprintf("https://space.bilibili.com/%s/video", uid)

	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		
		// 在页面脚本中查找视频数据
		if strings.Contains(scriptContent, "__INITIAL_STATE__") || strings.Contains(scriptContent, "window.__INITIAL_STATE__") {
			// 提取视频列表数据
			if strings.Contains(scriptContent, "\"list\":[") && strings.Contains(scriptContent, "\"bvid\":") {
				// 提取BV号
				bvidRegex := regexp.MustCompile(`"bvid":"([^"]+)"`)
				bvids := bvidRegex.FindAllStringSubmatch(scriptContent, limit)
				
				// 提取标题
				titleRegex := regexp.MustCompile(`"title":"([^"]+)"`)
				titles := titleRegex.FindAllStringSubmatch(scriptContent, limit)
				
				// 提取播放数
				playRegex := regexp.MustCompile(`"play":(\d+)`)
				plays := playRegex.FindAllStringSubmatch(scriptContent, limit)
				
				// 提取创建时间（Unix时间戳）
				createdRegex := regexp.MustCompile(`"created":(\d+)`)
				createds := createdRegex.FindAllStringSubmatch(scriptContent, limit)
				
				// 提取时长（秒）
				lengthRegex := regexp.MustCompile(`"length":"([^"]+)"`)
				lengths := lengthRegex.FindAllStringSubmatch(scriptContent, limit)
				
				// 提取描述
				descRegex := regexp.MustCompile(`"description":"([^"]*)"`)
				descs := descRegex.FindAllStringSubmatch(scriptContent, limit)
				
				// 提取封面图
				picRegex := regexp.MustCompile(`"pic":"([^"]+)"`)
				pics := picRegex.FindAllStringSubmatch(scriptContent, limit)

				// 组合视频信息
				maxLen := len(bvids)
				if limit > 0 && maxLen > limit {
					maxLen = limit
				}

				for i := 0; i < maxLen && i < len(titles); i++ {
					video := BilibiliVideoInfo{}
					
					if len(bvids[i]) > 1 {
						video.BvID = bvids[i][1]
						video.VideoURL = fmt.Sprintf("https://www.bilibili.com/video/%s", video.BvID)
					}
					
					if len(titles[i]) > 1 {
						video.Title = titles[i][1]
					}
					
					if i < len(plays) && len(plays[i]) > 1 {
						if count, err := strconv.ParseInt(plays[i][1], 10, 64); err == nil {
							video.ViewCount = count
						}
					}
					
					if i < len(createds) && len(createds[i]) > 1 {
						if timestamp, err := strconv.ParseInt(createds[i][1], 10, 64); err == nil {
							video.UploadDate = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
						}
					}
					
					if i < len(lengths) && len(lengths[i]) > 1 {
						video.Duration = lengths[i][1]
					}
					
					if i < len(descs) && len(descs[i]) > 1 {
						video.Description = descs[i][1]
					}
					
					if i < len(pics) && len(pics[i]) > 1 {
						video.CoverURL = "https:" + pics[i][1]
					}

					if video.BvID != "" {
						videos = append(videos, video)
					}
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Bilibili videos scraping error", "url", r.Request.URL, "error", err)
	})

	err := c.Visit(videosURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to visit Bilibili videos page: %s", videosURL)
	}

	return videos, nil
}

// GetBilibiliVideoDetails 获取单个视频的详细信息
func GetBilibiliVideoDetails(bvid string) (*BilibiliVideoInfo, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	video := &BilibiliVideoInfo{
		BvID:     bvid,
		VideoURL: fmt.Sprintf("https://www.bilibili.com/video/%s", bvid),
	}

	// 获取视频详细信息
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		
		if strings.Contains(scriptContent, "__INITIAL_STATE__") {
			// 提取视频统计数据
			if strings.Contains(scriptContent, "\"stat\":") {
				// 播放数
				if re := regexp.MustCompile(`"view":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
							video.ViewCount = count
						}
					}
				}
				
				// 点赞数
				if re := regexp.MustCompile(`"like":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
							video.LikeCount = count
						}
					}
				}
				
				// 投币数
				if re := regexp.MustCompile(`"coin":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
							video.CoinCount = count
						}
					}
				}
				
				// 收藏数
				if re := regexp.MustCompile(`"favorite":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
							video.FavoriteCount = count
						}
					}
				}
				
				// 分享数
				if re := regexp.MustCompile(`"share":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
							video.ShareCount = count
						}
					}
				}
				
				// 评论数
				if re := regexp.MustCompile(`"reply":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if count, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
							video.ReplyCount = count
						}
					}
				}
			}

			// 提取视频基本信息
			if strings.Contains(scriptContent, "\"videoData\":") {
				// 标题
				if re := regexp.MustCompile(`"title":"([^"]+)"`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						video.Title = matches[1]
					}
				}
				
				// 描述
				if re := regexp.MustCompile(`"desc":"([^"]*)"`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						video.Description = matches[1]
					}
				}
				
				// 封面
				if re := regexp.MustCompile(`"pic":"([^"]+)"`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						video.CoverURL = "https:" + matches[1]
					}
				}
				
				// 时长
				if re := regexp.MustCompile(`"duration":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if seconds, err := strconv.Atoi(matches[1]); err == nil {
							video.Duration = formatDuration(seconds)
						}
					}
				}
				
				// UP主
				if re := regexp.MustCompile(`"owner":\{[^}]*"name":"([^"]+)"`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						video.Author = matches[1]
					}
				}
				
				// 发布时间
				if re := regexp.MustCompile(`"pubdate":(\d+)`); re != nil {
					if matches := re.FindStringSubmatch(scriptContent); len(matches) > 1 {
						if timestamp, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
							video.UploadDate = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
						}
					}
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Bilibili video details scraping error", "url", r.Request.URL, "error", err)
	})

	err := c.Visit(video.VideoURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to visit Bilibili video: %s", video.VideoURL)
	}

	return video, nil
}

// parseBilibiliCount 解析B站的计数格式（如1.2万、234.5万）
func parseBilibiliCount(countStr string) (int64, error) {
	countStr = strings.TrimSpace(countStr)
	
	// 处理万、亿单位
	if strings.Contains(countStr, "万") {
		numStr := strings.Replace(countStr, "万", "", -1)
		if num, err := strconv.ParseFloat(numStr, 64); err == nil {
			return int64(num * 10000), nil
		}
	} else if strings.Contains(countStr, "亿") {
		numStr := strings.Replace(countStr, "亿", "", -1)
		if num, err := strconv.ParseFloat(numStr, 64); err == nil {
			return int64(num * 100000000), nil
		}
	} else {
		// 直接解析数字
		if num, err := strconv.ParseInt(countStr, 10, 64); err == nil {
			return num, nil
		}
	}
	
	return 0, errors.New("failed to parse count string: " + countStr)
}

// FormatBilibiliUserMessage 格式化Bilibili用户信息为消息
func FormatBilibiliUserMessage(userInfo *BilibiliUserInfo, videos []BilibiliVideoInfo) []string {
	var messages []string
	var message strings.Builder

	// 用户基本信息
	message.WriteString(fmt.Sprintf("# Bilibili 用户信息\n"))
	message.WriteString(fmt.Sprintf("**用户名**: %s\n", userInfo.Username))
	message.WriteString(fmt.Sprintf("**UID**: %s\n", userInfo.UserID))
	message.WriteString(fmt.Sprintf("**粉丝数**: %s\n", formatCount(userInfo.FollowerCount)))
	message.WriteString(fmt.Sprintf("**关注数**: %s\n", formatCount(userInfo.FollowingCount)))
	message.WriteString(fmt.Sprintf("**获赞数**: %s\n", formatCount(userInfo.LikeCount)))
	message.WriteString(fmt.Sprintf("**播放数**: %s\n", formatCount(userInfo.PlayCount)))
	message.WriteString(fmt.Sprintf("**等级**: Lv.%d\n", userInfo.Level))
	if userInfo.VipType > 0 {
		vipText := "月度大会员"
		if userInfo.VipType == 2 {
			vipText = "年度大会员"
		}
		message.WriteString(fmt.Sprintf("**会员类型**: %s\n", vipText))
	}
	message.WriteString(fmt.Sprintf("**个人空间**: %s\n", userInfo.SpaceURL))
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
			if video.ViewCount > 0 {
				message.WriteString(fmt.Sprintf("**播放量**: %s\n", formatCount(video.ViewCount)))
			}
			if video.LikeCount > 0 {
				message.WriteString(fmt.Sprintf("**点赞数**: %s\n", formatCount(video.LikeCount)))
			}
			if video.CoinCount > 0 {
				message.WriteString(fmt.Sprintf("**投币数**: %s\n", formatCount(video.CoinCount)))
			}
			if video.FavoriteCount > 0 {
				message.WriteString(fmt.Sprintf("**收藏数**: %s\n", formatCount(video.FavoriteCount)))
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

// formatCount 格式化数字显示（转换为万、亿等单位）
func formatCount(count int64) string {
	if count >= 100000000 { // 1亿
		return fmt.Sprintf("%.1f亿", float64(count)/100000000)
	} else if count >= 10000 { // 1万
		return fmt.Sprintf("%.1f万", float64(count)/10000)
	}
	return fmt.Sprintf("%d", count)
}

// SendBilibiliUserInfo 获取并发送Bilibili用户信息到Discord
func SendBilibiliUserInfo(uid string) error {
	// 获取用户信息
	userInfo, err := GetBilibiliUserInfo(uid)
	if err != nil {
		return errors.Wrapf(err, "failed to get Bilibili user info for %s", uid)
	}

	// 获取最新视频
	videos, err := GetBilibiliUserVideos(uid, 10)
	if err != nil {
		slog.Warn("Failed to get Bilibili videos", "uid", uid, "error", err)
		videos = []BilibiliVideoInfo{} // 继续处理，但没有视频信息
	}

	// 获取视频详细信息（包括点赞数、投币数等）
	for i := range videos {
		if videoDetails, err := GetBilibiliVideoDetails(videos[i].BvID); err == nil {
			videos[i].LikeCount = videoDetails.LikeCount
			videos[i].CoinCount = videoDetails.CoinCount
			videos[i].FavoriteCount = videoDetails.FavoriteCount
			videos[i].ShareCount = videoDetails.ShareCount
			videos[i].ReplyCount = videoDetails.ReplyCount
			if videos[i].Author == "" {
				videos[i].Author = videoDetails.Author
			}
		}
		// 添加延迟避免被限制
		time.Sleep(100 * time.Millisecond)
	}

	// 格式化消息
	messages := FormatBilibiliUserMessage(userInfo, videos)

	// 发送到Discord
	return SendMessageToDiscord(messages, config.GetDiscordChatWebhookUrl(), ServiceNameBilibiliUser)
}
