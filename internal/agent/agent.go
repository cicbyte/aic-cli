package agent

import (
	"path/filepath"
)

// projectAgent 是项目级 Agent 的基础类型（IsGlobal 返回 false）
type projectAgent struct{}

func (projectAgent) IsGlobal() bool { return false }

// globalAgent 是全局型 Agent 的基础类型（IsGlobal 返回 true）
type globalAgent struct{}

func (globalAgent) IsGlobal() bool { return true }

// joinPath 将 projectDir 和多个子路径片段拼接为绝对路径
func joinPath(projectDir string, parts ...string) string {
	all := append([]string{projectDir}, parts...)
	return filepath.Join(all...)
}

// AgentProfile 定义了一个 Coding Agent 的 skills 目录规范
type AgentProfile interface {
	// Name 返回 Agent 名称（如 "claude", "cursor"）
	Name() string

	// IsGlobal 返回是否为全局型 Agent（只有全局 skills 目录，无项目级目录）
	IsGlobal() bool

	// Detect 检测当前项目是否存在该 Agent 的目录/文件
	// 对于全局型 Agent，检测用户主目录下是否存在该 Agent
	Detect(projectDir string) bool

	// SkillsDir 返回 skills 安装的绝对路径
	// 对于项目型 Agent，基于 projectDir 计算
	// 对于全局型 Agent，返回用户主目录下的固定路径
	SkillsDir(projectDir string) string
}

