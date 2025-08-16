# GitHub Actions 使用指南

本项目包含多个 GitHub Actions workflow，用于自动化运行和测试各种服务。

## 📋 可用的 Workflows

### 1. 定时任务 (Scheduled Tasks) - `main.yml`
- **触发方式**: 每天自动运行 + 手动触发
- **功能**: 运行标准服务（Oricon排行、GitHub趋势、Vocaloid排行）
- **用途**: 日常自动化数据收集

### 2. 手动服务测试 (Manual Service Testing) - `manual-test.yml`
- **触发方式**: 仅手动触发
- **功能**: 测试单个服务，支持所有服务类型
- **用途**: 开发和调试时的精确测试

### 3. 快速预设测试 (Quick Service Tests) - `quick-test.yml`  
- **触发方式**: 仅手动触发
- **功能**: 使用预设参数快速测试常见场景
- **用途**: 快速验证服务功能

## 🚀 使用方法

### 在 GitHub 网页上手动运行

1. **访问 Actions 页面**
   - 进入你的 GitHub 仓库
   - 点击 "Actions" 标签页
   
2. **选择 Workflow**
   - 在左侧列表中选择要运行的 workflow
   - 点击 "Run workflow" 按钮

3. **填写参数**（根据选择的 workflow）

### Manual Service Testing 使用方法

这是最灵活的测试方式，支持所有服务：

#### 参数说明
- **Service**: 从下拉菜单选择服务
  - `oricon_ranking`: Oricon音乐排行榜
  - `github_trending`: GitHub趋势项目
  - `vocaloid_ranking`: Vocaloid排行榜  
  - `youtube_user`: YouTube用户信息
  - `bilibili_user`: Bilibili用户信息

- **User ID**: 用户标识符（仅YouTube/Bilibili服务需要）
  - YouTube格式: `@MrBeast`, `UCX6OQ3DkcsbYNE6H8uQQuVA`, `pewdiepie`
  - Bilibili格式: `1`, `946974`, `123456`

- **Additional params**: 额外参数（JSON格式，可选）

#### 使用示例

**测试 YouTube 服务：**
```
Service: youtube_user
User ID: @MrBeast
```

**测试 Bilibili 服务：**
```
Service: bilibili_user  
User ID: 946974
```

**测试标准服务：**
```
Service: github_trending
User ID: (留空)
```

### Quick Service Tests 使用方法

提供预设的测试组合，无需手动输入复杂参数：

#### 可用预设
- **all_standard**: 运行所有标准服务
- **youtube_mrbeast**: 测试 MrBeast 频道
- **youtube_pewdiepie**: 测试 PewDiePie 频道
- **bilibili_admin**: 测试 Bilibili 管理员账户
- **bilibili_popular_up**: 测试热门UP主
- **custom_youtube**: 自定义 YouTube 测试
- **custom_bilibili**: 自定义 Bilibili 测试

#### 使用示例

**快速测试所有标准服务：**
```
Test Preset: all_standard
```

**测试特定 YouTube 频道：**
```  
Test Preset: youtube_mrbeast
```

**自定义测试：**
```
Test Preset: custom_youtube
Custom User ID: @your_favorite_channel
```

## 🔧 本地命令行使用

编译并运行项目：

```bash
# 编译
go build -o main

# 运行标准服务
./main -task=github_trending
./main -task=oricon_ranking
./main -task=vocaloid_ranking

# 运行 YouTube 用户服务
./main -task=youtube_user -user-id=@MrBeast
./main -task=youtube_user -user-id=UCX6OQ3DkcsbYNE6H8uQQuVA

# 运行 Bilibili 用户服务  
./main -task=bilibili_user -uid=1
./main -task=bilibili_user -user-id=946974

# 查看帮助
./main -h
```

## 📊 输出和结果

### 日志输出
所有 workflow 都会在 GitHub Actions 的日志中显示详细输出：
- 服务运行状态
- 获取的数据内容
- 错误信息（如有）
- 执行时间和统计

### Discord 集成
如果配置了 Discord webhook secrets，结果也会发送到指定的 Discord 频道：
- `DISCORD_CHAT_WEBHOOK_URL`: 聊天频道webhook
- `DISCORD_SYS_WEBHOOK_URL`: 系统通知webhook

### 数据格式
输出采用 Markdown 格式，包含：
- 用户基本信息（订阅数、粉丝数等）
- 最新视频信息（标题、播放量、发布时间等）
- 统计数据和链接

## ⚙️ 环境变量和 Secrets

需要在 GitHub 仓库设置中配置以下 secrets：

```
DISCORD_CHAT_WEBHOOK_URL=https://discord.com/api/webhooks/your-webhook-url
DISCORD_SYS_WEBHOOK_URL=https://discord.com/api/webhooks/your-system-webhook-url  
```

### 设置步骤
1. 进入 GitHub 仓库
2. Settings → Secrets and variables → Actions
3. 点击 "New repository secret"
4. 添加上述变量

## 🛠️ 故障排除

### 常见问题

**1. YouTube/Bilibili 服务失败**
- 检查用户ID格式是否正确
- 确认用户/频道是否存在且公开
- 查看网络连接是否正常

**2. Discord 消息未发送**
- 检查 webhook secrets 是否正确配置
- 确认 Discord webhook URL 有效

**3. 编译错误**
- 检查 Go 版本是否为 1.24.1
- 运行 `go mod tidy` 更新依赖

**4. 权限错误**
- 确保仓库有足够的 Actions 权限
- 检查分支保护规则

### 调试技巧

1. **查看详细日志**
   - 在 Actions 页面点击失败的运行
   - 展开各个步骤查看详细输出

2. **本地测试**
   - 先在本地环境测试命令
   - 确认服务正常后再在 Actions 中运行

3. **逐步测试**
   - 从简单的标准服务开始测试
   - 逐步添加复杂的参数化服务

## 📝 扩展和自定义

### 添加新的预设
在 `quick-test.yml` 中添加新的测试预设：

```yaml
- name: Test New Preset
  if: ${{ inputs.test_preset == 'your_new_preset' }}
  run: |
    echo "Testing new preset..."
    ./main -task=your_service -param=your_value
```

### 修改运行频率
在 `main.yml` 中修改 cron 表达式：

```yaml
schedule:
  - cron: "0 */6 * * *"  # 每6小时运行一次
```

### 添加新的服务参数
在 `manual-test.yml` 中添加新的输入参数：

```yaml
inputs:
  your_new_param:
    description: 'Your parameter description'
    required: false
    type: string
```

这样你就可以通过 GitHub 网页界面轻松测试所有服务了！
