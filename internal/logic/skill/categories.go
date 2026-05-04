package skill

import (
	"context"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/models"
)

type CategoriesConfig struct{}

type CategoriesResult struct {
	Categories []api.Category
}

type CategoriesProcessor struct {
	config    *CategoriesConfig
	appConfig *models.AppConfig
}

func NewCategoriesProcessor(config *CategoriesConfig, appConfig *models.AppConfig) *CategoriesProcessor {
	return &CategoriesProcessor{config: config, appConfig: appConfig}
}

func (p *CategoriesProcessor) Execute(ctx context.Context) (*CategoriesResult, error) {
	client := api.NewClient(p.appConfig.AIC.BaseURL)
	categories, err := client.ListCategories()
	if err != nil {
		return nil, err
	}
	return &CategoriesResult{Categories: categories}, nil
}
