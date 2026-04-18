---
status: testing
phase: 01-core-tui-app
source: [01-01-SUMMARY.md, 01-02-SUMMARY.md, 01-03-SUMMARY.md, 01-04-SUMMARY.md, 01-05-SUMMARY.md]
started: 2026-04-18T19:25:00Z
updated: 2026-04-18T20:15:00Z
---

## Current Test

number: 10
name: Exit app with Escape
expected: |
  Press Escape in normal navigation mode to quit the app.
  ESC also cancels model picker and edit mode.
awaiting: user response

## Critical Issues Fixed (Second Round)

**Issue 1: ESC doesn't exit**
- `tui/keys.go`: ESC in normal mode now returns `tea.Quit` to exit app
- ESC also properly cancels model picker or edit mode

**Issue 2: Type-to-search conflicts with h/j/k/l**
- Removed character input filtering in model picker mode
- j/k navigation only works via arrow keys now (tea.KeyDown/KeyUp)
- Character input no longer interferes with navigation

**Issue 3: Model list incomplete**
- `tui/models.go`: `GetAvailableModels()` now runs `opencode models` command
- Gets the full list of 200+ available models from opencode
- Falls back to config-based extraction if command fails

**Files changed:**
- `tui/models.go` (modified - now calls opencode models)
- `tui/keys.go` (modified - removed type-to-search, added ESC quit)
- `tui/view.go` (modified - updated footers)

**Build:** `go build ./...` succeeds

## Tests

### 1. Launch and see tree view
expected: |
  App launches in terminal. You see a tree view showing:
  - Providers at top level (e.g., "ollama-cloud", "minimax-cn-coding-plan", "google")
  - Agents and categories indented under each provider
  - Model values shown in gray for each agent/category
result: pending

### 2. Navigate with vim keys
expected: |
  Press h/j/k/l to move cursor up/down/left/right through the tree.
  Arrow keys also work.
result: pending

### 3. Select with spacebar
expected: |
  Press Space on an agent or category to toggle selection.
  Selected items show [x] marker, unselected show [ ].
result: pending

### 4. Batch selection
expected: |
  Select multiple agents/categories with Space.
  Selection persists as you navigate.
result: pending

### 5. Edit model value
expected: |
  Press Enter on a selected agent/category.
  Current model value appears as editable text.
  Type to edit the model string (format: "provider/model-name").
  Press Enter to confirm, Escape to cancel.
result: pending

### 6. Fallback models update
expected: |
  When you change a model, ALL fallback_models entries for that
  agent/category also update to the same new model value.
result: pending

### 7. Auto-save persists
expected: |
  After confirming edit, changes are immediately written to
  ~/.config/opencode/oh-my-openagent.json
result: pending

### 8. Custom config path
expected: |
  Run: go run main.go -config /path/to/config.json
  App loads the specified config file instead of default.
result: pending

### 9. Model picker - two level provider > model
expected: |
  After pressing Enter on an agent/category:
  - Two-level model picker opens
  - LEFT: Provider list (j/k navigate)
  - RIGHT: Preview of selected provider's models
  - "recently" appears at top if agent has history
  - Enter enters provider → shows model list (j/k navigate)
  - Enter on model → confirms selection
  - Esc → back to provider list, or cancel picker
result: pending

### 10. Exit app with Escape
expected: |
  Press Escape in normal mode to quit the app.
  Esc also cancels model picker or edit mode.
result: pending

### 11. Recently used models
expected: |
  "recently" provider shows at top of provider list.
  Lists the last 10 models used for this agent.
  Automatically updates when a model is selected.
result: pending

## Summary

total: 11
passed: 0
issues: 0
pending: 11
skipped: 0

## Gaps

- [OMO-01](01-03-SUMMARY.md): Model picker feature was incomplete (only free-text entry, no list selection) — **FIXED** — needs UAT verification
- Type-to-search conflicted with navigation — **FIXED** — removed typing filter, kept j/k navigation
- ESC didn't exit app — **FIXED** — ESC now quits in normal mode
- Recently used models — **NEW** — added "recently" provider showing last 10 models used per agent
