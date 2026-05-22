package agent

// CursorAgent 实现 Cursor 的 AgentProfile
type CursorAgent struct{ projectAgent }

func (a *CursorAgent) Name() string { return "cursor" }

func (a *CursorAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".cursor")
}

func (a *CursorAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".cursor", "rules")
}
