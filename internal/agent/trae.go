package agent

// TraeAgent 实现 Trae 的 AgentProfile
type TraeAgent struct{}

func (a *TraeAgent) Name() string { return "trae" }

func (a *TraeAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".trae")
}

func (a *TraeAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".trae", "skills")
}
