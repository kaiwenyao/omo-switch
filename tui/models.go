package tui

import (
	"bufio"
	"context"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kaiwenyao/omo-switch/config"
)

const opencodeModelsTimeout = 3 * time.Second

// Provider represents a model provider with its available models.
type Provider struct {
	Name      string
	Models    []string // Display names (without provider prefix)
	FullPaths []string // Full model paths (for "recently" provider)
}

var (
	cachedAllModels []string
	allModelsMutex  sync.RWMutex
)

func fetchAllModelsFromOpenCode() []string {
	ctx, cancel := context.WithTimeout(context.Background(), opencodeModelsTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "opencode", "models")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var models []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "[") && !strings.HasPrefix(line, "Error:") {
			models = append(models, line)
		}
	}

	return models
}

func GetAllModelsFromOpenCode() []string {
	allModelsMutex.RLock()
	if cachedAllModels != nil {
		result := make([]string, len(cachedAllModels))
		copy(result, cachedAllModels)
		allModelsMutex.RUnlock()
		return result
	}
	allModelsMutex.RUnlock()

	models := fetchAllModelsFromOpenCode()
	if models != nil && len(models) > 0 {
		sort.Strings(models)
		allModelsMutex.Lock()
		cachedAllModels = models
		allModelsMutex.Unlock()
	}

	allModelsMutex.RLock()
	result := make([]string, len(cachedAllModels))
	copy(result, cachedAllModels)
	allModelsMutex.RUnlock()
	return result
}

func GetAvailableModels(cfg *config.Config) []string {
	allModels := GetAllModelsFromOpenCode()
	if len(allModels) > 0 {
		return allModels
	}

	modelSet := make(map[string]bool)

	for _, agent := range cfg.Agents {
		if agent.Model != "" {
			modelSet[agent.Model] = true
		}
		for _, fm := range agent.FallbackModels {
			if fm.Model != "" {
				modelSet[fm.Model] = true
			}
		}
	}

	for _, category := range cfg.Categories {
		if category.Model != "" {
			modelSet[category.Model] = true
		}
		for _, fm := range category.FallbackModels {
			if fm.Model != "" {
				modelSet[fm.Model] = true
			}
		}
	}

	models := make([]string, 0, len(modelSet))
	for model := range modelSet {
		models = append(models, model)
	}
	sort.Strings(models)
	return models
}

// GetProvidersWithModels groups models by provider.
func GetProvidersWithModels(cfg *config.Config) []Provider {
	models := GetAvailableModels(cfg)
	providerMap := make(map[string][]string)

	for _, model := range models {
		provider := extractProviderFromModel(model)
		modelName := strings.TrimPrefix(model, provider+"/")
		providerMap[provider] = append(providerMap[provider], modelName)
	}

	providers := make([]string, 0, len(providerMap))
	for provider := range providerMap {
		providers = append(providers, provider)
	}
	sort.Strings(providers)

	result := make([]Provider, 0, len(providers))
	for _, provider := range providers {
		models := providerMap[provider]
		sort.Strings(models)
		result = append(result, Provider{
			Name:      provider,
			Models:    models,
			FullPaths: nil,
		})
	}

	return result
}

// GetAllProviders returns providers with "recently" pseudo-provider at the top.
func GetAllProviders(cfg *config.Config, nodeID string) []Provider {
	providers := GetProvidersWithModels(cfg)

	parts := strings.SplitN(nodeID, "/", 2)
	if len(parts) != 2 {
		return providers
	}

	nodeType := parts[0]
	nodeName := parts[1]

	var recentModels []string
	switch nodeType {
	case "agent":
		if agent, ok := cfg.Agents[nodeName]; ok {
			recentModels = agent.RecentModels
		}
	case "category":
		if category, ok := cfg.Categories[nodeName]; ok {
			recentModels = category.RecentModels
		}
	}

	if len(recentModels) == 0 {
		return providers
	}

	recentModelNames := make([]string, 0, len(recentModels))
	for _, model := range recentModels {
		if idx := strings.Index(model, "/"); idx != -1 {
			recentModelNames = append(recentModelNames, model[idx+1:])
		} else {
			recentModelNames = append(recentModelNames, model)
		}
	}

	result := make([]Provider, 0, len(providers)+1)
	result = append(result, Provider{
		Name:      "recently",
		Models:    recentModelNames,
		FullPaths: recentModels,
	})
	result = append(result, providers...)

	return result
}

// extractProviderFromModel extracts the provider name from a model string.
func extractProviderFromModel(model string) string {
	if idx := strings.Index(model, "/"); idx != -1 {
		return model[:idx]
	}
	return model
}
