package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const MaxRecentModels = 10

type FallbackModel struct {
	Model string `json:"model"`
	extra map[string]json.RawMessage
}

func (f FallbackModel) MarshalJSON() ([]byte, error) {
	m := make(map[string]json.RawMessage, len(f.extra)+1)
	for k, v := range f.extra {
		m[k] = v
	}
	data, err := json.Marshal(f.Model)
	if err != nil {
		return nil, err
	}
	m["model"] = data
	return json.Marshal(m)
}

func (f *FallbackModel) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["model"]; ok {
		if err := json.Unmarshal(v, &f.Model); err != nil {
			return err
		}
		delete(raw, "model")
	}
	f.extra = raw
	return nil
}

type Agent struct {
	Model          string
	FallbackModels []FallbackModel
	RecentModels   []string
	extra          map[string]json.RawMessage
}

func (a Agent) MarshalJSON() ([]byte, error) {
	m := make(map[string]json.RawMessage, len(a.extra)+3)
	for k, v := range a.extra {
		m[k] = v
	}
	model, err := json.Marshal(a.Model)
	if err != nil {
		return nil, err
	}
	m["model"] = model
	fbs, err := json.Marshal(a.FallbackModels)
	if err != nil {
		return nil, err
	}
	m["fallback_models"] = fbs
	if len(a.RecentModels) > 0 {
		recent, err := json.Marshal(a.RecentModels)
		if err != nil {
			return nil, err
		}
		m["recent_models"] = recent
	} else {
		delete(m, "recent_models")
	}
	return json.Marshal(m)
}

func (a *Agent) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["model"]; ok {
		if err := json.Unmarshal(v, &a.Model); err != nil {
			return err
		}
		delete(raw, "model")
	}
	if v, ok := raw["fallback_models"]; ok {
		if err := json.Unmarshal(v, &a.FallbackModels); err != nil {
			return err
		}
		delete(raw, "fallback_models")
	}
	if v, ok := raw["recent_models"]; ok {
		if err := json.Unmarshal(v, &a.RecentModels); err != nil {
			return err
		}
		delete(raw, "recent_models")
	}
	a.extra = raw
	return nil
}

type Category struct {
	Model          string
	Variant        string
	FallbackModels []FallbackModel
	RecentModels   []string
	extra          map[string]json.RawMessage
}

func (c Category) MarshalJSON() ([]byte, error) {
	m := make(map[string]json.RawMessage, len(c.extra)+4)
	for k, v := range c.extra {
		m[k] = v
	}
	model, err := json.Marshal(c.Model)
	if err != nil {
		return nil, err
	}
	m["model"] = model
	if c.Variant != "" {
		variant, err := json.Marshal(c.Variant)
		if err != nil {
			return nil, err
		}
		m["variant"] = variant
	} else {
		delete(m, "variant")
	}
	fbs, err := json.Marshal(c.FallbackModels)
	if err != nil {
		return nil, err
	}
	m["fallback_models"] = fbs
	if len(c.RecentModels) > 0 {
		recent, err := json.Marshal(c.RecentModels)
		if err != nil {
			return nil, err
		}
		m["recent_models"] = recent
	} else {
		delete(m, "recent_models")
	}
	return json.Marshal(m)
}

func (c *Category) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["model"]; ok {
		if err := json.Unmarshal(v, &c.Model); err != nil {
			return err
		}
		delete(raw, "model")
	}
	if v, ok := raw["variant"]; ok {
		if err := json.Unmarshal(v, &c.Variant); err != nil {
			return err
		}
		delete(raw, "variant")
	}
	if v, ok := raw["fallback_models"]; ok {
		if err := json.Unmarshal(v, &c.FallbackModels); err != nil {
			return err
		}
		delete(raw, "fallback_models")
	}
	if v, ok := raw["recent_models"]; ok {
		if err := json.Unmarshal(v, &c.RecentModels); err != nil {
			return err
		}
		delete(raw, "recent_models")
	}
	c.extra = raw
	return nil
}

type Config struct {
	Agents     map[string]Agent
	Categories map[string]Category
	GoogleAuth bool
	extra      map[string]json.RawMessage
}

func (c Config) MarshalJSON() ([]byte, error) {
	m := make(map[string]json.RawMessage, len(c.extra)+3)
	for k, v := range c.extra {
		m[k] = v
	}
	agents, err := json.Marshal(c.Agents)
	if err != nil {
		return nil, err
	}
	m["agents"] = agents
	cats, err := json.Marshal(c.Categories)
	if err != nil {
		return nil, err
	}
	m["categories"] = cats
	gauth, err := json.Marshal(c.GoogleAuth)
	if err != nil {
		return nil, err
	}
	m["google_auth"] = gauth
	return json.Marshal(m)
}

func (c *Config) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["agents"]; ok {
		if err := json.Unmarshal(v, &c.Agents); err != nil {
			return fmt.Errorf("parse agents: %w", err)
		}
		delete(raw, "agents")
	}
	if v, ok := raw["categories"]; ok {
		if err := json.Unmarshal(v, &c.Categories); err != nil {
			return fmt.Errorf("parse categories: %w", err)
		}
		delete(raw, "categories")
	}
	if v, ok := raw["google_auth"]; ok {
		if err := json.Unmarshal(v, &c.GoogleAuth); err != nil {
			return fmt.Errorf("parse google_auth: %w", err)
		}
		delete(raw, "google_auth")
	}
	c.extra = raw
	return nil
}

// AddRecentModel adds a model to the recent list for an agent or category.
func (c *Config) AddRecentModel(nodeID string, model string) {
	parts := strings.SplitN(nodeID, "/", 2)
	if len(parts) != 2 {
		return
	}
	nodeType := parts[0]
	nodeName := parts[1]

	switch nodeType {
	case "agent":
		if agent, ok := c.Agents[nodeName]; ok {
			agent.RecentModels = prependUnique(agent.RecentModels, model, MaxRecentModels)
			c.Agents[nodeName] = agent
		}
	case "category":
		if category, ok := c.Categories[nodeName]; ok {
			category.RecentModels = prependUnique(category.RecentModels, model, MaxRecentModels)
			c.Categories[nodeName] = category
		}
	}
}

func prependUnique(list []string, item string, max int) []string {
	out := make([]string, 0, len(list)+1)
	out = append(out, item)
	for _, m := range list {
		if m != item {
			out = append(out, m)
		}
	}
	if len(out) > max {
		out = out[:max]
	}
	return out
}

func Load(path string) (*Config, error) {
	expandedPath := expandPath(path)

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", expandedPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

// Save writes the config atomically: marshal → temp file in same dir → fsync → rename.
// Permissions are tightened to 0600 since the file may contain auth tokens.
func (c *Config) Save(path string) error {
	expandedPath := expandPath(path)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	dir := filepath.Dir(expandedPath)
	tmp, err := os.CreateTemp(dir, ".omo-switch-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	cleanup := func() { _ = os.Remove(tmpPath) }

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("fsync temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, expandedPath); err != nil {
		cleanup()
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}
