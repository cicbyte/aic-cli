package skill

import (
	"context"
	"fmt"

	"github.com/cicbyte/aic-cli/internal/models"
)

type ImportConfig struct {
	ZipPath     string
	Description string
	CategoryID  int
	Overwrite   bool
}

type ImportResult struct {
	SkillName string
	SkillID   int
}

type ImportProcessor struct {
	config    *ImportConfig
	appConfig *models.AppConfig
}

func NewImportProcessor(config *ImportConfig, appConfig *models.AppConfig) *ImportProcessor {
	return &ImportProcessor{config: config, appConfig: appConfig}
}

func (p *ImportProcessor) Execute(ctx context.Context) (*ImportResult, error) {
	client := NewClient(p.appConfig)

	resp, err := client.ImportZip(p.config.ZipPath, p.config.Description, p.config.CategoryID, p.config.Overwrite)
	if err != nil {
		return nil, fmt.Errorf("导入失败: %w", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	return &ImportResult{
		SkillName: resp.Data.Name,
		SkillID:   resp.Data.SkillID,
	}, nil
}
