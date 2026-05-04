# AIC CLI

> 🚀 管理 Claude Skills 的强大命令行工具 — 搜索、安装、管理，一站式解决方案

[![Go Report Card](https://goreportcard.com/badge/github.com/cicbyte/aic-cli)](https://goreportcard.com/report/github.com/cicbyte/aic-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23.2-blue)](https://go.dev/dl/)

**[🇨🇳 中文版](README.md)** | **[🇺🇸 English](README_EN.md)**

<!-- screenshot: 在此处添加终端录屏或截图 -->

## ✨ 功能特性

- **🔍 智能搜索** — 从 AIC 服务器快速搜索 skills，支持关键词和分类筛选
- **📦 一键安装** — 下载并自动安装到项目或全局目录，支持 zip 包和链接导入
- **🎨 精美 TUI** — 内置交互式终端界面，可视化浏览和操作
- **🗂️ 分类管理** — 按 category 组织 skills，轻松找到所需工具
- **🔗 符号链接** — 支持软链接管理，节省存储空间
- **📊 详细信息** — 显示版本、下载量、收藏数等详细信息
- **🧹 清理工具** — 自动清理未使用的 skills，保持环境整洁
- **⚙️ 灵活配置** — 支持项目级和全局 skills 目录，自动检测 `.claude` 目录

## 📦 安装

### 环境要求

- Go 1.23.2 或更高版本

### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/cicbyte/aic-cli.git
cd aic-cli

# 编译安装
go build -o aic-cli main.go

# 移动到 PATH
sudo mv aic-cli /usr/local/bin/
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
aic-cli search claude

# 安装 skill
aic-cli add 123

# 列出已安装的 skills
aic-cli list
```

## 📖 使用方法

### `aic-cli`

启动交互式 TUI 界面，可视化浏览和管理 skills。

```bash
aic-cli
```

不带任何参数运行时，自动进入交互模式，提供友好的终端用户界面。

### `aic-cli search [keyword]`

搜索 skills，支持关键词和分类筛选。

```bash
aic-cli search [关键词] [选项]
```

| 选项 | 别名 | 说明 | 默认值 |
|---|---|---|---|
| `--category` | `-c` | 按分类 ID 筛选 | `0` |
| `--page` | `-p` | 页码 | `1` |
| `--size` | `-n` | 每页数量 | `20` |

**示例：**

```bash
# 搜索所有 skills
aic-cli search

# 按关键词搜索
aic-cli search claude

# 按分类筛选
aic-cli search --category 2

# 分页浏览
aic-cli search --page 2 --size 10
```

### `aic-cli add <id>`

下载并安装 skill 到本地。

```bash
aic-cli add <skill-id>
```

**示例：**

```bash
# 安装指定 skill
aic-cli add 123

# 安装到全局目录
aic-cli add 123 --global
```

### `aic-cli list`

列出本地已安装的 skills。

```bash
aic-cli list [选项]
```

| 选项 | 别名 | 说明 | 默认值 |
|---|---|---|---|
| `--global` | `-g` | 显示全局 skills 目录 | `false` |

**示例：**

```bash
# 列出当前项目 skills
aic-cli list

# 列出全局 skills
aic-cli list -g
```

### `aic-cli remove <name>`

删除已安装的 skill。

```bash
aic-cli remove <skill-name>
```

### `aic-cli categories`

列出所有 skill 分类。

```bash
aic-cli categories
```

### `aic-cli download <id>`

下载 skill zip 包到当前目录。

```bash
aic-cli download <skill-id>
```

### `aic-cli clean`

清理未使用的 skills 和缓存。

```bash
aic-cli clean
```

### `aic-cli import <url>`

从 URL 导入 skill。

```bash
aic-cli import <url>
```

### 全部命令

```bash
aic-cli --help
```

| 命令 | 说明 |
|---|---|
| `search` | 搜索 skills |
| `add` | 安装 skill |
| `remove` | 删除 skill |
| `list` | 列出已安装的 skills |
| `categories` | 查看分类 |
| `download` | 下载 zip 包 |
| `clean` | 清理缓存 |
| `import` | 从 URL 导入 |
| `tui` | 交互式界面 |

## ⚙️ 配置

### 目录结构

AIC CLI 支持两种 skills 目录：

- **项目级目录**：`<project>/.claude/skills/`
  - 自动检测当前项目中的 `.claude` 目录
  - 适合项目特定的 skills

- **全局目录**：`~/.aic-cli/skills/`
  - 跨项目共享的通用 skills
  - 使用 `--global` 或 `-g` 参数访问

### 配置文件

配置文件位于 `~/.aic-cli/config.yaml`：

```yaml
# AIC 服务器配置
aic:
  base_url: "https://api.example.com"  # AIC 服务器地址
  timeout: 30

log:
  level: "info"
  path: "~/.aic-cli/logs/aic-cli.log"
```

### 环境变量

| 变量 | 说明 | 默认值 |
|---|---|---|
| `AIC_CLI_BASE_URL` | AIC 服务器 API 基础 URL | — |
| `AIC_CLI_CONFIG` | 配置文件路径 | `~/.aic-cli/config.yaml` |
| `AIC_CLI_LOG_LEVEL` | 日志级别 | `info` |

## 🏗️ 项目结构

```
aic-cli/
├── cmd/              # CLI 命令定义
│   ├── root.go       # 根命令
│   ├── search.go     # 搜索命令
│   ├── add.go        # 添加命令
│   ├── list.go       # 列表命令
│   └── tui.go        # TUI 界面
├── internal/         # 内部包
│   ├── api/          # API 客户端
│   ├── common/       # 公共定义
│   ├── log/          # 日志模块
│   └── utils/        # 工具函数
└── main.go           # 程序入口
```

## 🛠️ 开发

### 构建项目

```bash
# 快速构建（仅 Go）
go build -o aic-cli main.go

# 完整构建（web + Go）
python build.py

# 交叉编译三平台
python build_new.py
```

详细的构建和发布流程请查看：[scripts/README.md](scripts/README.md)

### 运行测试

```bash
go test ./...
```

### 贡献指南

欢迎提交 Issue 和 Pull Request！在提交 PR 前，请确保：

1. 代码通过 `go test` 和 `go vet`
2. 添加必要的测试用例
3. 更新相关文档

## 📄 开源许可证

[MIT](LICENSE) © 2026 Cicbyte

## 🙏 致谢

- [Cobra](https://github.com/spf13/cobra) — 强大的 CLI 框架
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — 精美的 TUI 框架
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — 终端样式库

---

**项目地址：** [https://github.com/cicbyte/aic-cli](https://github.com/cicbyte/aic-cli)

**反馈与建议：** 请提交 [Issue](https://github.com/cicbyte/aic-cli/issues)
