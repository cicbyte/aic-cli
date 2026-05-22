package agent

import "path/filepath"

// ClaudeAgent 实现 Claude Code 的 AgentProfile
type ClaudeAgent struct{ projectAgent }

func (a *ClaudeAgent) Name() string { return "claude" }

func (a *ClaudeAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".claude")
}

func (a *ClaudeAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".claude", "skills")
}

func (a *ClaudeAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".claude", "skills")
}
