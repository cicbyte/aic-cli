package agent

import (
	"os"
	"path/filepath"
)

// CodeBuddyAgent 实现 CodeBuddy 的 AgentProfile（全局型）
type CodeBuddyAgent struct{ globalAgent }

func (a *CodeBuddyAgent) Name() string { return "codebuddy" }

func (a *CodeBuddyAgent) Detect(projectDir string) bool {
	home, _ := os.UserHomeDir()
	return detectDir(home, ".codebuddy")
}

func (a *CodeBuddyAgent) SkillsDir(projectDir string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".codebuddy", "skills")
}
