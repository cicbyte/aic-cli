package utils

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/cicbyte/aic-cli/internal/agent"
)

// ResolveOutputDir 根据 --agent 和 --outputDir 解析最终输出目录
// 优先级: --outputDir > --agent > config default_agent > 检测 > 默认 claude
func ResolveOutputDir(outputDir string, agentName string) (string, error) {
	if outputDir != "" {
		return GetSkillsOutputDir(outputDir)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	a, err := SelectAgent(cwd, agentName)
	if err != nil {
		return "", err
	}

	return a.SkillsDir(cwd), nil
}

// SelectAgent 选择目标 Agent
func SelectAgent(projectDir string, agentName string) (agent.AgentProfile, error) {
	if agentName != "" {
		a, err := agent.GetAgent(agentName)
		if err != nil {
			return nil, fmt.Errorf("未知的 Agent: %s，支持: claude, cursor, continue, amazonq, copilot, windsurf, cline, opencode, codex, gemini, roo, amp", agentName)
		}
		return a, nil
	}

	defaultName := GetDefaultAgentName()
	if defaultName != "claude" {
		a, err := agent.GetAgent(defaultName)
		if err == nil {
			return a, nil
		}
	}

	detected := agent.DetectAgents(projectDir)

	switch len(detected) {
	case 0:
		return agent.GetAgent("claude")
	case 1:
		return detected[0], nil
	default:
		return PromptSelectAgent(detected)
	}
}

// PromptSelectAgent 交互式选择 Agent
func PromptSelectAgent(agents []agent.AgentProfile) (agent.AgentProfile, error) {
	options := make([]huh.Option[agent.AgentProfile], len(agents))
	for i, a := range agents {
		options[i] = huh.NewOption(a.Name(), a)
	}

	var selected agent.AgentProfile
	err := huh.NewSelect[agent.AgentProfile]().
		Title("检测到多个 Coding Agent，请选择目标:").
		Options(options...).
		Value(&selected).
		Run()
	if err != nil {
		return nil, fmt.Errorf("选择已取消")
	}
	return selected, nil
}
