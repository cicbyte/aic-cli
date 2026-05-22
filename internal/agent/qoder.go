package agent

import "path/filepath"

// QoderAgent 实现 Qoder 的 AgentProfile（同时支持项目级和全局）
type QoderAgent struct{ projectAgent }

func (a *QoderAgent) Name() string { return "qoder" }

func (a *QoderAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".qoder")
}

func (a *QoderAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".qoder", "skills")
}

func (a *QoderAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".qoder", "skills")
}
