package config

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置项及当前值",
	Long: `列出所有配置项及当前值。敏感字段（如 token）会显示脱敏后的值。

示例:
  aic-cli config list`,
	Run: runList,
}

func runList(cmd *cobra.Command, args []string) {
	appConfig := common.AppConfigModel

	type configEntry struct {
		key       string
		section   string
		value     string
		sensitive bool
	}

	entries := []configEntry{
		{key: "aic.base_url", section: "AIC", value: appConfig.AIC.BaseURL},
		{key: "aic.token", section: "AIC", value: appConfig.AIC.Token, sensitive: true},
		{key: "skills.global_dir", section: "Skills", value: appConfig.Skills.GlobalDir},
		{key: "skills.default_mode", section: "Skills", value: appConfig.Skills.DefaultMode},
		{key: "log.level", section: "Log", value: appConfig.Log.Level},
		{key: "log.max_size", section: "Log", value: strconv.Itoa(appConfig.Log.MaxSize)},
		{key: "log.max_backups", section: "Log", value: strconv.Itoa(appConfig.Log.MaxBackups)},
		{key: "log.max_age", section: "Log", value: strconv.Itoa(appConfig.Log.MaxAge)},
		{key: "log.compress", section: "Log", value: strconv.FormatBool(appConfig.Log.Compress)},
	}

	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	currentSection := ""
	for _, e := range entries {
		if e.section != currentSection {
			currentSection = e.section
			fmt.Println(sectionStyle.Render(fmt.Sprintf("[%s]", currentSection)))
		}

		displayVal := e.value
		if e.sensitive {
			displayVal = maskValue(e.value)
		}
		if displayVal == "" {
			displayVal = "(未设置)"
		}

		fmt.Printf("  %s = %s\n", keyStyle.Render(e.key), displayVal)
	}
}
