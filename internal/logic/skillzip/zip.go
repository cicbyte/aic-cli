package skillzip

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cicbyte/aic-cli/internal/models"
)

type ZipConfig struct {
	InputPath  string
	OutputPath string
	Format     string // "zip" (default) or "skill"
}

type ZipResult struct {
	OutputPath string
	FileCount  int
}

type ZipProcessor struct {
	config    *ZipConfig
	appConfig *models.AppConfig
}

func NewZipProcessor(config *ZipConfig, appConfig *models.AppConfig) *ZipProcessor {
	return &ZipProcessor{config: config, appConfig: appConfig}
}

func (p *ZipProcessor) Execute(ctx context.Context) (*ZipResult, error) {
	inputPath := p.config.InputPath

	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("目录不存在: %s", inputPath)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("路径不是目录: %s", inputPath)
	}

	skillMD := filepath.Join(inputPath, "skill.md")
	if _, err := os.Stat(skillMD); os.IsNotExist(err) {
		return nil, fmt.Errorf("目录中未找到 skill.md 文件")
	}

	outputPath := p.config.OutputPath
	if outputPath == "" {
		ext := ".zip"
		if p.config.Format == "skill" {
			ext = ".skill"
		}
		outputPath = filepath.Base(inputPath) + ext
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("创建 ZIP 文件失败: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	fileCount := 0
	err = filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := filepath.Base(path)
		if len(name) > 1 && name[0] == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(inputPath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			_, err = zipWriter.Create(relPath + "/")
			return err
		}

		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}

		fileCount++
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("打包失败: %w", err)
	}

	return &ZipResult{OutputPath: outputPath, FileCount: fileCount}, nil
}
