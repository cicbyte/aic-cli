package agent

import "path/filepath"

// AmpAgent 实现 Amp 的 AgentProfile
type AmpAgent struct{ projectAgent }

func (a *AmpAgent) Name() string { return "amp" }

func (a *AmpAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".amp")
}

func (a *AmpAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".amp", "skills")
}

func (a *AmpAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".config", "amp", "skills")
}
