package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cicbyte/aic-cli/internal/models"
	"github.com/cicbyte/aic-cli/internal/utils"
)

type AddConfig struct {
	SkillID   int
	SkillName string
	OutputDir string // 已解析的绝对路径（由命令层传入）
	Mode      string
	Overwrite bool
}

type AddResult struct {
	SkillName   string
	InstallPath string
}

type AddProcessor struct {
	config    *AddConfig
	appConfig *models.AppConfig
}

func NewAddProcessor(config *AddConfig, appConfig *models.AppConfig) *AddProcessor {
	return &AddProcessor{config: config, appConfig: appConfig}
}

func (p *AddProcessor) Execute(ctx context.Context) (*AddResult, error) {
	client := NewClient(p.appConfig)

	skillID, err := ResolveSkillID(client, p.config.SkillID, p.config.SkillName)
	if err != nil {
		return nil, err
	}

	detail, err := client.GetSkillDetail(skillID)
	if err != nil {
		return nil, fmt.Errorf("获取 skill 信息失败: %w", err)
	}
	if detail.Code != 0 {
		return nil, fmt.Errorf("%s", detail.Message)
	}

	outputDir := p.config.OutputDir
	if outputDir == "" {
		return nil, fmt.Errorf("未指定输出目录")
	}

	tmpFile, err := os.CreateTemp("", "skill-*.zip")
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if _, err := client.DownloadSkill(skillID, tmpPath); err != nil {
		return nil, fmt.Errorf("下载失败: %w", err)
	}

	if p.config.Mode == "symlink" {
		return p.installSymlink(detail.Data.Name, tmpPath, outputDir)
	}
	return p.installCopy(detail.Data.Name, tmpPath, outputDir)
}

func (p *AddProcessor) installSymlink(name, zipPath, outputDir string) (*AddResult, error) {
	globalSkillsDir := utils.GetGlobalSkillsDir()
	globalSkillDir := filepath.Join(globalSkillsDir, name)

	if !utils.DirExists(globalSkillDir) {
		if err := utils.Unzip(zipPath, globalSkillDir); err != nil {
			return nil, fmt.Errorf("解压失败: %w", err)
		}
	}

	symlinkPath := filepath.Join(outputDir, name)
	if utils.FileExists(symlinkPath) || utils.DirExists(symlinkPath) {
		if !p.config.Overwrite {
			return nil, fmt.Errorf("目标目录已存在: %s（需要确认覆盖）", symlinkPath)
		}
		os.Remove(symlinkPath)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目标目录失败: %w", err)
	}

	if err := utils.CreateLink(globalSkillDir, symlinkPath); err != nil {
		return nil, fmt.Errorf("创建连接失败: %w", err)
	}

	_ = utils.AddLink(name, globalSkillDir, symlinkPath)

	return &AddResult{SkillName: name, InstallPath: symlinkPath}, nil
}

func (p *AddProcessor) installCopy(name, zipPath, outputDir string) (*AddResult, error) {
	skillDir := filepath.Join(outputDir, name)

	if utils.DirExists(skillDir) {
		if !p.config.Overwrite {
			return nil, fmt.Errorf("skill 目录已存在: %s（需要确认覆盖）", skillDir)
		}
		os.RemoveAll(skillDir)
	}

	if err := utils.Unzip(zipPath, skillDir); err != nil {
		return nil, fmt.Errorf("解压失败: %w", err)
	}

	return &AddResult{SkillName: name, InstallPath: skillDir}, nil
}
