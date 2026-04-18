# ROADMAP.md

## Phase 1: Core TUI App

**Goal:** Read oh-my-openagent JSON config and render editable TUI with vim navigation.

### Requirements

- [ ] **OMO-01**: Parse `~/.config/opencode/oh-my-openagent.json` into in-memory model
- [ ] **OMO-02**: Render tree view: Providers → Agents/Categories → Model field
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

---

**1 phases** | **7 requirements mapped** | All v1 requirements covered ✓

## Traceability

| REQ | Phase |
|-----|-------|
| OMO-01 | 1 |
| OMO-02 | 1 |
| OMO-03 | 1 |
| OMO-04 | 1 |
| OMO-05 | 1 |
| OMO-06 | 1 |
| OMO-07 | 1 |
