package agent

import "path/filepath"

// AmazonQAgent 实现 Amazon Q 的 AgentProfile
type AmazonQAgent struct{ projectAgent }

func (a *AmazonQAgent) Name() string { return "amazonq" }

func (a *AmazonQAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".amazonq")
}

func (a *AmazonQAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".amazonq", "rules")
}

func (a *AmazonQAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".amazonq", "rules")
}
