package local

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/cicbyte/aic-cli/internal/models"
	"github.com/cicbyte/aic-cli/internal/utils"
)

type SkillEntry struct {
	Name   string
	IsDir  bool
	IsLink bool
}

type ListConfig struct {
	Global     bool
	WorkingDir string
}

type ListResult struct {
	Skills []SkillEntry
	Dir    string
}

type ListProcessor struct {
	config    *ListConfig
	appConfig *models.AppConfig
}

func NewListProcessor(config *ListConfig, appConfig *models.AppConfig) *ListProcessor {
	return &ListProcessor{config: config, appConfig: appConfig}
}

func (p *ListProcessor) Execute(ctx context.Context) (*ListResult, error) {
	var skillsDir string

	if p.config.Global {
		skillsDir = utils.GetGlobalSkillsDir()
	} else {
		if p.config.WorkingDir == "" {
			wd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			p.config.WorkingDir = wd
		}
		skillsDir = filepath.Join(p.config.WorkingDir, ".claude", "skills")
	}

	if !utils.DirExists(skillsDir) {
		return &ListResult{Dir: skillsDir}, nil
	}

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}

	var skills []SkillEntry
	for _, entry := range entries {
		info, _ := entry.Info()
		isLink := info != nil && info.Mode()&os.ModeSymlink != 0
		skills = append(skills, SkillEntry{
			Name:   entry.Name(),
			IsDir:  entry.IsDir(),
			IsLink: isLink,
		})
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	return &ListResult{Skills: skills, Dir: skillsDir}, nil
}
