package agent

import "path/filepath"

// RooAgent 实现 Roo Code 的 AgentProfile
type RooAgent struct{ projectAgent }

func (a *RooAgent) Name() string { return "roo" }

func (a *RooAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".roo")
}

func (a *RooAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".roo", "skills")
}

func (a *RooAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".roo", "skills")
}
