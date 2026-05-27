# AIC CLI

> 🚀 管理 Coding Agent Skills 的命令行工具 — 搜索、安装、管理，一站式解决方案

[![Go Report Card](https://goreportcard.com/badge/github.com/cicbyte/aic-cli)](https://goreportcard.com/report/github.com/cicbyte/aic-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23.2-blue)](https://go.dev/dl/)

**[🇨🇳 中文版](README.md)** | **[🇺🇸 English](README_EN.md)**

## ✨ 功能特性

- **🔍 智能搜索** — 从 AIC 服务器快速搜索 skills，支持关键词和分类筛选
- **📦 一键安装** — 下载并自动安装到项目或全局目录
- **🤖 多 Agent 支持** — 支持 18 种 Coding Agent（Claude Code、Cursor、Windsurf、Cline 等）
- **🎨 精美 TUI** — 内置交互式终端界面，可视化浏览和操作
- **🗂️ 技能包** — 按技能包批量安装，一键配齐开发环境
- **🔗 符号链接** — 支持软链接管理，多项目共享同一份 skill，节省存储空间
- **🧹 清理工具** — 自动清理失效的软链接，保持环境整洁

## 📦 安装

### 环境要求

- Go 1.23.2 或更高版本

### 从源码安装

```bash
git clone https://github.com/cicbyte/aic-cli.git
cd aic-cli
go build -o aic-cli main.go
sudo mv aic-cli /usr/local/bin/  # Linux/macOS
```

### 使用 Go install

```bash
go install github.com/cicbyte/aic-cli@latest
```

### 从发布版本下载

访问 [Releases](https://github.com/cicbyte/aic-cli/releases) 页面下载适合你系统的二进制文件。

## 🚀 快速开始

```bash
# 交互式 TUI 界面（推荐）
aic-cli

# 搜索 skills
aic-cli skill search "前端"

# 安装 skill
aic-cli skill add obsidian

# 安装技能包
aic-cli package add 1

# 列出已安装的 skills
aic-cli skill list
```

## 📖 使用方法

### Skill 管理

#### 搜索 skill

```bash
aic-cli skill search "关键词"
aic-cli skill search "前端" -c 1          # 按分类筛选
aic-cli skill search "react" -p 1 -n 10   # 分页
```

#### 添加 skill

```bash
aic-cli skill add <skill-id>
aic-cli skill add <skill-name>
aic-cli skill add "obsidian" --agent cursor    # 指定目标 Agent
aic-cli skill add "obsidian" -o ./custom-dir   # 指定输出目录
aic-cli skill add "obsidian" -m symlink        # symlink 模式
```

![skill-add](images/skill-add.gif)

使用 `--agent` 指定目标 Agent 时，skill 会安装到对应 Agent 的 skills 目录：

![skill-add-with-agent](images/skill-add-with-agent.gif)

#### 下载 skill zip 包

```bash
aic-cli skill download <skill-id>
aic-cli skill download "obsidian" -o ./downloads
```

#### 导入本地技能到服务器

```bash
aic-cli skill import ./my-skill.zip
aic-cli skill import ./my-skill.zip -d "技能描述" -c 1
```

#### 列出已安装的 skill

```bash
aic-cli skill list               # 当前项目
aic-cli skill list -g            # 全局目录
aic-cli skill list --agent cursor # 指定 Agent
```

#### 移除已安装的 skill

```bash
aic-cli skill remove <skill-name>
aic-cli skill remove <skill-name> -g            # 同时删除全局源文件
```

#### 清理失效软链接

```bash
aic-cli skill clean
```

#### 打包 skill

```bash
aic-cli skill pack ./my-skill/                  # 默认 .zip
aic-cli skill pack ./my-skill/ --format skill   # .skill 格式
```

#### 安装模式

```bash
aic-cli skill mode               # 查看当前模式
aic-cli skill mode symlink       # 切换到 symlink 模式
aic-cli skill mode copy          # 切换到 copy 模式
```

### 技能包管理

```bash
aic-cli package list                    # 列出所有技能包
aic-cli package list -s "前端"          # 搜索技能包
aic-cli package add <package-id>        # 按 ID 添加
aic-cli package add "前端工具"           # 按名称添加
aic-cli package add 1 --agent cursor    # 指定目标 Agent
```

![package-add](images/package-add.gif)

### 其他命令

```bash
aic-cli server open    # 在浏览器中打开 AIC 服务页面
aic-cli server status  # 查看服务器连接状态
aic-cli version        # 显示版本信息
```

### 全部命令

```bash
aic-cli --help
```

| 命令 | 说明 |
|---|---|
| `skill search` | 搜索远程技能 |
| `skill add` | 下载并安装技能 |
| `skill download` | 下载技能 ZIP 包 |
| `skill import` | 导入本地技能到服务器 |
| `skill categories` | 列出所有分类 |
| `skill list` | 列出已安装的技能 |
| `skill remove` | 移除已安装的技能 |
| `skill clean` | 清理失效的软链接 |
| `skill pack` | 打包技能文件夹 |
| `skill mode` | 查看或切换安装模式 |
| `package list` | 列出技能包 |
| `package add` | 安装技能包中的所有技能 |
| `server open` | 打开 AIC 服务页面 |
| `server status` | 查看服务器状态 |
| `tui` | 交互式界面 |

## 🤖 多 Agent 支持

AIC CLI 支持将 skill 安装到不同 Coding Agent 的目录中。使用 `--agent` 参数指定目标 Agent：

```bash
aic-cli skill add "react-helper" --agent cursor
aic-cli skill add "react-helper" --agent windsurf
```

支持的 Agent：

| Agent | 项目级目录 | 全局目录 |
|-------|-----------|---------|
| claude | `.claude/skills/` | `~/.claude/skills/` |
| cursor | `.cursor/rules/` | `~/.cursor/rules/` |
| windsurf | `.windsurf/skills/` | `~/.windsurf/skills/` |
| cline | `.cline/skills/` | `~/.cline/skills/` |
| continue | `.continue/skills/` | `~/.continue/skills/` |
| opencode | `.opencode/skills/` | `~/.config/opencode/skills/` |
| trae | `.trae/skills/` | `~/.trae/skills/` |
| gemini | `.gemini/skills/` | `~/.gemini/skills/` |
| codex | `.codex/skills/` | `~/.codex/skills/` |
| roo | `.roo/skills/` | `~/.roo/skills/` |
| amp | `.amp/skills/` | `~/.amp/skills/` |
| amazonq | `.amazonq/skills/` | `~/.amazonq/skills/` |
| copilot | `.github/prompts/` | `~/.github/prompts/` |
| qoder | `.qoder/skills/` | `~/.qoder/skills/` |
| openclaw | `.openclaw/skills/` | `~/.openclaw/skills/` |
| hermes | — | `~/.hermes/skills/` |
| codebuddy | — | `~/.codebuddy/skills/` |
| qclaw | — | `~/.qclaw/skills/` |

不指定 `--agent` 时，自动检测当前项目的 Agent，多个时交互选择，默认使用 Claude Code。

## ⚙️ 配置

### 配置文件

配置文件位于 `~/.ciclebyte/aic-cli/config/config.yaml`：

```yaml
aic:
  base_url: "https://aic.cicbyte.com"
  token: ""

skills:
  default_mode: "symlink"    # 安装模式: symlink 或 copy
  default_agent: "claude"    # 默认 Agent
```

### 安装模式

- **symlink**（默认）— 下载到全局目录，创建软连接到目标目录。多项目共享，节省磁盘空间。Windows 使用 Junction，无需管理员权限。
- **copy** — 直接复制文件到目标目录。各项目独立，互不影响。

## 🏗️ 项目结构

```
aic-cli/
├── cmd/                  # CLI 命令定义（Cobra）
│   ├── skill/            # skill 子命令
│   ├── package/          # package 子命令
│   ├── server/           # server 子命令
│   ├── local/            # list/remove/clean
│   ├── skillzip/         # pack
│   └── tui/              # 交互式界面
├── internal/
│   ├── agent/            # 18 个 Agent Profile 实现
│   ├── api/              # AIC 服务器 API 客户端
│   ├── logic/            # 业务逻辑层
│   ├── utils/            # 工具函数
│   ├── models/           # 数据模型
│   └── log/              # 日志模块
├── skills/               # 项目内置 skills
└── main.go
```

## 🛠️ 开发

```bash
# 构建
go build -o aic-cli main.go

# 测试
go vet ./...
go test ./...
```

## 📄 开源许可证

[MIT](LICENSE) © 2026 Cicbyte

## 🙏 致谢

- [Cobra](https://github.com/spf13/cobra) — CLI 框架
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — TUI 框架
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — 终端样式库

---

**项目地址：** [https://github.com/cicbyte/aic-cli](https://github.com/cicbyte/aic-cli)

**反馈与建议：** 请提交 [Issue](https://github.com/cicbyte/aic-cli/issues)
