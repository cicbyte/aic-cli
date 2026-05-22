package agent

import "path/filepath"

// CodexAgent 实现 Codex CLI 的 AgentProfile
type CodexAgent struct{ projectAgent }

func (a *CodexAgent) Name() string { return "codex" }

func (a *CodexAgent) Detect(projectDir string) bool {
	return detectDirOrFile(projectDir, ".codex") || detectFile(projectDir, "AGENTS.md")
}

func (a *CodexAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".codex", "skills")
}

func (a *CodexAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".codex", "skills")
}
