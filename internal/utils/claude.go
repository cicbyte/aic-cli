package utils

import (
	"os"
	"path/filepath"

	"github.com/cicbyte/aic-cli/internal/agent"
	"github.com/cicbyte/aic-cli/internal/common"
)

func IsClaudeCodeProject(dir string) bool {
	claudeDir := filepath.Join(dir, ".claude")
	info, err := os.Stat(claudeDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func GetGlobalSkillsDir() string {
	if common.AppConfigModel != nil && common.AppConfigModel.Skills.GlobalDir != "" {
		absPath, err := filepath.Abs(common.AppConfigModel.Skills.GlobalDir)
		if err == nil {
			return absPath
		}
	}
	return ConfigInstance.GetGlobalSkillsDir()
}

// GetSkillsOutputDir 保持向后兼容，默认使用 Claude Code
func GetSkillsOutputDir(outputDir string) (string, error) {
	if outputDir != "" {
		absPath, err := filepath.Abs(outputDir)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if IsClaudeCodeProject(cwd) {
		skillsDir := filepath.Join(cwd, ".claude", "skills")
		return skillsDir, nil
	}

	return "", nil
}

// GetSkillsDirForAgent 获取指定 Agent 的 skills 目录
func GetSkillsDirForAgent(a agent.AgentProfile, projectDir string) (string, error) {
	skillsDir := a.SkillsDir(projectDir)
	return skillsDir, nil
}

// GetDefaultAgentName 从配置中获取默认 Agent 名称
func GetDefaultAgentName() string {
	if common.AppConfigModel != nil && common.AppConfigModel.Skills.DefaultAgent != "" {
		return common.AppConfigModel.Skills.DefaultAgent
	}
	return "claude"
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
