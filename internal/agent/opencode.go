package agent

// OpenCodeAgent 实现 OpenCode 的 AgentProfile
type OpenCodeAgent struct{}

func (a *OpenCodeAgent) Name() string { return "opencode" }

func (a *OpenCodeAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".opencode")
}

func (a *OpenCodeAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".opencode", "skills")
}
