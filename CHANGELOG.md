# Changelog

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
