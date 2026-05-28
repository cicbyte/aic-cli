---
name: aic-cli
description: |
  AIC CLI 命令行工具使用指南。管理 Claude Code 等 Coding Agent 的 skills。
  当用户需要搜索、下载、安装、导入、打包 skills，或管理技能包时使用。
  触发词："skill"、"技能"、"技能包"、"aic-cli"、"下载skill"、"安装skill"、"搜索skill"。
---

# AIC CLI 使用指南

AIC CLI 是用于管理 Coding Agent skills 的命令行工具，支持从 AIC 服务器搜索、下载、安装 skills，以及管理技能包。

## 命令总览

```
aic-cli skill           # skills 管理（搜索、下载、安装、导入、打包）
aic-cli skill remote    # 远程技能管理（AI Agent 专用：cat/tree/patch/validate/publish）
aic-cli package         # 技能包管理（列表、添加）
aic-cli server          # 服务器操作（打开页面、查看状态）
aic-cli version         # 显示版本信息
```

## Skills 管理（常用命令）

### 搜索 skill

```bash
aic-cli skill search "关键词"
aic-cli skill search "前端" -c 1        # 按分类筛选
aic-cli skill search "react" -p 1 -n 10 # 分页
```

### 添加 skill 到本地

```bash
aic-cli skill add <skill-id>
aic-cli skill add <skill-name>
aic-cli skill add "obsidian" --agent cursor    # 指定目标 Agent
aic-cli skill add "obsidian" -o ./custom-dir   # 指定输出目录
aic-cli skill add "obsidian" -m symlink        # symlink 模式
```

支持的 `--agent` 值：
- `claude`（默认）、`cursor`、`continue`、`amazonq`、`copilot`
- `windsurf`、`cline`、`opencode`、`codex`、`gemini`
- `roo`、`amp`、`trae`、`openclaw`、`qclaw`、`hermes`、`codebuddy`、`qoder`

不指定 `--agent` 时，自动检测当前项目的 Agent，多个时交互选择。

### 下载 skill zip 包

```bash
aic-cli skill download <skill-id>
aic-cli skill download "obsidian" -o ./downloads
```

### 导入本地技能到服务器

```bash
aic-cli skill import ./my-skill.zip
aic-cli skill import ./my-skill.skill
aic-cli skill import ./my-skill-folder/
aic-cli skill import ./my-skill.zip -d "技能描述" -c 1
```

### 查看分类

```bash
aic-cli skill categories
```

### 列出已安装的 skill

```bash
aic-cli skill list               # 当前项目
aic-cli skill list -g            # 全局目录
aic-cli skill list --agent cursor # 指定 Agent
```

### 移除已安装的 skill

```bash
aic-cli skill remove <skill-name>
aic-cli skill remove <skill-name> -g            # 同时删除全局源文件
aic-cli skill remove <skill-name> --agent cursor # 指定 Agent
```

### 清理失效软链接

```bash
aic-cli skill clean
```

### 打包 skill

```bash
aic-cli skill pack ./my-skill/                  # 默认 .zip
aic-cli skill pack ./my-skill/ --format skill   # .skill 格式
aic-cli skill pack ./my-skill/ -o output.zip    # 指定输出文件名
```

### 安装模式管理

```bash
aic-cli skill mode               # 查看当前模式及说明
aic-cli skill mode symlink       # 切换到 symlink 模式
aic-cli skill mode copy          # 切换到 copy 模式
```

## 远程技能管理（AI Agent 专用）

以下命令通过 `aic-cli skill remote` 访问，主要用于 AI Agent 自动化操作。

### 查看远程文件内容

```bash
aic-cli skill remote cat 42 SKILL.md               # 输出到 stdout
aic-cli skill remote cat 42 SKILL.md -o local.md   # 保存到本地文件
```

### 显示远程文件树

```bash
aic-cli skill remote tree 42
```

输出示例：
```
my-skill (ID: 42)
├── SKILL.md (1.2KB)
├── prompts/
│   ├── review.md (512B)
│   └── refactor.md (384B)
└── references/
    └── api.md (2.1KB)
```

### 增量编辑技能文件

```bash
# 纯内容匹配
aic-cli skill remote patch 42 --path SKILL.md \
  --old "description: 旧描述" \
  --new "description: 新描述"

# 行号 + 内容校验（推荐）
aic-cli skill remote patch 42 --path SKILL.md \
  --line 5 \
  --old "description: 旧描述" \
  --new "description: 新描述"

# 纯行号替换
aic-cli skill remote patch 42 --path SKILL.md \
  --line 20-25 \
  --new "## 新增章节\n\n替换全部内容"

# 通过管道传入 JSON（自动检测）
echo '{"path":"SKILL.md","edits":[...]}' | aic-cli skill remote patch 42

# 批量 patch 多个文件
aic-cli skill remote patch 42 --batch edits.json

# 查看 JSON 格式说明
aic-cli skill remote patch --schema
```

### 编辑远程技能文件

```bash
aic-cli skill remote edit 42                        # 交互式：列出文件供选择
aic-cli skill remote edit 42 SKILL.md              # 直接编辑指定文件
aic-cli skill remote edit 42 SKILL.md --dry-run    # 只显示差异，不提交
```

### 校验技能文件

```bash
aic-cli skill remote validate 42                    # 标准校验
aic-cli skill remote validate 42 --strict          # 严格模式
```

### 发布技能

```bash
aic-cli skill remote publish 42                     # 发布
aic-cli skill remote publish 42 --version 1.0.0    # 指定版本
aic-cli skill remote publish 42 --changelog "新增功能"  # 带说明
```

### 取消发布

```bash
aic-cli skill remote unpublish 42
```

### 更新技能元数据

```bash
aic-cli skill remote update 42 --name "新名称"
aic-cli skill remote update 42 --desc "新描述"
aic-cli skill remote update 42 --tags "go,cli,tool"
aic-cli skill remote update 42 --category 2
```

## 技能包管理

```bash
aic-cli package list                    # 列出所有技能包
aic-cli package list -s "前端"          # 搜索技能包
aic-cli package add <package-id>        # 按 ID 添加
aic-cli package add "前端工具"           # 按名称添加
aic-cli package add 1 --agent cursor    # 指定目标 Agent
```

## 服务器操作

```bash
aic-cli server open    # 在浏览器中打开 AIC 服务页面
aic-cli server status  # 查看服务器连接状态
```

## 多 Agent 支持

aic-cli 支持将 skill 安装到不同 Coding Agent 的目录中。每个 Agent 有项目级和全局两套 skills 目录：

| Agent | 项目级目录 | 全局目录 |
|-------|-----------|---------|
| claude | `.claude/skills/` | `~/.claude/skills/` |
| cursor | `.cursor/rules/` | `~/.cursor/rules/` |
| opencode | `.opencode/skills/` | `~/.config/opencode/skills/` |
| trae | `.trae/skills/` | `~/.trae/skills/` |

安装模式：
- `symlink`（默认）：下载到全局目录，创建软连接到目标目录（Windows 使用 Junction，无需管理员权限）
- `copy`：直接复制文件到目标目录

## 注意事项

- 使用前需先通过 `aic-cli server status` 确认服务器连接正常
- 名称匹配如有多个结果，会交互式提示选择
- 已存在的 skill 会自动跳过，不会重复安装
- 不要使用此工具执行 `config` 修改、`login`、`logout` 等敏感操作
