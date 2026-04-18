package tui

import (
	"sort"
	"strings"

	"github.com/kaiwenyao/omo-switch/config"
)

// TreeNode represents a node in the config tree hierarchy.
// Level 0: Provider nodes (e.g., "ollama-cloud", "minimax-cn-coding-plan")
// Level 1: Agent or Category nodes under a provider
type TreeNode struct {
	ID       string      // unique ID for selection (e.g., "agent/sisyphus")
	Type     string      // "provider", "agent", "category"
	Label    string      // display name
	Model    string      // current model value (for agents/categories)
	Depth    int         // 0=provider, 1=agent/category
	Children []*TreeNode // provider's agents/categories
	Parent   *TreeNode
	Selected bool // for batch selection
}

// providerFromModel extracts the provider name from a model string.
// Examples:
//   - "ollama-cloud/glm-5.1" -> "ollama-cloud"
//   - "minimax-cn-coding-plan/....." -> "minimax-cn-coding-plan"
//   - "google/gemini-2.0" -> "google"
func providerFromModel(model string) string {
	if idx := strings.Index(model, "/"); idx != -1 {
		return model[:idx]
	}
	// If no slash, the entire model string is the provider
	return model
}

// BuildTree constructs the tree hierarchy from config.
// Returns root provider nodes with agent/category children.
func BuildTree(cfg *config.Config) []*TreeNode {
	// Group agents and categories by provider
	providerAgents := make(map[string][]*TreeNode)
	providerCategories := make(map[string][]*TreeNode)

	// Process agents
	for name, agent := range cfg.Agents {
		provider := providerFromModel(agent.Model)
		node := &TreeNode{
			ID:    "agent/" + name,
			Type:  "agent",
			Label: name,
			Model: agent.Model,
			Depth: 1,
		}
		providerAgents[provider] = append(providerAgents[provider], node)
	}

	// Process categories
	for name, category := range cfg.Categories {
		provider := providerFromModel(category.Model)
		node := &TreeNode{
			ID:    "category/" + name,
			Type:  "category",
			Label: name,
			Model: category.Model,
			Depth: 1,
		}
		providerCategories[provider] = append(providerCategories[provider], node)
	}

	// Collect all unique providers
	providerSet := make(map[string]bool)
	for provider := range providerAgents {
		providerSet[provider] = true
	}
	for provider := range providerCategories {
		providerSet[provider] = true
	}

	// Sort providers for consistent ordering
	providers := make([]string, 0, len(providerSet))
	for provider := range providerSet {
		providers = append(providers, provider)
	}
	sort.Strings(providers)

	// Build provider nodes with sorted children
	var roots []*TreeNode
	for _, provider := range providers {
		// Sort agents and categories
		agents := providerAgents[provider]
		categories := providerCategories[provider]
		sort.Slice(agents, func(i, j int) bool {
			return agents[i].Label < agents[j].Label
		})
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Label < categories[j].Label
		})

		// Interleave agents and categories under the provider
		children := make([]*TreeNode, 0, len(agents)+len(categories))
		children = append(children, agents...)
		children = append(children, categories...)

		providerNode := &TreeNode{
			ID:       "provider/" + provider,
			Type:     "provider",
			Label:    provider,
			Depth:    0,
			Children: children,
		}

		// Set parent reference for children
		for _, child := range children {
			child.Parent = providerNode
		}

		roots = append(roots, providerNode)
	}

	return roots
}

// Flatten performs depth-first traversal of the tree,
// returning a flat list suitable for linear navigation.
func Flatten(tree []*TreeNode) []TreeNode {
	var result []TreeNode
	var dfs func(nodes []*TreeNode)
	dfs = func(nodes []*TreeNode) {
		for _, node := range nodes {
			// Append a copy of the node (not the pointer)
			result = append(result, *node)
			if len(node.Children) > 0 {
				dfs(node.Children)
			}
		}
	}
	dfs(tree)
	return result
}
