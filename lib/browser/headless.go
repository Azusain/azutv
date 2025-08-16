package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/chromedp/chromedp"
)

// HeadlessBrowser 无头浏览器封装
type HeadlessBrowser struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// PageResult 页面结果
type PageResult struct {
	HTML      string                 `json:"html"`
	Title     string                 `json:"title"`
	URL       string                 `json:"url"`
	JSResults map[string]interface{} `json:"js_results"`
	Error     error                  `json:"error"`
}

// NewHeadlessBrowser 创建新的无头浏览器实例
func NewHeadlessBrowser() *HeadlessBrowser {
	// 创建 Chrome 上下文
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-images", true), // 加速加载
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocCtx)

	return &HeadlessBrowser{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Close 关闭浏览器
func (hb *HeadlessBrowser) Close() {
	if hb.cancel != nil {
		hb.cancel()
	}
}

// LoadPage 加载页面并等待JavaScript执行
func (hb *HeadlessBrowser) LoadPage(url string, waitTime time.Duration) (*PageResult, error) {
	result := &PageResult{
		URL:       url,
		JSResults: make(map[string]interface{}),
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(hb.ctx, 30*time.Second)
	defer cancel()

	var html, title string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery), // 等待body元素可见
		chromedp.Sleep(waitTime),                       // 等待JavaScript执行
		chromedp.Title(&title),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		result.Error = err
		return result, err
	}

	result.HTML = html
	result.Title = title
	return result, nil
}

// ExecuteScript 在页面上执行JavaScript脚本
func (hb *HeadlessBrowser) ExecuteScript(url string, script string, waitTime time.Duration) (*PageResult, error) {
	result := &PageResult{
		URL:       url,
		JSResults: make(map[string]interface{}),
	}

	ctx, cancel := context.WithTimeout(hb.ctx, 30*time.Second)
	defer cancel()

	var html, title string
	var scriptResult interface{}

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(waitTime), // 等待页面加载完成
		chromedp.Evaluate(script, &scriptResult),
		chromedp.Title(&title),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		result.Error = err
		return result, err
	}

	result.HTML = html
	result.Title = title
	result.JSResults["main"] = scriptResult
	return result, nil
}

// WaitForElementAndExecute 等待特定元素出现后执行脚本
func (hb *HeadlessBrowser) WaitForElementAndExecute(url string, selector string, script string, maxWait time.Duration) (*PageResult, error) {
	result := &PageResult{
		URL:       url,
		JSResults: make(map[string]interface{}),
	}

	ctx, cancel := context.WithTimeout(hb.ctx, maxWait+10*time.Second)
	defer cancel()

	var html, title string
	var scriptResult interface{}

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // 额外等待，确保数据加载完成
		chromedp.Evaluate(script, &scriptResult),
		chromedp.Title(&title),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		result.Error = err
		slog.Warn("Failed to wait for element and execute script",
			"url", url,
			"selector", selector,
			"error", err)
		return result, err
	}

	result.HTML = html
	result.Title = title
	result.JSResults["main"] = scriptResult
	return result, nil
}

// ExtractBilibiliStats 专门用于提取Bilibili统计数据的方法
func (hb *HeadlessBrowser) ExtractBilibiliStats(uid string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://space.bilibili.com/%s", uid)
	
	// JavaScript 脚本来提取统计数据
	script := `
(function() {
    var stats = {
        username: '',
        follower: 0,
        following: 0,
        like: 0,
        play: 0,
        level: 0,
        vipType: 0,
        description: '',
        face: ''
    };
    
    try {
        // 尝试从页面标题获取用户名
        var title = document.title;
        if (title.includes('的个人空间')) {
            stats.username = title.split('的个人空间')[0];
        }
        
        // 尝试从_render_data_全局变量获取数据
        if (window._render_data_) {
            console.log('Found _render_data_');
        }
        
        // 尝试从__INITIAL_STATE__获取数据
        if (window.__INITIAL_STATE__) {
            var state = window.__INITIAL_STATE__;
            console.log('Found __INITIAL_STATE__');
            
            // 尝试不同的路径来获取用户信息
            if (state.userInfo) {
                stats.username = state.userInfo.name || stats.username;
                stats.description = state.userInfo.sign || '';
                stats.face = state.userInfo.face || '';
                stats.level = state.userInfo.level || 0;
                stats.vipType = state.userInfo.vipType || 0;
            }
            
            // 尝试获取统计信息
            if (state.userInfo && state.userInfo.stats) {
                stats.follower = state.userInfo.stats.follower || 0;
                stats.following = state.userInfo.stats.following || 0;
            }
        }
        
        // 尝试从页面DOM元素获取统计数据
        var numElements = document.querySelectorAll('.n-num');
        numElements.forEach(function(el) {
            var text = el.textContent.trim();
            var parent = el.parentElement;
            if (parent) {
                var parentText = parent.textContent.toLowerCase();
                if (parentText.includes('粉丝') && text.match(/\d/)) {
                    stats.follower = parseCount(text);
                } else if (parentText.includes('关注') && text.match(/\d/)) {
                    stats.following = parseCount(text);
                } else if (parentText.includes('获赞') && text.match(/\d/)) {
                    stats.like = parseCount(text);
                } else if (parentText.includes('播放') && text.match(/\d/)) {
                    stats.play = parseCount(text);
                }
            }
        });
        
        // 尝试其他可能的选择器
        var statSelectors = [
            '[class*="fans"]',
            '[class*="follow"]', 
            '[class*="stat"]',
            '[data-v-*][class*="num"]'
        ];
        
        statSelectors.forEach(function(selector) {
            var elements = document.querySelectorAll(selector);
            elements.forEach(function(el) {
                var text = el.textContent.trim();
                if (text.match(/\d/) && text.length < 20) {
                    console.log('Found potential stat:', selector, text);
                }
            });
        });
        
        // 解析计数的辅助函数
        function parseCount(countStr) {
            if (!countStr) return 0;
            countStr = countStr.trim();
            
            if (countStr.includes('万')) {
                var num = parseFloat(countStr.replace('万', ''));
                return Math.floor(num * 10000);
            } else if (countStr.includes('亿')) {
                var num = parseFloat(countStr.replace('亿', ''));
                return Math.floor(num * 100000000);
            } else {
                var num = parseInt(countStr.replace(/[^\d]/g, ''));
                return isNaN(num) ? 0 : num;
            }
        }
        
        // 尝试等待异步数据加载
        setTimeout(function() {
            // 再次尝试获取数据
        }, 2000);
        
    } catch (e) {
        console.error('Error extracting stats:', e);
        stats.error = e.toString();
    }
    
    return stats;
})();
`

	result, err := hb.WaitForElementAndExecute(url, "body", script, 10*time.Second)
	if err != nil {
		return nil, err
	}

	if result.JSResults["main"] != nil {
		return result.JSResults["main"].(map[string]interface{}), nil
	}

	return make(map[string]interface{}), nil
}

// 辅助方法：将interface{}转换为具体类型
func GetIntFromResult(result map[string]interface{}, key string) int64 {
	if val, ok := result[key]; ok {
		switch v := val.(type) {
		case float64:
			return int64(v)
		case int:
			return int64(v)
		case int64:
			return v
		}
	}
	return 0
}

func GetStringFromResult(result map[string]interface{}, key string) string {
	if val, ok := result[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
