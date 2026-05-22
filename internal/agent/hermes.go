package agent

import "path/filepath"

// HermesAgent 实现 Hermes Agent 的 AgentProfile（全局型）
type HermesAgent struct{ globalOnlyAgent }

func (a *HermesAgent) Name() string { return "hermes" }

func (a *HermesAgent) Detect(projectDir string) bool {
	return detectDir(homeDir(), ".hermes")
}

func (a *HermesAgent) SkillsDir(projectDir string) string {
	return filepath.Join(homeDir(), ".hermes", "skills")
}

func (a *HermesAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".hermes", "skills")
}
