package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var getShowFlag bool

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "查看配置项的值",
	Long: `查看指定配置项的当前值。

敏感字段（如 token）默认显示为脱敏值，使用 --show 查看明文。

示例:
  aic-cli config get aic.base_url
  aic-cli config get skills.default_mode
  aic-cli config get aic.token
  aic-cli config get aic.token --show`,
	Args: cobra.ExactArgs(1),
	Run:  runGet,
}

func init() {
	getCmd.Flags().BoolVar(&getShowFlag, "show", false, "显示敏感字段的明文值")
}

func runGet(cmd *cobra.Command, args []string) {
	key := args[0]

	value, ok, sensitive := getConfigValue(key)
	if !ok {
		fmt.Printf("错误: 未知配置项 '%s'\n", key)
		fmt.Println("使用 'aic-cli config list' 查看所有配置项")
		os.Exit(1)
	}

	if value == "" {
		fmt.Printf("%s: (未设置)\n", key)
		return
	}

	if sensitive && !getShowFlag {
		fmt.Printf("%s: %s\n", key, maskValue(value))
		fmt.Println("使用 --show 查看明文")
		return
	}

	fmt.Printf("%s: %s\n", key, value)
}

func getConfigValue(key string) (string, bool, bool) {
	appConfig := common.AppConfigModel

	switch key {
	case "aic.base_url":
		return appConfig.AIC.BaseURL, true, false
	case "aic.token":
		return appConfig.AIC.Token, true, true
	case "skills.global_dir":
		return appConfig.Skills.GlobalDir, true, false
	case "skills.default_mode":
		return appConfig.Skills.DefaultMode, true, false
	case "log.level":
		return appConfig.Log.Level, true, false
	case "log.max_size":
		return strconv.Itoa(appConfig.Log.MaxSize), true, false
	case "log.max_backups":
		return strconv.Itoa(appConfig.Log.MaxBackups), true, false
	case "log.max_age":
		return strconv.Itoa(appConfig.Log.MaxAge), true, false
	case "log.compress":
		return strconv.FormatBool(appConfig.Log.Compress), true, false
	default:
		return "", false, false
	}
}

func maskValue(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= 8 {
		return "******"
	}
	return string(runes[:4]) + "..." + string(runes[len(runes)-4:])
}
