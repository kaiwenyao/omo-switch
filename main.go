package main

import (
	"flag"
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

	model := tui.NewAppModel(cfg)
	model.ConfigPath = *configPath

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Failed to start TUI: %v", err)
	}
}
