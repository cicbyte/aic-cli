package agent

// GeminiAgent 实现 Gemini CLI 的 AgentProfile
type GeminiAgent struct{}

func (a *GeminiAgent) Name() string { return "gemini" }

func (a *GeminiAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".gemini")
}

func (a *GeminiAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".gemini", "skills")
}
