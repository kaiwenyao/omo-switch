# CLAUDE.md

## Project Overview

**omo-switch** — A terminal UI app for editing oh-my-openagent JSON config with vim navigation and batch model switching.

## Tech Stack

- **Language:** Go
- **TUI:** Bubble Tea (github.com/charmbracelet/bubbletea)
- **Target:** macOS, single binary

## Config File

`~/.config/opencode/oh-my-openagent.json`

Structure:
```json
{
  "agents": { "<name>": { "model": "...", "fallback_models": [...] } },
  "categories": { "<name>": { "model": "...", "variant": "...", "fallback_models": [...] } }
}
```

## Requirements (v1)

1. Parse JSON config into in-memory model
2. Render as navigable tree (Providers → Agents/Categories → Model)
3. Vim-style navigation (h/j/k/l) + arrow keys
4. Batch selection (spacebar)
5. Type-to-filter OR free-form model input
6. Auto-save with fallback_models update
7. CLI arg for custom config path

## Key Decisions

- **Bubble Tea:** Modern Go TUI library
- **Vim keys:** User-specified navigation
- **Auto-save:** Immediate JSON write on edit
- **Fallback update:** All fallback_models get same model when batch updating

## Workflow

This project uses GSD workflow:
1. `/gsd-plan-phase 1` — plan the phase
2. `/gsd-execute-phase 1` — execute plans
3. `/gsd-verify-phase 1` — verify (if verifier enabled)
