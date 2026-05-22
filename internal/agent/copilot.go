package agent

// CopilotAgent 实现 GitHub Copilot 的 AgentProfile
type CopilotAgent struct{}

func (a *CopilotAgent) Name() string { return "copilot" }

func (a *CopilotAgent) Detect(projectDir string) bool {
	return detectFile(projectDir, ".github/copilot-instructions.md")
}

func (a *CopilotAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".github", "prompts")
}
