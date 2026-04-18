package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// normalizeJSON re-marshals input through a map for stable key-order-agnostic comparison.
func normalizeJSON(t *testing.T, data []byte) string {
	t.Helper()
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("normalize: %v", err)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("normalize marshal: %v", err)
	}
	return string(out)
}

func TestLoadSaveRoundtripPreservesUnknownFields(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "top-level unknown fields",
			json: `{
  "version": 3,
  "telemetry": { "enabled": false, "id": "abc-123" },
  "agents": {
    "foo": {
      "model": "ollama-cloud/gpt-4",
      "fallback_models": [{"model": "anthropic/claude"}]
    }
  },
  "categories": {
    "writing": {
      "model": "google/gemini",
      "variant": "fast",
      "fallback_models": []
    }
  },
  "google_auth": true
}`,
		},
		{
			name: "per-agent unknown fields",
			json: `{
  "agents": {
    "bar": {
      "model": "openai/gpt-5",
      "fallback_models": [{"model": "x/y", "weight": 5, "headers": {"Authorization": "Bearer z"}}],
      "temperature": 0.7,
      "system_prompt": "be helpful"
    }
  },
  "categories": {},
  "google_auth": false
}`,
		},
		{
			name: "per-category unknown fields",
			json: `{
  "agents": {},
  "categories": {
    "code": {
      "model": "anthropic/claude-4",
      "fallback_models": [],
      "max_tokens": 8000
    }
  },
  "google_auth": false
}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "cfg.json")
			if err := os.WriteFile(path, []byte(tc.json), 0600); err != nil {
				t.Fatal(err)
			}
			cfg, err := Load(path)
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			if err := cfg.Save(path); err != nil {
				t.Fatalf("save: %v", err)
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			before := normalizeJSON(t, []byte(tc.json))
			got := normalizeJSON(t, after)
			if before != got {
				t.Errorf("unknown fields lost on roundtrip\n--- before ---\n%s\n--- after ---\n%s", before, got)
			}
		})
	}
}

func TestSaveUsesAtomicRename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	initial := `{"agents":{},"categories":{},"google_auth":false}`
	if err := os.WriteFile(path, []byte(initial), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if err := cfg.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".omo-switch-") {
			t.Errorf("temp file leaked after save: %s", e.Name())
		}
	}
}

func TestSavePermissions0600(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	if err := os.WriteFile(path, []byte(`{"agents":{},"categories":{},"google_auth":false}`), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.Save(path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0600 {
		t.Errorf("want 0600, got %o", mode)
	}
}

func TestAddRecentModelPrependsAndDedupes(t *testing.T) {
	cfg := &Config{
		Agents: map[string]Agent{
			"a1": {Model: "p/m1", FallbackModels: []FallbackModel{}, RecentModels: []string{"p/m2", "p/m3"}},
		},
		Categories: map[string]Category{},
	}
	cfg.AddRecentModel("agent/a1", "p/m3") // duplicate should move to front
	want := []string{"p/m3", "p/m2"}
	if got := cfg.Agents["a1"].RecentModels; !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}

	cfg.AddRecentModel("agent/a1", "p/m4") // new entry prepended
	want = []string{"p/m4", "p/m3", "p/m2"}
	if got := cfg.Agents["a1"].RecentModels; !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestAddRecentModelHonoursMaxSize(t *testing.T) {
	cfg := &Config{
		Agents: map[string]Agent{"a": {Model: "x/y"}},
	}
	for i := 0; i < MaxRecentModels+5; i++ {
		cfg.AddRecentModel("agent/a", "p/m"+string(rune('a'+i)))
	}
	if got := len(cfg.Agents["a"].RecentModels); got != MaxRecentModels {
		t.Errorf("want cap %d, got %d", MaxRecentModels, got)
	}
}

func TestAddRecentModelIgnoresBadNodeID(t *testing.T) {
	cfg := &Config{Agents: map[string]Agent{}, Categories: map[string]Category{}}
	// Must not panic / not mutate anything.
	cfg.AddRecentModel("bogus", "p/m")
	cfg.AddRecentModel("agent/missing", "p/m")
	cfg.AddRecentModel("", "p/m")
}

func TestExpandPathTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	if got := expandPath("~/foo.json"); got != filepath.Join(home, "foo.json") {
		t.Errorf("got %s", got)
	}
	if got := expandPath("/tmp/x.json"); got != "/tmp/x.json" {
		t.Errorf("got %s", got)
	}
}
