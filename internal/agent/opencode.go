package agent

import "path/filepath"

// OpenCodeAgent 实现 OpenCode 的 AgentProfile
type OpenCodeAgent struct{ projectAgent }

func (a *OpenCodeAgent) Name() string { return "opencode" }

func (a *OpenCodeAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".opencode")
}

func (a *OpenCodeAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".opencode", "skills")
}

func (a *OpenCodeAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".config", "opencode", "skills")
}
