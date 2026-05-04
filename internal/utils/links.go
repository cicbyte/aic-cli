package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/cicbyte/aic-cli/internal/models"
)

func GetLinksPath() string {
	return filepath.Join(ConfigInstance.GetAppDir(), "links.json")
}

func LoadLinks() *models.LinksConfig {
	linksPath := GetLinksPath()

	data, err := os.ReadFile(linksPath)
	if err != nil {
		return &models.LinksConfig{Links: []models.SkillLink{}}
	}

	var config models.LinksConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return &models.LinksConfig{Links: []models.SkillLink{}}
	}

	if config.Links == nil {
		config.Links = []models.SkillLink{}
	}

	return &config
}

func SaveLinks(config *models.LinksConfig) error {
	linksPath := GetLinksPath()

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(linksPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(linksPath, data, 0644)
}

func AddLink(name, globalPath, linkPath string) error {
	config := LoadLinks()

	// 移除已存在的同名记录
	for i, link := range config.Links {
		if link.Name == name {
			config.Links = append(config.Links[:i], config.Links[i+1:]...)
			break
		}
	}

	// 添加新记录
	config.Links = append(config.Links, models.SkillLink{
		Name:       name,
		GlobalPath: globalPath,
		LinkPath:   linkPath,
		AddedAt:    time.Now().Format(time.RFC3339),
	})

	return SaveLinks(config)
}

func RemoveLink(name string) error {
	config := LoadLinks()

	for i, link := range config.Links {
		if link.Name == name {
			config.Links = append(config.Links[:i], config.Links[i+1:]...)
			return SaveLinks(config)
		}
	}

	return nil
}

func GetLink(name string) *models.SkillLink {
	config := LoadLinks()

	for _, link := range config.Links {
		if link.Name == name {
			return &link
		}
	}

	return nil
}

func CleanBrokenLinks() []string {
	config := LoadLinks()
	var cleaned []string
	var valid []models.SkillLink

	for _, link := range config.Links {
		// 检查全局目录是否存在
		globalExists := DirExists(link.GlobalPath)
		// 检查软连接是否存在且有效
		linkValid := false
		if FileExists(link.LinkPath) || DirExists(link.LinkPath) {
			// 尝试读取目标，判断是否为有效链接
			_, err := os.Stat(link.LinkPath)
			linkValid = err == nil
		}

		if !globalExists || !linkValid {
			// 清理无效链接
			if linkValid {
				os.Remove(link.LinkPath)
			}
			cleaned = append(cleaned, link.Name)
		} else {
			valid = append(valid, link)
		}
	}

	config.Links = valid
	SaveLinks(config)

	return cleaned
}
