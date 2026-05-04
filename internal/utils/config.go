package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/cicbyte/aic-cli/internal/models"
	"go.yaml.in/yaml/v3"
)

var ConfigInstance = Config{}

type Config struct {
	HomeDir         string
	AppSeriesDir    string
	AppDir          string
	ConfigDir       string
	ConfigPath      string
	GlobalSkillsDir string
	LogDir          string
	LogPath         string
}

func (c *Config) GetHomeDir() string {
	if c.HomeDir != "" {
		return c.HomeDir
	}
	usr, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current user: %v", err))
	}
	c.HomeDir = usr.HomeDir
	return c.HomeDir
}

func (c *Config) GetAppSeriesDir() string {
	if c.AppSeriesDir != "" {
		return c.AppSeriesDir
	}
	c.AppSeriesDir = c.GetHomeDir() + "/.ciclebyte"
	return c.AppSeriesDir
}

func (c *Config) GetAppDir() string {
	if c.AppDir != "" {
		return c.AppDir
	}
	c.AppDir = c.GetAppSeriesDir() + "/aic-cli"
	return c.AppDir
}

func (c *Config) GetConfigDir() string {
	if c.ConfigDir != "" {
		return c.ConfigDir
	}
	c.ConfigDir = c.GetAppDir() + "/config"
	return c.ConfigDir
}

func (c *Config) GetConfigPath() string {
	if c.ConfigPath != "" {
		return c.ConfigPath
	}
	c.ConfigPath = c.GetConfigDir() + "/config.yaml"
	return c.ConfigPath
}

func (c *Config) GetGlobalSkillsDir() string {
	if c.GlobalSkillsDir != "" {
		return c.GlobalSkillsDir
	}
	c.GlobalSkillsDir = filepath.Join(c.GetAppDir(), "skills")
	return c.GlobalSkillsDir
}

func (c *Config) GetLogDir() string {
	if c.LogDir == "" {
		c.LogDir = filepath.Join(c.GetAppDir(), "logs")
	}
	return c.LogDir
}

func (c *Config) GetLogPath() string {
	if c.LogPath == "" {
		now := time.Now().Format("20060102")
		c.LogPath = filepath.Join(c.GetLogDir(), fmt.Sprintf("aic-cli_log_%s.log", now))
	}
	return c.LogPath
}

func (c *Config) LoadConfig() *models.AppConfig {
	configPath := c.GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := GetDefaultConfig()
		data, err := yaml.Marshal(defaultConfig)
		if err == nil {
			_ = os.WriteFile(configPath, data, 0644)
		}
		return defaultConfig
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return GetDefaultConfig()
	}

	var config models.AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return GetDefaultConfig()
	}

	return &config
}

func (c *Config) SaveConfig(config *models.AppConfig) {
	configPath := c.GetConfigPath()
	data, err := yaml.Marshal(config)
	if err != nil {
		return
	}
	os.WriteFile(configPath, data, 0644)
}

func GetDefaultConfig() *models.AppConfig {
	config := &models.AppConfig{}

	config.AIC.BaseURL = "http://localhost:8000"

	config.Skills.GlobalDir = ""       // 留空表示使用默认目录 ~/.ciclebyte/aic-cli/skills
	config.Skills.DefaultMode = "copy" // 默认安装模式: copy 或 symlink

	config.Log.Level = "info"
	config.Log.MaxSize = 10
	config.Log.MaxBackups = 30
	config.Log.MaxAge = 30
	config.Log.Compress = true

	return config
}
