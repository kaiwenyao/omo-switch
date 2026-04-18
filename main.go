package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbletea"
	"github.com/kaiwenyao/omo-switch/config"
	"github.com/kaiwenyao/omo-switch/tui"
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

	// Create the TUI model with the config path for auto-save
	model := tui.NewAppModel(cfg)
	model.ConfigPath = *configPath

	// Run the TUI
	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		log.Fatalf("Failed to start TUI: %v", err)
	}
}
