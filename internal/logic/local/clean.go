package local

import (
	"context"

	"github.com/cicbyte/aic-cli/internal/models"
	"github.com/cicbyte/aic-cli/internal/utils"
)

type CleanConfig struct{}

type CleanResult struct {
	Cleaned []string
}

type CleanProcessor struct {
	config    *CleanConfig
	appConfig *models.AppConfig
}

func NewCleanProcessor(config *CleanConfig, appConfig *models.AppConfig) *CleanProcessor {
	return &CleanProcessor{config: config, appConfig: appConfig}
}

func (p *CleanProcessor) Execute(ctx context.Context) (*CleanResult, error) {
	cleaned := utils.CleanBrokenLinks()
	return &CleanResult{Cleaned: cleaned}, nil
}
