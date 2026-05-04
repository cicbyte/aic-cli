package skill

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/models"
)

type DownloadConfig struct {
	SkillID   int
	SkillName string
	OutputDir string
}

type DownloadResult struct {
	FilePath  string
	SkillName string
}

type DownloadProcessor struct {
	config    *DownloadConfig
	appConfig *models.AppConfig
}

func NewDownloadProcessor(config *DownloadConfig, appConfig *models.AppConfig) *DownloadProcessor {
	return &DownloadProcessor{config: config, appConfig: appConfig}
}

func (p *DownloadProcessor) Execute(ctx context.Context) (*DownloadResult, error) {
	client := api.NewClient(p.appConfig.AIC.BaseURL)

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
		outputDir = "."
	}

	outputPath := filepath.Join(outputDir, detail.Data.Name+".zip")

	savedPath, err := client.DownloadSkill(skillID, outputPath)
	if err != nil {
		return nil, fmt.Errorf("下载失败: %w", err)
	}

	return &DownloadResult{FilePath: savedPath, SkillName: detail.Data.Name}, nil
}
