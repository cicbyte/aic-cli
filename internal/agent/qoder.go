package agent

import (
	"os"
	"path/filepath"
)

// QoderAgent 实现 Qoder 的 AgentProfile（全局型）
type QoderAgent struct{ globalAgent }

func (a *QoderAgent) Name() string { return "qoder" }

func (a *QoderAgent) Detect(projectDir string) bool {
	home, _ := os.UserHomeDir()
	return detectDir(home, ".qoder")
}

func (a *QoderAgent) SkillsDir(projectDir string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".qoder", "skills")
}
