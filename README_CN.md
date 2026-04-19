# omo-switch

[English](README.md) | 中文

用于切换 [Oh My OpenAgent (OpenCode)](https://opencode.ai) 中各角色所用 AI 模型的终端 UI 工具——例如将西西弗斯（Sisyphus）、先知（Oracle）等角色的模型从 GPT-4 换成 Claude，或为多个角色批量切换模型。Vim 风格导航、按供应商分组的模型选择器、原子保存。

<img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go" alt="Go 1.21+">
<img src="https://img.shields.io/badge/License-MIT-green" alt="MIT License">

## 安装

需要 Go 1.21+。

```bash
go install github.com/kaiwenyao/omo-switch@latest
```

或从源码构建：

```bash
git clone https://github.com/kaiwenyao/omo-switch.git
cd omo-switch
go build -o omo-switch .
```

## 使用

```bash
omo-switch                                         # 默认: ~/.config/opencode/oh-my-openagent.json
omo-switch -config /path/to/oh-my-openagent.json   # 自定义路径
```

程序读取 `oh-my-openagent.json`，将所有角色（西西弗斯、先知、摩莫斯等）和分类按当前模型供应商分组展示。导航到任意节点，选择新模型，配置即刻原子保存——无需手动编辑 JSON。

### 快捷键

| 按键 | 操作 |
|---|---|
| `j` / `k` / 方向键 | 移动光标 |
| `h` / `l` / 方向键 | 收起 / 展开到父节点或第一个子节点 |
| `空格` | 切换当前节点的批量选择 |
| `回车` | 打开当前节点的模型选择器（若光标在选择区域内，则对整个选择区域生效） |
| `x` | 清除选择 |
| `q` / `Esc` | 返回 / 退出 |

选择器内：`j`/`k` + `回车` 逐级选择 供应商 → 模型；`Esc` 返回上一级。

### 行为说明

- **原子保存。** 写入时先写同级临时文件，`fsync` 后 `rename`，临时文件权限为 `0600`——配置文件可能包含认证令牌，因此权限收紧，且部分写入不会破坏原文件。
- **未知字段保留。** UI 不识别的任意字段在保存往返中完整保留。
- **批量编辑。** 当光标位于当前选择区域*内部*时，选择模型会更新所有已选节点的主模型*及其* `fallback_models`。当光标位于选择区域*外部*时，选择被清除，编辑仅针对光标所在节点。
- **最近模型。** 每个角色 / 分类维护一个节点级 `recent_models` 列表（去重，上限 10 个），在选择器顶部以合成的"recently"供应商展示。

## 开发

```bash
go build ./...
go vet  ./...
go test -race -cover ./...
```

## 项目结构

```
config/   JSON 读写 + 最近模型追踪
tui/      Bubble Tea 模型、按键处理、视图、模型选择器
main.go   入口 + CLI 参数解析
```

## 许可证

[MIT](LICENSE)