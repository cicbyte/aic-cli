package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cicbyte/aic-cli/internal/models"
	"github.com/cicbyte/aic-cli/internal/utils"
)

type ImportConfig struct {
	InputPath   string
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
	inputPath := p.config.InputPath
	uploadPath := inputPath
	var tempZip string

	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("路径不存在: %s", inputPath)
	}

	if info.IsDir() {
		skillMD := filepath.Join(inputPath, "skill.md")
		if _, err := os.Stat(skillMD); os.IsNotExist(err) {
			return nil, fmt.Errorf("目录中未找到 skill.md 文件: %s", inputPath)
		}

		tmpFile, err := os.CreateTemp("", "aic-import-*.zip")
		if err != nil {
			return nil, fmt.Errorf("创建临时文件失败: %w", err)
		}
		tmpFile.Close()
		tempZip = tmpFile.Name()

		if err := utils.ZipDir(inputPath, tempZip); err != nil {
			os.Remove(tempZip)
			return nil, fmt.Errorf("压缩目录失败: %w", err)
		}
		uploadPath = tempZip
	} else {
		ext := strings.ToLower(filepath.Ext(inputPath))
		if ext != ".zip" && ext != ".skill" {
			return nil, fmt.Errorf("不支持的文件格式: %s，支持 .zip、.skill 或文件夹", ext)
		}
	}

	defer func() {
		if tempZip != "" {
			os.Remove(tempZip)
		}
	}()

	client := NewClient(p.appConfig)

	resp, err := client.ImportZip(uploadPath, p.config.Description, p.config.CategoryID, p.config.Overwrite)
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
