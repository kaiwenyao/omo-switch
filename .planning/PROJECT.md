# omo-switch

## What This Is

A terminal UI application for visually editing the oh-my-openagent JSON config file. Read the config, render it in a TUI, let users navigate and edit model assignments via keyboard, and write changes back to the JSON file.

## Core Value

Quickly switch which model any agent or category uses, with batch selection and auto-save.

## Requirements

### Active

- [ ] **TUI rendering**: Read `~/.config/opencode/oh-my-openagent.json` and display as navigable tree (providers → agents/categories → model field)
- [ ] **Vim-style navigation**: `h/j/k/l` or arrow keys to move, `enter` to select/edit
- [ ] **Batch selection**: Multi-select multiple agents/categories, then assign one model to all selected
- [ ] **Model input**: Type-to-filter from a list OR type freely (no predefined list needed)
- [ ] **Fallback update**: When assigning model, also update all `fallback_models` entries to same model
- [ ] **Auto-save**: Changes written immediately to JSON file on edit
- [ ] **Config path**: Default to `~/.config/opencode/oh-my-openagent.json`, allow override via CLI arg

### Out of Scope

- Schema validation (trust the file is valid JSON)
- Backup/undo (auto-save only, no version control)
- Adding/removing agents or categories (edit existing entries only)
- Creating new model providers (read-only provider list from config)

## Context

- The oh-my-openagent config file contains `agents` and `categories` dictionaries
- Each agent/category has a `model` field and optional `fallback_models` array
- Model strings use format: `provider/model-name` (e.g., `ollama-cloud/glm-5.1`, `minimax-cn-coding-plan/MiniMax-M2.7`)
- Existing file has 10 agents and 8 categories across ~10 providers

## Constraints

- **Tech**: Go + Bubble Tea TUI library
- **Platform**: macOS (existing config path)
- **Config**: Single JSON file, path from env or arg

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Bubble Tea | Modern Go TUI, good for interactive lists, user deferred to me | — Pending |
| Vim keys | User specified vim-style navigation | — Pending |
| Auto-save | User wants quick workflow, no confirmation step | — Pending |
| Update fallbacks | User confirmed all fallback_models should also update | — Pending |

---
*Last updated: 2026-04-18 after initialization*
