package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kaiwenyao/omo-switch/config"
)

func main() {
	configPath := flag.String("config", "~/.config/opencode/oh-my-openagent.json", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded config: %d agents, %d categories, google_auth=%v\n",
		len(cfg.Agents), len(cfg.Categories), cfg.GoogleAuth)
}
