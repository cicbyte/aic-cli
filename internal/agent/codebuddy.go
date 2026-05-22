package agent

import "path/filepath"

// CodeBuddyAgent 实现 CodeBuddy 的 AgentProfile（全局型）
type CodeBuddyAgent struct{ globalOnlyAgent }

func (a *CodeBuddyAgent) Name() string { return "codebuddy" }

func (a *CodeBuddyAgent) Detect(projectDir string) bool {
	return detectDir(homeDir(), ".codebuddy")
}

func (a *CodeBuddyAgent) SkillsDir(projectDir string) string {
	return filepath.Join(homeDir(), ".codebuddy", "skills")
}

func (a *CodeBuddyAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".codebuddy", "skills")
}
