package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	uid := "36615703"
	
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// 添加延迟和调试
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("访问: %s\n", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("响应状态: %d\n", r.StatusCode)
		fmt.Printf("响应长度: %d 字节\n", len(r.Body))
		
		// 检查页面是否包含用户信息
		body := string(r.Body)
		if strings.Contains(body, "TimeTom790") {
			fmt.Println("✅ 找到用户名 TimeTom790")
		}
		if strings.Contains(body, "__INITIAL_STATE__") {
			fmt.Println("✅ 找到 __INITIAL_STATE__")
		}
		if strings.Contains(body, "window.__INITIAL_STATE__") {
			fmt.Println("✅ 找到 window.__INITIAL_STATE__")
		}
		if strings.Contains(body, "_render_data_") {
			fmt.Println("✅ 找到 _render_data_")
		}
	})

	// 尝试从页面标题获取用户名
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text
		fmt.Printf("页面标题: %s\n", title)
		
		// 从标题提取用户名 "TimeTom790的个人空间-TimeTom790个人主页-哔哩哔哩视频"
		if strings.Contains(title, "的个人空间") {
			parts := strings.Split(title, "的个人空间")
			if len(parts) > 0 {
				username := parts[0]
				fmt.Printf("从标题提取的用户名: %s\n", username)
			}
		}
	})

	// 尝试从meta标签获取描述
	c.OnHTML("meta[name='description']", func(e *colly.HTMLElement) {
		desc := e.Attr("content")
		fmt.Printf("页面描述: %s\n", desc)
	})

	// 尝试从script标签获取数据
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text
		
		// 检查各种可能的数据格式
		if strings.Contains(scriptContent, "__INITIAL_STATE__") {
			fmt.Println("在script中找到 __INITIAL_STATE__")
			
			// 尝试提取用户名
			nameRegex := regexp.MustCompile(`"name":"([^"]+)"`)
			matches := nameRegex.FindAllStringSubmatch(scriptContent, -1)
			for _, match := range matches {
				if len(match) > 1 {
					fmt.Printf("找到name字段: %s\n", match[1])
				}
			}
			
			// 尝试提取等级
			levelRegex := regexp.MustCompile(`"level":(\d+)`)
			levelMatches := levelRegex.FindAllStringSubmatch(scriptContent, -1)
			for _, match := range levelMatches {
				if len(match) > 1 {
					fmt.Printf("找到level字段: %s\n", match[1])
				}
			}
		}
		
		if strings.Contains(scriptContent, "_render_data_") {
			fmt.Println("在script中找到 _render_data_")
		}
		
		if strings.Contains(scriptContent, "TimeTom790") {
			fmt.Println("在script中找到用户名")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("错误: %v\n", err)
	})

	spaceURL := fmt.Sprintf("https://space.bilibili.com/%s", uid)
	err := c.Visit(spaceURL)
	if err != nil {
		log.Fatal(err)
	}
}
