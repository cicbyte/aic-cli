package agent

// WindsurfAgent 实现 Windsurf 的 AgentProfile
// Windsurf 原为单文件型（.windsurfrules），统一目录结构后使用 .windsurf/skills/
type WindsurfAgent struct{}

func (a *WindsurfAgent) Name() string { return "windsurf" }

func (a *WindsurfAgent) Detect(projectDir string) bool {
	return detectFile(projectDir, ".windsurfrules")
}

func (a *WindsurfAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".windsurf", "skills")
}
