package main

import (
	"azuserver/config"
	"azuserver/service"
	"flag"
	"fmt"
	"log/slog"
)

// TODO: test bin options.
func main() {
	// TODO: Load plugin configuration only if the plugin is enabled.
	if err := config.LoadConfig(); err != nil {
		slog.Error(err.Error())
		return
	}

	task := flag.String("task", "", "oricon_ranking, github_trending, vocaloid_ranking, youtube_user, bilibili_user")
	userID := flag.String("user-id", "", "User ID for YouTube (@username, UCxxxx, or username) or Bilibili (numeric UID)")
	uid := flag.String("uid", "", "Bilibili user UID (alias for user-id)")
	flag.Parse()

	// Handle parameterized services
	if *task == "youtube_user" || *task == "bilibili_user" {
		params := make(map[string]string)
		
		if *task == "youtube_user" {
			if *userID == "" {
				slog.Error("YouTube user service requires -user-id parameter")
				slog.Info("Examples: -user-id=@MrBeast, -user-id=UCX6OQ3DkcsbYNE6H8uQQuVA, -user-id=pewdiepie")
				return
			}
			params["userID"] = *userID
		} else if *task == "bilibili_user" {
			// Use uid parameter if provided, otherwise fall back to user-id
			biliUID := *uid
			if biliUID == "" {
				biliUID = *userID
			}
			if biliUID == "" {
				slog.Error("Bilibili user service requires -uid or -user-id parameter")
				slog.Info("Examples: -uid=1, -uid=946974, -user-id=123456")
				return
			}
			params["uid"] = biliUID
		}
		
		if err := service.RunServiceWithParams(*task, params); err != nil {
			slog.Error(fmt.Sprintf("Failed to run parameterized service %s: %v", *task, err))
			return
		}
	} else {
		// Handle standard services
		service.RunService(task)
	}

	slog.Info(fmt.Sprintf("task %s completed", *task))
}
