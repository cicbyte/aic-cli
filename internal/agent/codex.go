package agent

// CodexAgent 实现 Codex CLI 的 AgentProfile
type CodexAgent struct{}

func (a *CodexAgent) Name() string { return "codex" }

func (a *CodexAgent) Detect(projectDir string) bool {
	return detectDirOrFile(projectDir, ".codex") || detectFile(projectDir, "AGENTS.md")
}

func (a *CodexAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".codex", "skills")
}
