package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cicbyte/aic-cli/internal/models"
	"github.com/cicbyte/aic-cli/internal/utils"
)

type RemoveConfig struct {
	SkillName  string
	Global     bool
	WorkingDir string
}

type RemoveResult struct {
	SkillName     string
	RemovedGlobal bool
}

type RemoveProcessor struct {
	config    *RemoveConfig
	appConfig *models.AppConfig
}

func NewRemoveProcessor(config *RemoveConfig, appConfig *models.AppConfig) *RemoveProcessor {
	return &RemoveProcessor{config: config, appConfig: appConfig}
}

func (p *RemoveProcessor) Execute(ctx context.Context) (*RemoveResult, error) {
	if p.config.WorkingDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		p.config.WorkingDir = wd
	}

	skillsDir := filepath.Join(p.config.WorkingDir, ".claude", "skills")
	skillPath := filepath.Join(skillsDir, p.config.SkillName)

	if !utils.FileExists(skillPath) && !utils.DirExists(skillPath) {
		return nil, fmt.Errorf("skill '%s' 不存在", p.config.SkillName)
	}

	if err := os.Remove(skillPath); err != nil {
		return nil, fmt.Errorf("删除失败: %w", err)
	}

	utils.RemoveLink(p.config.SkillName)

	result := &RemoveResult{SkillName: p.config.SkillName}

	if p.config.Global {
		globalDir := utils.GetGlobalSkillsDir()
		globalSkillDir := filepath.Join(globalDir, p.config.SkillName)
		if utils.DirExists(globalSkillDir) {
			os.RemoveAll(globalSkillDir)
			result.RemovedGlobal = true
		}
	}

	return result, nil
}
