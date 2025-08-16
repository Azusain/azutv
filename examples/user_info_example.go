package main

import (
	"azuserver/service"
	"fmt"
	"log"
	"log/slog"
	"os"
)

func init() {
	// 设置日志格式
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func main() {
	fmt.Println("=== AzuTV User Info Services Demo ===")
	
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  go run examples/user_info_example.go youtube <user_id>")
		fmt.Println("  go run examples/user_info_example.go bilibili <uid>")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run examples/user_info_example.go youtube @MrBeast")
		fmt.Println("  go run examples/user_info_example.go youtube UCX6OQ3DkcsbYNE6H8uQQuVA")
		fmt.Println("  go run examples/user_info_example.go bilibili 1")
		fmt.Println("  go run examples/user_info_example.go bilibili 946974")
		return
	}
	
	platform := os.Args[1]
	userID := os.Args[2]
	
	switch platform {
	case "youtube", "yt":
		fmt.Printf("Fetching YouTube user info for: %s\n\n", userID)
		
		params := map[string]string{
			"userID": userID,
		}
		
		if err := service.RunServiceWithParams("youtube_user", params); err != nil {
			log.Fatalf("Failed to get YouTube user info: %v", err)
		}
		
	case "bilibili", "bili":
		fmt.Printf("Fetching Bilibili user info for UID: %s\n\n", userID)
		
		params := map[string]string{
			"uid": userID,
		}
		
		if err := service.RunServiceWithParams("bilibili_user", params); err != nil {
			log.Fatalf("Failed to get Bilibili user info: %v", err)
		}
		
	default:
		fmt.Printf("Unsupported platform: %s\n", platform)
		fmt.Println("Supported platforms: youtube, bilibili")
		return
	}
	
	fmt.Println("\n=== Demo completed ===")
}
