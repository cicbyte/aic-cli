# Changelog

## [0.0.4] - 2026-05-28

### Features

- 新增 `skill create` 命令，创建新技能并上传到服务器，支持交互式和命令行两种模式
- 新增 `skill remote` 子命令组（AI Agent 专用），与常用命令隔离
  - `skill remote cat` - 查看远程文件内容，支持保存到本地
  - `skill remote tree` - 显示远程技能文件树
  - `skill remote patch` - 增量编辑文件，支持纯内容匹配、行号替换、混合模式
  - `skill remote edit` - 交互式编辑远程文件
  - `skill remote validate` - 校验技能文件，支持严格模式
  - `skill remote publish` - 发布技能，支持指定版本和说明
  - `skill remote unpublish` - 取消发布
  - `skill remote update` - 更新技能元数据（名称、描述、标签、分类）
- 新增 `patch --schema` 命令，输出 JSON 格式说明，方便 AI Agent 使用
- 新增管道输入自动检测
- 新增 `batch` 模式，支持批量文件操作

### Testing

- 新增 pytest 集成测试框架
- 新增 122 个测试用例，覆盖所有命令
- 新增测试辅助模块：CLI 客户端封装、通用断言
- 新增测试数据管理：`tests/data/test_data.yaml`

### Docs

- 新增 `aic-author` skill，创建者视角的使用指南
- 更新 `aic-cli` skill，专注于使用者视角
- 更新中英文 README，补充新功能说明
- 新增 `tests/README.md`，测试框架使用说明
- 移除文档中的 emoji

## [0.0.3] - 2026-05-28

### Fixes

- 修复 `skill import` 在 Windows 上 zip 内路径使用反斜杠导致服务端目录结构丢失的问题
- 修复 `server login` 输入带尾部斜杠的 URL 后请求地址异常的问题

### Docs

- 更新中英文 README，补充 AIC 仓库关联说明

## [0.0.2] - 2026-05-27

### Fixes

- 修复 `gen_release_notes.py` 在 Windows 中文系统下的 GBK 编码问题

### CI

- 重构 release 相关脚本代码
- 优化 GitHub Release 工作流，修正模块路径和输出文件名

### Docs

- 更新项目说明文档，补充 AIC CLI 介绍

### Fixes

- 修复 `skill import` 在 Windows 上 zip 内路径使用反斜杠导致服务端目录结构丢失的问题
- 修复 `server login` 输入带尾部斜杠的 URL 后请求地址异常的问题

### CI

- 优化 GitHub Release 工作流，修正模块路径和输出文件名
- 修复 `gen_release_notes.py` 在 Windows 中文系统下的编码问题

### Docs

- 更新中英文 README，补充 AIC 仓库关联说明
- 维护 CHANGELOG

## [0.0.1] - 2026-05-27

### Features

- 初始化 AIC CLI 项目，实现基础 skill 管理功能（搜索、下载、安装）
- 新增多 Agent 支持，覆盖 18 种 Coding Agent（Claude Code、Cursor、Windsurf、Cline、Continue、OpenCode、Codex、Gemini、Roo、Amp、Trae、Qoder、OpenClaw、Amazon Q、Copilot、Hermes、CodeBuddy、QClaw）
- 新增 `skill add` 命令，支持按 ID 或名称搜索，交互式选择
- 新增 `skill search` 命令，支持关键词和分类筛选
- 新增 `skill download` 命令，下载 skill zip 包
- 新增 `skill import` 命令，支持导入 zip、.skill 格式和文件夹
- 新增 `skill list` 命令，列出已安装的技能
- 新增 `skill remove` 命令，移除已安装的技能
- 新增 `skill clean` 命令，清理失效的软链接
- 新增 `skill pack` 命令，将技能文件夹打包为 .zip 或 .skill
- 新增 `skill mode` 命令，查看或切换安装模式（symlink / copy）
- 新增 `package list` 命令，列出和搜索技能包
- 新增 `package add` 命令，批量安装技能包中的所有技能
- 新增 `server open` 命令，在浏览器中打开 AIC 服务页面
- 新增 `server status` 命令，查看服务器连接状态
- 新增 `server login` / `server logout` 命令，管理认证
- 新增 `config list` 命令，查看当前配置
- 新增版本管理，支持 `--version` 参数
- 支持 symlink 和 copy 两种安装模式，symlink 为默认模式
- 支持 `--agent` 参数指定目标 Agent，自动检测当前项目的 Agent
- 支持项目级和全局 skills 目录
- Windows 平台使用 Junction 创建软链接，无需管理员权限
- 新增 aic-cli 技能文档，指导 AI 使用本工具

### Refactors

- 重构项目结构为三层架构：cmd → internal/logic → internal/api + internal/utils
- 重构 Agent 体系，使用 `AgentProfile` 接口统一管理，区分项目级和全局型 Agent
- 移除交互式 TUI 功能

### Docs

- 更新中英文 README，适配多 Agent 架构
- 添加运行截图（skill add、package add）
