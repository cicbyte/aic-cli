package utils

import (
	"os"
	"path/filepath"

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
