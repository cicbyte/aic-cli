package agent

import "path/filepath"

// ZCodeAgent 实现 ZCode 的 AgentProfile
// 文档来源：https://zcode.z.ai/cn/docs/skill
// 全局级（用户级）：~/.zcode/skills/<skill-name>/SKILL.md
// 项目级：.zcode/skills（文档提到支持项目级导入，路径按惯例推断）
type ZCodeAgent struct{ projectAgent }

func (a *ZCodeAgent) Name() string { return "zcode" }

func (a *ZCodeAgent) Detect(projectDir string) bool {
	return detectDir(projectDir, ".zcode")
}

func (a *ZCodeAgent) SkillsDir(projectDir string) string {
	return joinPath(projectDir, ".zcode", "skills")
}

func (a *ZCodeAgent) GlobalSkillsDir() string {
	return filepath.Join(homeDir(), ".zcode", "skills")
}
