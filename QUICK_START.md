# 🚀 AzuTV 快速开始指南

欢迎使用 AzuTV！这是一个多功能的数据抓取和通知服务，支持多个平台的信息获取。

## ✨ 新功能亮点

### 🎬 YouTube 用户信息服务
- 获取频道基本信息（订阅数、视频总数）
- 抓取最新视频数据（播放量、点赞数、发布时间）
- 支持多种 ID 格式：`@MrBeast`、`UCX6OQ3DkcsbYNE6H8uQQuVA`、`pewdiepie`

### 📺 Bilibili 用户信息服务  
- 获取用户统计数据（粉丝数、获赞数、等级）
- 抓取视频详情（播放量、点赞数、投币数、收藏数）
- 支持中文数字格式（万、亿）自动解析

## 🏃‍♂️ 5分钟快速体验

### 1. 本地运行

```bash
# 编译项目
go build -o main

# 测试 YouTube 服务
./main -task=youtube_user -user-id=@MrBeast

# 测试 Bilibili 服务  
./main -task=bilibili_user -uid=946974

# 运行标准服务
./main -task=github_trending
```

### 2. GitHub Actions 网页测试

1. **进入 GitHub 仓库 → Actions 页面**

2. **选择 "Manual Service Testing"**

3. **配置参数：**
   ```
   Service: youtube_user  (从下拉菜单选择)
   User ID: @MrBeast     (输入用户ID)
   ```

4. **点击 "Run workflow" 并查看结果**

## 📊 支持的服务

| 服务类型 | 描述 | 参数要求 | 示例 |
|---------|------|---------|------|
| `oricon_ranking` | 日本音乐排行榜 | 无 | `./main -task=oricon_ranking` |
| `github_trending` | GitHub趋势项目 | 无 | `./main -task=github_trending` |
| `vocaloid_ranking` | Vocaloid音乐排行 | 无 | `./main -task=vocaloid_ranking` |
| `youtube_user` | YouTube用户信息 | `user-id` | `./main -task=youtube_user -user-id=@MrBeast` |
| `bilibili_user` | Bilibili用户信息 | `uid` | `./main -task=bilibili_user -uid=946974` |

## 🔧 配置要求

### 基础配置
```yaml
# config.yaml
discord:
  chat_webhook_url: "your_webhook_url_here"
  sys_webhook_url: "your_system_webhook_url_here"
```

### GitHub Actions Secrets
在仓库设置中添加：
- `DISCORD_CHAT_WEBHOOK_URL`
- `DISCORD_SYS_WEBHOOK_URL`

## 📈 输出示例

### YouTube 用户信息
```markdown
# YouTube 用户信息
**频道名称**: MrBeast
**用户ID**: @MrBeast  
**订阅数**: 329M subscribers
**视频总数**: 700 videos

## 最新视频
### 1. [I Spent $1,000,000 On This Video](https://youtube.com/watch?v=abc)
**观看次数**: 50M views
**点赞数**: 2.1M likes
**发布时间**: 2 days ago
**时长**: 15:30
```

### Bilibili 用户信息
```markdown
# Bilibili 用户信息  
**用户名**: 某UP主
**UID**: 946974
**粉丝数**: 100.5万
**获赞数**: 500.2万
**等级**: Lv.6

## 最新视频
### 1. [视频标题](https://bilibili.com/video/BV1xx411xxx)
**播放量**: 50.2万
**点赞数**: 1.5万
**投币数**: 5000
**收藏数**: 8000
**发布时间**: 2024-01-15 14:30:00
```

## 🔍 常用命令速查

```bash
# 查看帮助
./main -h

# YouTube 服务示例
./main -task=youtube_user -user-id=@MrBeast
./main -task=youtube_user -user-id=UCX6OQ3DkcsbYNE6H8uQQuVA  
./main -task=youtube_user -user-id=pewdiepie

# Bilibili 服务示例
./main -task=bilibili_user -uid=1
./main -task=bilibili_user -uid=946974
./main -task=bilibili_user -user-id=123456

# 标准服务
./main -task=oricon_ranking
./main -task=github_trending
./main -task=vocaloid_ranking
```

## 📚 详细文档

- 📖 **[用户信息服务详解](USER_INFO_SERVICES.md)**: 技术实现和API说明
- 🤖 **[GitHub Actions指南](GITHUB_ACTIONS_GUIDE.md)**: 自动化和CI/CD配置
- 🔧 **[开发者文档](README.md)**: 项目结构和开发指南

## 🆘 需要帮助？

### 常见问题
1. **编译失败**: 检查 Go 版本是否为 1.24.1+
2. **服务无响应**: 确认网络连接和用户ID格式
3. **Discord未收到消息**: 检查webhook配置和secrets

### 调试步骤
1. 先在本地测试命令行版本
2. 检查配置文件和参数格式
3. 查看详细的错误日志
4. 在GitHub Actions中逐步测试

### 获取支持
- 查看 Issues 页面寻找相似问题
- 提交新的 Issue 描述问题
- 参考详细文档进行故障排除

---

🎉 **开始探索 AzuTV 的强大功能吧！**
