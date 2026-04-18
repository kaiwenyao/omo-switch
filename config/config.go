package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FallbackModel struct {
	Model string `json:"model"`
}

type Agent struct {
	Model          string          `json:"model"`
	FallbackModels []FallbackModel `json:"fallback_models"`
}

type Category struct {
	Model          string          `json:"model"`
	Variant        string          `json:"variant,omitempty"`
	FallbackModels []FallbackModel `json:"fallback_models"`
}

type Config struct {
	Agents     map[string]Agent     `json:"agents"`
	Categories map[string]Category `json:"categories"`
	GoogleAuth bool                 `json:"google_auth"`
}

func Load(path string) (*Config, error) {
	expandedPath := expandPath(path)

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	expandedPath := expandPath(path)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(expandedPath, data, 0644)
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
