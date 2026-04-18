package tui

import (
	"fmt"
	"strings"

	"github.com/kaiwenyao/omo-switch/config"
)

// ApplyChange updates the model value for an agent or category node.
// When batch is true (multiple nodes selected), it also updates all fallback_models
// to the same model per the spec. When batch is false (single-node edit), only the
// primary model is changed, preserving the existing fallback chain.
func ApplyChange(cfg *config.Config, nodeID string, newModel string, batch bool) {
	parts := strings.SplitN(nodeID, "/", 2)
	if len(parts) != 2 {
		return
	}

	nodeType := parts[0]
	nodeName := parts[1]

	switch nodeType {
	case "agent":
		if agent, ok := cfg.Agents[nodeName]; ok {
			agent.Model = newModel
			if batch {
				for i := range agent.FallbackModels {
					agent.FallbackModels[i].Model = newModel
				}
			}
			cfg.Agents[nodeName] = agent
		}
	case "category":
		if category, ok := cfg.Categories[nodeName]; ok {
			category.Model = newModel
			if batch {
				for i := range category.FallbackModels {
					category.FallbackModels[i].Model = newModel
				}
			}
			cfg.Categories[nodeName] = category
		}
	}
}

// SaveConfig saves the config to the specified path.
func SaveConfig(cfg *config.Config, path string) error {
	if err := cfg.Save(path); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}
