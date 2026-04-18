# ROADMAP.md

## Phase 1: Core TUI App

**Goal:** Read oh-my-openagent JSON config and render editable TUI with vim navigation.

### Requirements

- [ ] **OMO-01**: Parse `~/.config/opencode/oh-my-openagent.json` into in-memory model
- [x] **OMO-02**: Render tree view: Providers → Agents/Categories → Model field
- [ ] **OMO-03**: Vim-style navigation (h/j/k/l) + arrow key support
- [ ] **OMO-04**: Batch selection via spacebar or similar
- [ ] **OMO-05**: Type-to-filter model input OR free-form entry
- [ ] **OMO-06**: Auto-save changes to JSON on edit (with fallback_models update)
- [ ] **OMO-07**: CLI arg for custom config path

### Success Criteria

1. App launches and displays config structure in TUI
2. User can navigate with vim keys and select multiple entries
3. Model assignment updates primary + all fallback_models
4. Changes persist to JSON file immediately
5. Works with existing config structure (10 agents, 8 categories)

### Plans

- [x] 01-01-PLAN.md — Go module init + config structs + JSON parsing
- [x] 01-02-PLAN.md — Bubble Tea app shell + tree model + view rendering
- [ ] 01-03-PLAN.md — Vim navigation (h/j/k/l) + batch selection (spacebar)
- [ ] 01-04-PLAN.md — Model editing + fallback_models update + auto-save

---

**1 phases** | **7 requirements mapped** | **4 plans created**

## Traceability

| REQ | Phase | Plan |
|-----|-------|------|
| OMO-01 | 1 | 01-01 |
| OMO-02 | 1 | 01-02 |
| OMO-03 | 1 | 01-03 |
| OMO-04 | 1 | 01-03 |
| OMO-05 | 1 | 01-04 |
| OMO-06 | 1 | 01-04 |
| OMO-07 | 1 | 01-01 |
