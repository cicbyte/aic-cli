package skill

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/spf13/cobra"
)

var downloadOutputDir string

func GetDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [skill-id|skill-name]",
		Short: "下载 skill zip 包",
		Long: `从 AIC 服务器下载指定 skill 的 zip 包。

支持通过 ID 或名称匹配。
名称匹配时如有多个结果，会提示选择。

示例:
  aic-cli skill download 42
  aic-cli skill download obsidian
  aic-cli skill download "前端" -o ./downloads`,
		Args: cobra.ExactArgs(1),
		Run:  runDownload,
	}
	cmd.Flags().StringVarP(&downloadOutputDir, "outputDir", "o", "", "输出目录 (默认: 当前目录)")
	return cmd
}

func runDownload(cmd *cobra.Command, args []string) {
	client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

	skill, err := resolveSkill(client, args[0])
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		os.Exit(1)
	}

	config := &logicskill.DownloadConfig{
		SkillID:   skill.ID,
		OutputDir: downloadOutputDir,
	}

	processor := logicskill.NewDownloadProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("下载失败: %v\n", err)
		os.Exit(1)
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	fmt.Printf("%s: %s\n", infoStyle.Render("下载 skill"), result.SkillName)
	fmt.Printf("%s: %s\n", successStyle.Render("保存到"), result.FilePath)
}
