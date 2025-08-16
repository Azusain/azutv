# 无头浏览器功能设置指南

本项目包含一个可选的无头浏览器功能，用于获取需要 JavaScript 执行的动态内容，特别是 Bilibili 用户的统计数据（粉丝数、关注数等）。

## 功能说明

无头浏览器库 (`lib/browser`) 提供以下功能：
- 启动无头 Chrome 浏览器
- 执行 JavaScript 代码
- 等待动态内容加载
- 专门针对 Bilibili 统计数据的提取

## 安装步骤

### 1. 安装依赖

```bash
# 安装 chromedp 依赖
go get github.com/chromedp/chromedp

# 如果网络连接有问题，尝试使用代理或直接下载
go env -w GOPROXY=https://goproxy.cn,direct
go get github.com/chromedp/chromedp
```

### 2. 安装 Chrome/Chromium

确保系统中安装了 Chrome 或 Chromium 浏览器：

**Windows:**
- 下载并安装 Chrome 浏览器：https://www.google.com/chrome/
- 或下载 Chromium：https://chromium.org/getting-involved/download-chromium

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y chromium-browser

# 或安装 Chrome
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
sudo sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'
sudo apt-get update
sudo apt-get install -y google-chrome-stable
```

**macOS:**
```bash
# 使用 Homebrew
brew install --cask google-chrome
# 或
brew install --cask chromium
```

### 3. 启用浏览器功能

一旦依赖安装完成，需要取消注释 `service/bilibili_user.go` 中的浏览器代码：

1. 在 `getBilibiliUserStatsWithBrowser` 函数中，取消注释带有 `/*` 和 `*/` 的代码块
2. 添加必要的导入：

```go
import (
    // ... 其他导入
    "azuserver/lib/browser"
)
```

### 4. 编译和测试

```bash
# 编译项目
go build

# 测试 Bilibili 用户信息获取
go run main.go -task bilibili_user -uid 36615703
```

## 使用说明

### 基本使用

```go
// 创建浏览器实例
browser := browser.NewHeadlessBrowser()
defer browser.Close()

// 提取 Bilibili 统计数据
stats, err := browser.ExtractBilibiliStats("36615703")
if err != nil {
    log.Fatal(err)
}

// 获取特定数据
followerCount := browser.GetIntFromResult(stats, "follower")
username := browser.GetStringFromResult(stats, "username")
```

### 自定义 JavaScript 执行

```go
browser := browser.NewHeadlessBrowser()
defer browser.Close()

// 执行自定义脚本
script := `document.title`
result, err := browser.ExecuteScript("https://example.com", script, 5*time.Second)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Page title:", result.JSResults["main"])
```

## 配置选项

无头浏览器支持以下配置：
- `--headless`: 无头模式（默认启用）
- `--disable-gpu`: 禁用 GPU 加速
- `--no-sandbox`: 禁用沙盒模式
- `--disable-images`: 禁用图片加载以提高速度

## 故障排除

### 常见问题

1. **Chrome 未找到**
   - 确保安装了 Chrome 或 Chromium
   - 检查 PATH 环境变量是否包含浏览器路径

2. **依赖安装失败**
   - 使用代理：`go env -w GOPROXY=https://goproxy.cn,direct`
   - 手动下载依赖包

3. **权限问题**
   - Linux 上可能需要 `--no-sandbox` 选项
   - 确保有足够的磁盘空间和内存

4. **网络连接问题**
   - 检查防火墙设置
   - 确保目标网站可访问

### 调试模式

如需调试浏览器操作，可以启用非无头模式：

```go
opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.Flag("headless", false), // 禁用无头模式以查看浏览器窗口
    // ... 其他选项
)
```

## 性能考虑

- 浏览器启动有一定开销，建议复用实例
- 禁用不必要的功能（如图片、插件）可以提高速度
- 设置合理的超时时间
- 在高并发场景下注意资源限制

## 安全注意事项

- 在生产环境中启用沙盒模式
- 验证执行的 JavaScript 代码
- 限制可访问的域名
- 定期更新 Chrome/Chromium 版本
