package agent

import (
	"os"
	"path/filepath"
)

// QClawAgent 实现 QClaw 的 AgentProfile（全局型）
type QClawAgent struct{ globalAgent }

func (a *QClawAgent) Name() string { return "qclaw" }

func (a *QClawAgent) Detect(projectDir string) bool {
	home, _ := os.UserHomeDir()
	return detectDir(home, ".qclaw")
}

func (a *QClawAgent) SkillsDir(projectDir string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".qclaw", "skills")
}
