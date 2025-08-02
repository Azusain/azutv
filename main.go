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

	task := flag.String("task", "", "oricon_ranking, github_trending, vocaloid_ranking")
	flag.Parse()

	service.RunService(task)

	slog.Info(fmt.Sprintf("task %s completed", *task))
}
