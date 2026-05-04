package skill

import (
	"context"
	"fmt"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/models"
)

type SearchConfig struct {
	Keyword    string
	CategoryID int
	PageNum    int
	PageSize   int
}

type SearchResult struct {
	Skills []api.Skill
	Total  int
}

type SearchProcessor struct {
	config    *SearchConfig
	appConfig *models.AppConfig
}

func NewSearchProcessor(config *SearchConfig, appConfig *models.AppConfig) *SearchProcessor {
	return &SearchProcessor{config: config, appConfig: appConfig}
}

func (p *SearchProcessor) Execute(ctx context.Context) (*SearchResult, error) {
	client := NewClient(p.appConfig)

	resp, err := client.ListSkills(p.config.PageNum, p.config.PageSize, p.config.CategoryID, p.config.Keyword)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	return &SearchResult{Skills: resp.Data.List, Total: resp.Data.Total}, nil
}
