package agent

import "path/filepath"

// ClineAgent 实现 Cline 的 AgentProfile
// Cline 原为单文件型（.clinerules），统一目录结构后使用 .cline/skills/
type ClineAgent struct{ projectAgent }

func (a *ClineAgent) Name() string { return "cline" }

func (a *ClineAgent) Detect(projectDir string) bool {
	return detectFile(projectDir, ".clinerules")
}

func (a *ClineAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".cline", "skills")
}

func (a *ClineAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".cline", "skills")
}
