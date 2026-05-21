package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/models"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setCmd = &cobra.Command{
	Use:   "set <key> [value]",
	Short: "设置配置项的值",
	Long: `设置指定配置项的值。

敏感字段（如 token）如果不提供 value 参数，会以不回显方式交互式输入。

示例:
  aic-cli config set aic.base_url http://192.168.1.100:8000
  aic-cli config set aic.token your-token-here
  aic-cli config set aic.token                    # 交互式输入（不回显）
  aic-cli config set skills.default_mode symlink
  aic-cli config set log.level debug
  aic-cli config set log.max_size 20`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runSet,
}

var sensitiveKeys = map[string]bool{
	"aic.token": true,
}

func runSet(cmd *cobra.Command, args []string) {
	key := args[0]

	_, ok, _ := getConfigValue(key)
	if !ok {
		fmt.Printf("错误: 未知配置项 '%s'\n", key)
		fmt.Println("使用 'aic-cli config list' 查看所有配置项")
		os.Exit(1)
	}

	var value string

	if len(args) >= 2 {
		value = args[1]
	} else if sensitiveKeys[key] {
		fmt.Printf("请输入 %s: ", key)
		raw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("错误: 读取输入失败")
			os.Exit(1)
		}
		value = string(raw)
	} else {
		fmt.Printf("请输入 %s: ", key)
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		value = strings.TrimSpace(line)
	}

	if value == "" {
		fmt.Println("错误: 值不能为空")
		os.Exit(1)
	}

	appConfig := common.AppConfigModel

	if err := setConfigValue(appConfig, key, value); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	utils.ConfigInstance.SaveConfig(appConfig)
	common.AppConfigModel = appConfig

	fmt.Printf("%s 已更新\n", key)
}

func setConfigValue(c *models.AppConfig, key, value string) error {
	switch key {
	case "aic.base_url":
		c.AIC.BaseURL = value
	case "aic.token":
		c.AIC.Token = value
	case "skills.global_dir":
		c.Skills.GlobalDir = value
	case "skills.default_mode":
		if value != "copy" && value != "symlink" {
			return fmt.Errorf("无效的安装模式: %s (copy/symlink)", value)
		}
		c.Skills.DefaultMode = value
	case "log.level":
		c.Log.Level = value
	case "log.max_size":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("无效的整数值: %s", value)
		}
		c.Log.MaxSize = v
	case "log.max_backups":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("无效的整数值: %s", value)
		}
		c.Log.MaxBackups = v
	case "log.max_age":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("无效的整数值: %s", value)
		}
		c.Log.MaxAge = v
	case "log.compress":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("无效的布尔值: %s (true/false)", value)
		}
		c.Log.Compress = v
	}
	return nil
}
