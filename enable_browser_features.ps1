# Bilibili 爬虫无头浏览器功能启用脚本
# Enable Browser Features for Bilibili Scraper

Write-Host "=== Bilibili 爬虫无头浏览器功能启用 ===" -ForegroundColor Green
Write-Host ""

# 检查当前状态
Write-Host "检查当前状态..." -ForegroundColor Yellow

# 1. 检查 Chrome 安装
$chromeFound = $false
$chromePaths = @(
    "$env:ProgramFiles\Google\Chrome\Application\chrome.exe",
    "$env:ProgramFiles(x86)\Google\Chrome\Application\chrome.exe",
    "$env:LOCALAPPDATA\Google\Chrome\Application\chrome.exe"
)

foreach ($path in $chromePaths) {
    if (Test-Path $path) {
        Write-Host "✓ 找到 Chrome 浏览器: $path" -ForegroundColor Green
        $chromeFound = $true
        break
    }
}

if (-not $chromeFound) {
    Write-Host "❌ 未找到 Chrome 浏览器" -ForegroundColor Red
    Write-Host "请先安装 Chrome: https://www.google.com/chrome/" -ForegroundColor Yellow
    Write-Host ""
}

# 2. 安装 Go 依赖
Write-Host "安装 chromedp 依赖..." -ForegroundColor Yellow
try {
    # 设置代理以解决网络问题
    go env -w GOPROXY=https://goproxy.cn,direct
    go get github.com/chromedp/chromedp
    Write-Host "✓ chromedp 依赖安装成功" -ForegroundColor Green
} catch {
    Write-Host "❌ chromedp 依赖安装失败: $_" -ForegroundColor Red
    Write-Host "请手动执行: go get github.com/chromedp/chromedp" -ForegroundColor Yellow
}

# 3. 修改源码以启用浏览器功能
Write-Host "启用浏览器功能..." -ForegroundColor Yellow

$bilibiliFile = "service/bilibili_user.go"
if (Test-Path $bilibiliFile) {
    $content = Get-Content $bilibiliFile -Raw
    
    # 检查是否需要修改
    if ($content -match '/\*[\s\S]*?\*/') {
        Write-Host "发现注释的浏览器代码，正在启用..." -ForegroundColor Yellow
        
        # 添加导入
        if ($content -notmatch '"azuserver/lib/browser"') {
            $content = $content -replace '(import \(\s*"azuserver/config"\s*"encoding/json")', '$1`n`t"azuserver/lib/browser"'
        }
        
        # 取消注释浏览器代码块
        $content = $content -replace '/\*\s*(.*?)\s*\*/', '$1'
        
        # 移除提示信息
        $content = $content -replace 'slog\.Warn\("Browser method is not enabled yet.*?\n.*?return errors\.New\("browser method not available"\)', ''
        
        # 保存修改
        $content | Set-Content $bilibiliFile
        Write-Host "✓ 浏览器功能已启用" -ForegroundColor Green
    } else {
        Write-Host "✓ 浏览器功能已经启用或无需修改" -ForegroundColor Green
    }
} else {
    Write-Host "❌ 找不到 $bilibiliFile 文件" -ForegroundColor Red
}

# 4. 编译测试
Write-Host "编译项目..." -ForegroundColor Yellow
try {
    go build
    Write-Host "✓ 项目编译成功" -ForegroundColor Green
} catch {
    Write-Host "❌ 项目编译失败: $_" -ForegroundColor Red
}

# 5. 运行测试
Write-Host ""
Write-Host "=== 运行测试 ===" -ForegroundColor Green
Write-Host "测试命令: go run main.go -task bilibili_user -uid 36615703"
Write-Host ""

# 提供手动步骤
Write-Host "如果自动启用失败，请手动执行以下步骤：" -ForegroundColor Yellow
Write-Host "1. 在 service/bilibili_user.go 中添加导入: import \"azuserver/lib/browser\""
Write-Host "2. 在 getBilibiliUserStatsWithBrowser 函数中取消注释 /* */ 代码块"
Write-Host "3. 移除或注释掉警告信息和错误返回"
Write-Host ""
Write-Host "详细文档请查看: docs/browser_setup.md" -ForegroundColor Cyan
