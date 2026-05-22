package agent

import "path/filepath"

// OpenClawAgent 实现 OpenClaw 的 AgentProfile（同时支持项目级和全局）
type OpenClawAgent struct{ projectAgent }

func (a *OpenClawAgent) Name() string { return "openclaw" }

func (a *OpenClawAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".openclaw")
}

func (a *OpenClawAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".openclaw", "skills")
}

func (a *OpenClawAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".openclaw", "skills")
}
