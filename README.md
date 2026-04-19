# omo-switch

English | [中文](README_CN.md)

Terminal UI for switching the AI model used by each agent and category in
[Oh My OpenAgent (OpenCode)](https://opencode.ai) — for example, changing the
model for Sisyphus, Oracle, or any other agent defined in your
`oh-my-openagent.json` config. Vim-style navigation, grouped provider → model
picker, batch switching, atomic save.

<img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go" alt="Go 1.21+">
<img src="https://img.shields.io/badge/License-MIT-green" alt="MIT License">

## Install

Requires Go 1.21+.

```bash
go install github.com/kaiwenyao/omo-switch@latest
```

Or build from source:

```bash
git clone https://github.com/kaiwenyao/omo-switch.git
cd omo-switch
go build -o omo-switch .
```

## Usage

```bash
omo-switch                                         # default: ~/.config/opencode/oh-my-openagent.json
omo-switch -config /path/to/oh-my-openagent.json   # custom path
```

The app reads your `oh-my-openagent.json` and presents all agents (Sisyphus,
Oracle, Momus, …) and categories grouped by their current model
provider. Navigate to any node, pick a new model, and the config is saved
atomically — no manual JSON editing required.

### Keys

| Key | Action |
|---|---|
| `j` / `k` / arrows | Move cursor |
| `h` / `l` / arrows | Collapse / expand to parent or first child |
| `Space` | Toggle batch selection on current node |
| `Enter` | Open model picker for current node (or for the selection, if the cursor is inside it) |
| `x` | Clear selection |
| `q` / `Esc` | Back / quit |

Inside the picker: `j`/`k` + `Enter` to drill provider → model; `Esc` steps back
one level.

### Behavior notes

- **Atomic save.** Writes go to a sibling temp file, `fsync`, then `rename`,
  with `chmod 0600` on the temp file — the config may contain auth tokens, so
  permissions are tightened and partial writes cannot corrupt the original.
- **Unknown fields preserved.** Any keys the UI does not know about are kept
  verbatim across save roundtrips.
- **Batch edit.** When the cursor is *inside* the current selection, picking a
  model updates every selected node's primary model *and* all its
  `fallback_models`. When the cursor is *outside* the selection, the selection
  is dropped and the edit targets only the cursor node.
- **Recent models.** Each agent / category keeps a per-node `recent_models`
  list (dedup, capped at 10) surfaced as a synthetic "recently" provider at
  the top of the picker.

## Development

```bash
go build ./...
go vet  ./...
go test -race -cover ./...
```

## Repository settings

| Setting | Value | Rationale |
|---|---|---|
| Visibility | **Public** | No secrets in the tree; the Go import path `github.com/kaiwenyao/omo-switch` already implies a public module. |
| Default branch | `main` | Matches local convention. |
| Description | "Terminal UI for switching Oh My OpenAgent models — vim keys, grouped picker, atomic save." | Short, searchable. |
| Topics | `go`, `golang`, `tui`, `bubble-tea`, `cli`, `opencode` | Discoverability. |
| Issues | Enabled | Bug tracker. |
| Wiki | Disabled | Docs live in this repo. |
| Discussions | Disabled | Low volume; revisit if community grows. |
| Merge style | Default (all enabled) | Not opinionated yet; may tighten later. |
| License | **MIT** | See [LICENSE](LICENSE). |

## Project structure

```
config/   JSON load/save + recent-model tracking
tui/      Bubble Tea model, key handler, view, model picker
main.go   Entry point + CLI flag parsing
```

## License

[MIT](LICENSE)