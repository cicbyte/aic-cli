package agent

import (
	"path/filepath"
)

// joinPath 将 projectDir 和多个子路径片段拼接为绝对路径
func joinPath(projectDir string, parts ...string) string {
	all := append([]string{projectDir}, parts...)
	return filepath.Join(all...)
}

// AgentProfile 定义了一个 Coding Agent 的 skills 目录规范
type AgentProfile interface {
	// Name 返回 Agent 名称（如 "claude", "cursor"）
	Name() string

	// Detect 检测当前项目是否存在该 Agent 的目录/文件
	Detect(projectDir string) bool

	// SkillsDir 返回 skills 安装的绝对路径
	SkillsDir(projectDir string) string
}

