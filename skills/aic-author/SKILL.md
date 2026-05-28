---
name: aic-author
description: |
  AIC 技能创建与维护工具。用于创建、编辑、校验、发布技能。
  当用户需要创建新技能、编辑技能文件、校验技能或发布技能时使用。
  触发词："创建skill"、"新建技能"、"编辑技能"、"发布技能"、"校验技能"、"skill模板"。
---

# AIC Author 使用指南

用于创建和维护 AIC 技能的命令行工具。

## 命令总览

```
aic-cli skill create    # 创建新技能
aic-cli skill remote    # 远程技能管理（文件操作、校验、发布）
```

## 创建新技能

```bash
# 交互式创建（推荐）
aic-cli skill create

# 指定名称创建
aic-cli skill create my-skill

# 完整选项
aic-cli skill create my-skill --category 1 --tags "go,cli" --desc "我的技能"
```

创建后会自动生成 SKILL.md 模板并上传到服务器。

## 文件操作

### 查看文件内容

```bash
aic-cli skill remote cat <ID> SKILL.md
aic-cli skill remote cat <ID> SKILL.md -o local.md  # 保存到本地
```

### 查看文件树

```bash
aic-cli skill remote tree <ID>
```

### 增量编辑

```bash
# 内容匹配
aic-cli skill remote patch <ID> --path SKILL.md \
  --old "旧内容" --new "新内容"

# 行号 + 内容校验（推荐）
aic-cli skill remote patch <ID> --path SKILL.md \
  --line 5 --old "旧内容" --new "新内容"

# 管道 JSON
echo '{"path":"SKILL.md","edits":[...]}' | aic-cli skill remote patch <ID>

# 批量操作
aic-cli skill remote patch <ID> --batch edits.json

# 查看 JSON 格式
aic-cli skill remote patch --schema
```

### 交互式编辑

```bash
aic-cli skill remote edit <ID>              # 选择文件
aic-cli skill remote edit <ID> SKILL.md     # 指定文件
aic-cli skill remote edit <ID> --dry-run    # 只看差异
```

## 校验与发布

### 校验技能

```bash
aic-cli skill remote validate <ID>          # 标准校验
aic-cli skill remote validate <ID> --strict # 严格模式
```

### 发布技能

```bash
aic-cli skill remote publish <ID>                     # 发布
aic-cli skill remote publish <ID> --version 1.0.0    # 指定版本
aic-cli skill remote publish <ID> --changelog "说明"  # 带说明
```

### 取消发布

```bash
aic-cli skill remote unpublish <ID>
```

## 元数据管理

```bash
aic-cli skill remote update <ID> --name "新名称"
aic-cli skill remote update <ID> --desc "新描述"
aic-cli skill remote update <ID> --tags "tag1,tag2"
aic-cli skill remote update <ID> --category 2
```

## 典型工作流

### 1. 创建新技能

```bash
# 创建技能
aic-cli skill create my-awesome-skill --desc "一个很棒的技能"
# 输出: ✓ 技能已创建 (ID: 42)

# 查看生成的模板
aic-cli skill remote cat 42 SKILL.md
```

### 2. 编辑技能内容

```bash
# 编辑 SKILL.md
aic-cli skill remote patch 42 --path SKILL.md \
  --old "[TODO: 填写技能描述]" \
  --new "这是一个自动化工具技能"

# 添加新文件
echo '{"path":"scripts/helper.py","content":"#!/usr/bin/python3\nprint(\"hello\")"}' \
  | aic-cli skill remote patch 42
```

### 3. 校验并发布

```bash
# 校验
aic-cli skill remote validate 42

# 发布
aic-cli skill remote publish 42 --version 1.0.0 --changelog "首次发布"
```

## 注意事项

- 技能名称使用 kebab-case（如 `my-skill`）
- SKILL.md 必须包含 YAML frontmatter（name 和 description）
- 发布前会自动校验，校验失败会拒绝发布
- 文件路径禁止 `..`、绝对路径、隐藏文件
