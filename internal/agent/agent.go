package agent

import (
	"os"
	"path/filepath"
)

// AgentProfile 定义了一个 Coding Agent 的 skills 目录规范
type AgentProfile interface {
	// Name 返回 Agent 名称（如 "claude", "cursor"）
	Name() string

	// HasProjectSkills 返回是否支持项目级 skills 目录
	HasProjectSkills() bool

	// Detect 检测当前项目是否存在该 Agent 的目录/文件
	Detect(projectDir string) bool

	// SkillsDir 返回项目级 skills 安装的绝对路径
	SkillsDir(projectDir string) string

	// GlobalSkillsDir 返回全局 skills 安装的绝对路径（如 ~/.claude/skills/）
	GlobalSkillsDir() string
}

// projectAgent 是项目级 Agent 的基础类型
type projectAgent struct{}

func (projectAgent) HasProjectSkills() bool { return true }

// globalOnlyAgent 是全局型 Agent 的基础类型（无项目级目录）
type globalOnlyAgent struct{}

func (globalOnlyAgent) HasProjectSkills() bool { return false }

// homeDir 返回用户主目录，失败时返回空字符串
func homeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// joinPath 将 projectDir 和多个子路径片段拼接为绝对路径
func joinPath(projectDir string, parts ...string) string {
	all := append([]string{projectDir}, parts...)
	return filepath.Join(all...)
}
