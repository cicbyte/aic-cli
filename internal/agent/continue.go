package agent

// ContinueAgent 实现 Continue 的 AgentProfile
type ContinueAgent struct{}

func (a *ContinueAgent) Name() string { return "continue" }

func (a *ContinueAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".continue")
}

func (a *ContinueAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".continue", "rules")
}
