package tui

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kaiwenyao/omo-switch/config"
)

func emptyConfig() *config.Config {
	return &config.Config{
		Agents:     map[string]config.Agent{},
		Categories: map[string]config.Category{},
	}
}

func TestArrowKeysOnEmptyTreeDoNotPanic(t *testing.T) {
	m := NewAppModel(emptyConfig())
	keys := []tea.KeyType{tea.KeyLeft, tea.KeyRight, tea.KeyEnter, tea.KeySpace}
	for _, k := range keys {
		_, _ = HandleKey(tea.KeyMsg{Type: k}, &m)
	}
	// j/k rune navigation also must not panic
	for _, r := range []rune{'j', 'k', 'h', 'l', 'x'} {
		_, _ = HandleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}, &m)
	}
}

func TestApplyChangeAgent(t *testing.T) {
	cfg := &config.Config{
		Agents: map[string]config.Agent{
			"foo": {
				Model: "old/model",
				FallbackModels: []config.FallbackModel{
					{Model: "old/fb1"},
					{Model: "old/fb2"},
				},
			},
		},
		Categories: map[string]config.Category{},
	}

	// Single-node edit: only primary model changes, fallbacks preserved
	ApplyChange(cfg, "agent/foo", "new/model", false)
	a := cfg.Agents["foo"]
	if a.Model != "new/model" {
		t.Errorf("primary: want new/model, got %s", a.Model)
	}
	if a.FallbackModels[0].Model != "old/fb1" {
		t.Errorf("fb[0]: want old/fb1, got %s", a.FallbackModels[0].Model)
	}
	if a.FallbackModels[1].Model != "old/fb2" {
		t.Errorf("fb[1]: want old/fb2, got %s", a.FallbackModels[1].Model)
	}

	// Batch edit: primary and all fallbacks change
	ApplyChange(cfg, "agent/foo", "batch/model", true)
	a = cfg.Agents["foo"]
	if a.Model != "batch/model" {
		t.Errorf("primary: want batch/model, got %s", a.Model)
	}
	for i, fb := range a.FallbackModels {
		if fb.Model != "batch/model" {
			t.Errorf("fb[%d]: want batch/model, got %s", i, fb.Model)
		}
	}
}

func TestApplyChangeCategory(t *testing.T) {
	cfg := &config.Config{
		Agents: map[string]config.Agent{},
		Categories: map[string]config.Category{
			"code": {
				Model:          "a/b",
				FallbackModels: []config.FallbackModel{{Model: "c/d"}},
			},
		},
	}
	// Single edit: fallbacks preserved
	ApplyChange(cfg, "category/code", "x/y", false)
	c := cfg.Categories["code"]
	if c.Model != "x/y" {
		t.Errorf("primary: want x/y, got %s", c.Model)
	}
	if c.FallbackModels[0].Model != "c/d" {
		t.Errorf("fb[0]: want c/d, got %s", c.FallbackModels[0].Model)
	}
	// Batch edit: fallbacks updated
	ApplyChange(cfg, "category/code", "z/w", true)
	c = cfg.Categories["code"]
	if c.Model != "z/w" || c.FallbackModels[0].Model != "z/w" {
		t.Errorf("batch: category not updated: %+v", c)
	}
}

func TestApplyChangeMissingNodeIsNoop(t *testing.T) {
	cfg := emptyConfig()
	ApplyChange(cfg, "agent/nope", "x", false)
	ApplyChange(cfg, "category/nope", "x", false)
	ApplyChange(cfg, "malformed", "x", false)
	ApplyChange(cfg, "provider/skipped", "x", false)
}

func TestBuildTreeGroupsByProvider(t *testing.T) {
	cfg := &config.Config{
		Agents: map[string]config.Agent{
			"a1": {Model: "p1/m"},
			"a2": {Model: "p2/m"},
		},
		Categories: map[string]config.Category{
			"c1": {Model: "p1/m"},
		},
	}
	roots := BuildTree(cfg)
	if len(roots) != 2 {
		t.Fatalf("want 2 providers, got %d", len(roots))
	}
	// Providers must be sorted; p1 has 2 children, p2 has 1.
	if roots[0].Label != "p1" || len(roots[0].Children) != 2 {
		t.Errorf("p1: got label=%s children=%d", roots[0].Label, len(roots[0].Children))
	}
	if roots[1].Label != "p2" || len(roots[1].Children) != 1 {
		t.Errorf("p2: got label=%s children=%d", roots[1].Label, len(roots[1].Children))
	}
}

// Regression: if the cursor is outside an existing selection, pressing Enter
// must edit the cursor node alone, not the stale selection. Previously
// ConfirmEdit would silently rewrite the selection while the user was staring
// at a different node.
func TestStartModelPickerOnUnselectedCursorClearsSelection(t *testing.T) {
	cfg := &config.Config{
		Agents: map[string]config.Agent{
			"a": {Model: "p/m"},
			"b": {Model: "p/m"},
			"d": {Model: "p/m"},
		},
		Categories: map[string]config.Category{},
	}
	m := NewAppModel(cfg)
	m.ToggleSelection("agent/a")
	m.ToggleSelection("agent/b")
	if m.SelectionCount != 2 {
		t.Fatalf("setup: want 2 selected, got %d", m.SelectionCount)
	}

	// Cursor lands on an unselected node; opening the picker must drop the selection.
	m.StartModelPicker("agent/d", "p/m")
	if m.SelectionCount != 0 {
		t.Errorf("selection not cleared: count=%d", m.SelectionCount)
	}
	if m.EditNodeID != "agent/d" {
		t.Errorf("EditNodeID: want agent/d, got %s", m.EditNodeID)
	}

	// Confirming now must modify only d, leaving a and b untouched.
	m.EditValue = "p/new"
	if err := m.ConfirmEdit(); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if got := cfg.Agents["d"].Model; got != "p/new" {
		t.Errorf("d not updated: got %s", got)
	}
	if got := cfg.Agents["a"].Model; got != "p/m" {
		t.Errorf("a wrongly updated: got %s", got)
	}
	if got := cfg.Agents["b"].Model; got != "p/m" {
		t.Errorf("b wrongly updated: got %s", got)
	}
}

// Cursor inside selection: batch edit proceeds and touches every selected node.
func TestStartModelPickerOnSelectedCursorKeepsSelection(t *testing.T) {
	cfg := &config.Config{
		Agents: map[string]config.Agent{
			"a": {Model: "p/m"},
			"b": {Model: "p/m"},
			"c": {Model: "p/m"},
		},
		Categories: map[string]config.Category{},
	}
	m := NewAppModel(cfg)
	m.ToggleSelection("agent/a")
	m.ToggleSelection("agent/b")

	m.StartModelPicker("agent/a", "p/m")
	if m.SelectionCount != 2 {
		t.Errorf("selection dropped on in-selection cursor: count=%d", m.SelectionCount)
	}

	m.EditValue = "p/new"
	if err := m.ConfirmEdit(); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if got := cfg.Agents["a"].Model; got != "p/new" {
		t.Errorf("a: got %s", got)
	}
	if got := cfg.Agents["b"].Model; got != "p/new" {
		t.Errorf("b: got %s", got)
	}
	if got := cfg.Agents["c"].Model; got != "p/m" {
		t.Errorf("c wrongly updated: got %s", got)
	}
	if m.SelectionCount != 0 {
		t.Errorf("selection not cleared after confirm: count=%d", m.SelectionCount)
	}
}

// Regression for BUG 3: save failure must roll back the in-memory mutation
// and keep the picker/selection open so the user's change is not silently
// lost when they later press q / Esc.
func TestConfirmEditRollsBackAndKeepsPickerOnSaveFailure(t *testing.T) {
	cfg := &config.Config{
		Agents: map[string]config.Agent{
			"a": {Model: "p/original", FallbackModels: []config.FallbackModel{{Model: "p/fb"}}},
			"b": {Model: "p/original", FallbackModels: []config.FallbackModel{{Model: "p/fb"}}},
		},
		Categories: map[string]config.Category{},
	}
	m := NewAppModel(cfg)
	// Two selected nodes, then open picker on one of them so ConfirmEdit
	// takes the batch path.
	m.ToggleSelection("agent/a")
	m.ToggleSelection("agent/b")
	m.StartModelPicker("agent/a", "p/original")
	// Point the save at a path whose parent does not exist — CreateTemp fails.
	m.ConfigPath = "/nonexistent-dir-omo-switch-test/config.json"
	m.EditValue = "p/new"

	err := m.ConfirmEdit()
	if err == nil {
		t.Fatal("want save error, got nil")
	}

	// Rollback: in-memory config must match original disk state.
	if got := cfg.Agents["a"].Model; got != "p/original" {
		t.Errorf("a not rolled back: got %s", got)
	}
	if got := cfg.Agents["b"].Model; got != "p/original" {
		t.Errorf("b not rolled back: got %s", got)
	}
	if got := cfg.Agents["a"].FallbackModels[0].Model; got != "p/fb" {
		t.Errorf("a fallback not rolled back: got %s", got)
	}

	// Picker and selection preserved so user can retry or cancel.
	if !m.ModelPickerMode {
		t.Error("ModelPickerMode dropped on save failure — user would lose context")
	}
	if m.SelectionCount != 2 {
		t.Errorf("selection lost on save failure: count=%d", m.SelectionCount)
	}
	if m.EditNodeID != "agent/a" {
		t.Errorf("EditNodeID lost: %s", m.EditNodeID)
	}
	if m.ErrorMsg == "" {
		t.Error("ErrorMsg empty — user has no indication of failure")
	}
}

// After a failed save, retrying on a now-writable path must succeed and
// fully commit the change, confirming the retry loop is not broken by the
// rollback path.
func TestConfirmEditRetrySucceedsAfterFailure(t *testing.T) {
	cfg := &config.Config{
		Agents: map[string]config.Agent{
			"a": {Model: "p/original"},
		},
		Categories: map[string]config.Category{},
	}
	m := NewAppModel(cfg)
	m.StartModelPicker("agent/a", "p/original")
	m.ConfigPath = "/nonexistent-dir-omo-switch-test/config.json"
	m.EditValue = "p/new"
	if err := m.ConfirmEdit(); err == nil {
		t.Fatal("want first attempt to fail")
	}

	// Retry with a real path.
	dir := t.TempDir()
	path := dir + "/cfg.json"
	if err := writeInitialConfig(path); err != nil {
		t.Fatal(err)
	}
	m.ConfigPath = path
	if err := m.ConfirmEdit(); err != nil {
		t.Fatalf("retry: %v", err)
	}
	if got := cfg.Agents["a"].Model; got != "p/new" {
		t.Errorf("retry did not apply: got %s", got)
	}
	if m.ModelPickerMode {
		t.Error("picker still open after successful retry")
	}
	if m.ErrorMsg != "" {
		t.Errorf("ErrorMsg not cleared: %s", m.ErrorMsg)
	}
}

func writeInitialConfig(path string) error {
	return os.WriteFile(path, []byte(`{"agents":{"a":{"model":"p/original","fallback_models":[]}},"categories":{},"google_auth":false}`), 0600)
}

func TestProviderFromModel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"ollama-cloud/glm-5.1", "ollama-cloud"},
		{"solo", "solo"},
		{"a/b/c", "a"},
	}
	for _, c := range cases {
		if got := providerFromModel(c.in); got != c.want {
			t.Errorf("providerFromModel(%q): want %q, got %q", c.in, c.want, got)
		}
	}
}
