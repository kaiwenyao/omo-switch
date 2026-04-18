package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/kaiwenyao/omo-switch/config"
)

// ApplyChange updates the model value for an agent or category node.
// It updates both the primary model and all fallback_models.
func ApplyChange(cfg *config.Config, nodeID string, newModel string) {
	parts := strings.SplitN(nodeID, "/", 2)
	if len(parts) != 2 {
		log.Printf("Invalid nodeID format: %s", nodeID)
		return
	}

	nodeType := parts[0]
	nodeName := parts[1]

	switch nodeType {
	case "agent":
		if agent, ok := cfg.Agents[nodeName]; ok {
			agent.Model = newModel
			// Update all fallback_models
			for i := range agent.FallbackModels {
				agent.FallbackModels[i].Model = newModel
			}
			cfg.Agents[nodeName] = agent
			fmt.Printf("Updated agent %s to model %s\n", nodeName, newModel)
		}
	case "category":
		if category, ok := cfg.Categories[nodeName]; ok {
			category.Model = newModel
			// Update all fallback_models
			for i := range category.FallbackModels {
				category.FallbackModels[i].Model = newModel
			}
			cfg.Categories[nodeName] = category
			fmt.Printf("Updated category %s to model %s\n", nodeName, newModel)
		}
	default:
		log.Printf("Unknown node type: %s", nodeType)
	}
}

// SaveConfig saves the config to the specified path.
// It logs errors but doesn't crash the app.
func SaveConfig(cfg *config.Config, path string) error {
	if err := cfg.Save(path); err != nil {
		log.Printf("Error saving config to %s: %v", path, err)
		return fmt.Errorf("failed to save config: %w", err)
	}
	fmt.Printf("Config saved to %s\n", path)
	return nil
}
