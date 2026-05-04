package skill

import (
	"context"
	"fmt"
	"strings"

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
	client := api.NewClient(p.appConfig.AIC.BaseURL)

	resp, err := client.ListSkills(p.config.PageNum, p.config.PageSize, p.config.CategoryID)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	var filtered []api.Skill
	if p.config.Keyword != "" {
		keywordLower := strings.ToLower(p.config.Keyword)
		for _, skill := range resp.Data.List {
			if strings.Contains(strings.ToLower(skill.Name), keywordLower) ||
				strings.Contains(strings.ToLower(skill.Description), keywordLower) {
				filtered = append(filtered, skill)
			}
		}
	} else {
		filtered = resp.Data.List
	}

	return &SearchResult{Skills: filtered, Total: len(filtered)}, nil
}
