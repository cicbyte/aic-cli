package agent

import (
	"os"
	"path/filepath"
)

// HermesAgent 实现 Hermes Agent 的 AgentProfile（全局型）
type HermesAgent struct{ globalAgent }

func (a *HermesAgent) Name() string { return "hermes" }

func (a *HermesAgent) Detect(projectDir string) bool {
	home, _ := os.UserHomeDir()
	return detectDir(home, ".hermes")
}

func (a *HermesAgent) SkillsDir(projectDir string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".hermes", "skills")
}
