package skill

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/spf13/cobra"
)

var downloadOutputDir string

func GetDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [skill-id|skill-name]",
		Short: "下载 skill zip 包",
		Long:  `从 AIC 服务器下载指定 skill 的 zip 包`,
		Args:  cobra.MaximumNArgs(1),
		Run:   runDownload,
	}
	cmd.Flags().StringVarP(&downloadOutputDir, "outputDir", "o", "", "输出目录 (默认: 当前目录)")
	return cmd
}

func runDownload(cmd *cobra.Command, args []string) {
	var skillID int
	var skillName string

	if len(args) > 0 {
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			skillName = args[0]
		}
	}

	if skillID == 0 && skillName == "" {
		fmt.Println("请指定要下载的 skill ID 或名称")
		os.Exit(1)
	}

	config := &logicskill.DownloadConfig{
		SkillID:   skillID,
		SkillName: skillName,
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
	fmt.Printf("%s: %s\n", infoStyle.Render("保存到"), result.FilePath)
	fmt.Printf("%s: %s\n", successStyle.Render("下载完成"), result.FilePath)
}
