package agent

import (
	"os"
	"path/filepath"
)

// OpenClawAgent 实现 OpenClaw 的 AgentProfile（全局型）
type OpenClawAgent struct{ globalAgent }

func (a *OpenClawAgent) Name() string { return "openclaw" }

func (a *OpenClawAgent) Detect(projectDir string) bool {
	home, _ := os.UserHomeDir()
	return detectDir(home, ".openclaw")
}

func (a *OpenClawAgent) SkillsDir(projectDir string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".openclaw", "skills")
}
