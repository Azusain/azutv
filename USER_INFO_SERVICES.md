# YouTube & Bilibili 用户信息爬虫服务

本项目新增了两个用户信息爬虫服务，用于获取 YouTube 和 Bilibili 用户的基本信息和视频数据。

## 功能特性

### YouTube 用户信息服务

- **用户基本信息**
  - 频道名称
  - 订阅数 (Subscriber Count)
  - 视频总数
  - 频道简介
  - 头像 URL
  - 频道链接

- **最新视频信息**（默认获取最新 10 个视频）
  - 视频标题
  - 播放量 (View Count)
  - 点赞数 (Like Count)
  - 发布时间
  - 视频时长
  - 视频链接和缩略图

### Bilibili 用户信息服务

- **用户基本信息**
  - 用户名
  - UID
  - 粉丝数 (关注者数)
  - 关注数 (Following)
  - 获赞数 (Total Likes)
  - 总播放量
  - 用户等级 (Lv.1-6)
  - 大会员类型 (普通/月度/年度)
  - 个人简介

- **最新视频信息**（默认获取最新 10 个视频）
  - 视频标题
  - 播放量
  - 点赞数
  - 投币数 (Coin Count)
  - 收藏数 (Favorite Count)
  - 发布时间
  - 视频时长
  - BV号和链接

## 使用方法

### 1. 程序化调用

```go
package main

import (
    "azuserver/service"
    "log"
)

func main() {
    // 获取 YouTube 用户信息
    params := map[string]string{
        "userID": "@MrBeast", // 支持 @username、UCxxxx 格式、或传统用户名
    }
    err := service.RunServiceWithParams("youtube_user", params)
    if err != nil {
        log.Fatal(err)
    }

    // 获取 Bilibili 用户信息
    params = map[string]string{
        "uid": "946974", // Bilibili 用户 UID
    }
    err = service.RunServiceWithParams("bilibili_user", params)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 2. 命令行测试

提供了示例程序用于测试：

```bash
# YouTube 用户信息
go run examples/user_info_example.go youtube @MrBeast
go run examples/user_info_example.go youtube UCX6OQ3DkcsbYNE6H8uQQuVA

# Bilibili 用户信息  
go run examples/user_info_example.go bilibili 1
go run examples/user_info_example.go bilibili 946974
```

## 支持的 ID 格式

### YouTube
- **新格式用户名**: `@username` (如 `@MrBeast`)
- **频道 ID**: `UCxxxxxxxxxxxxxxxxxxxx` (24位以UC开头)
- **传统用户名**: 直接用户名 (如 `pewdiepie`)

### Bilibili
- **用户 UID**: 数字 ID (如 `1`, `946974`)

## 技术实现

### 爬虫技术
- 使用 **Colly** 框架进行网页爬虫
- 采用正则表达式提取 JSON 数据中的用户信息
- 设置合理的 User-Agent 模拟真实浏览器访问
- 添加请求延迟避免被反爬虫机制限制

### 数据提取方式
- **YouTube**: 从页面的 `ytInitialData` 和 meta 标签中提取信息
- **Bilibili**: 从页面的 `__INITIAL_STATE__` 和 HTML 元素中提取信息

### 错误处理
- 网络请求失败时提供详细错误信息
- 数据解析失败时不中断程序执行，记录警告日志
- 支持部分数据获取失败的情况下继续执行

## 输出格式

输出采用 Markdown 格式，适用于 Discord 等支持 Markdown 的平台：

### YouTube 示例输出
```markdown
# YouTube 用户信息
**频道名称**: MrBeast
**用户ID**: @MrBeast
**订阅数**: 329M subscribers
**视频总数**: 700 videos
**频道链接**: https://www.youtube.com/@MrBeast

## 最新视频
### 1. [I Spent $1,000,000 On This Video](https://www.youtube.com/watch?v=abc123)
**观看次数**: 50M views
**点赞数**: 2.1M likes
**发布时间**: 2 days ago
**时长**: 15:30
```

### Bilibili 示例输出
```markdown
# Bilibili 用户信息
**用户名**: 某用户
**UID**: 946974
**粉丝数**: 100.5万
**关注数**: 150
**获赞数**: 500.2万
**播放数**: 1.2亿
**等级**: Lv.6
**个人空间**: https://space.bilibili.com/946974

## 最新视频
### 1. [视频标题](https://www.bilibili.com/video/BV1xx411xxx)
**播放量**: 50.2万
**点赞数**: 1.5万
**投币数**: 5000
**收藏数**: 8000
**发布时间**: 2024-01-15 14:30:00
**时长**: 10:25
```

## Discord 集成

服务支持发送到 Discord webhook：

1. 在 `config` 包中配置 Discord webhook URL
2. 取消注释相关代码行：
```go
return SendMessageToDiscord(messages, config.GetDiscordChatWebhookUrl(), ServiceNameYouTubeUser)
```

## 注意事项

1. **请求频率**: 代码中已添加适当的延迟 (100ms)，避免触发反爬虫机制
2. **网站结构变化**: 如果 YouTube 或 Bilibili 更改页面结构，可能需要更新正则表达式
3. **访问限制**: 某些地区或网络环境可能无法访问 YouTube
4. **用户隐私**: 只获取公开可见的用户信息，不涉及私人数据

## 扩展建议

1. 可以添加缓存机制减少重复请求
2. 支持批量获取多个用户信息
3. 添加数据持久化存储功能
4. 实现定时监控用户数据变化
5. 添加更多平台支持（如 Twitch、TikTok 等）
