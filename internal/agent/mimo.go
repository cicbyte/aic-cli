package agent

import "path/filepath"

// MimoAgent 实现 MiMo Code 的 AgentProfile
// 文档来源：https://mimo.xiaomi.com/zh/mimocode/skills
// 项目级：.mimocode/skills/**/SKILL.md（单数 skill/ 亦可）
// 全局级：~/.config/mimocode/skills/**/SKILL.md
type MimoAgent struct{ projectAgent }

func (a *MimoAgent) Name() string { return "mimo" }

func (a *MimoAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".mimocode")
}

func (a *MimoAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".mimocode", "skills")
}

func (a *MimoAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".config", "mimocode", "skills")
}
