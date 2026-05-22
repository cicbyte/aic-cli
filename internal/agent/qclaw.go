package agent

import "path/filepath"

// QClawAgent 实现 QClaw 的 AgentProfile（全局型）
type QClawAgent struct{ globalOnlyAgent }

func (a *QClawAgent) Name() string { return "qclaw" }

func (a *QClawAgent) Detect(projectDir string) bool {
	return detectDir(homeDir(), ".qclaw")
}

func (a *QClawAgent) SkillsDir(projectDir string) string {
	return filepath.Join(homeDir(), ".qclaw", "skills")
}

func (a *QClawAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".qclaw", "skills")
}
