package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var (
	downloadOutputDir string
)

var downloadCmd = &cobra.Command{
	Use:   "download [skill-id|skill-name]",
	Short: "下载 skill zip 包",
	Long:  `从 AIC 服务器下载指定 skill 的 zip 包`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL)

		if skillID == 0 && skillName != "" {
			resp, err := client.ListSkills(1, 100, 0)
			if err != nil {
				fmt.Printf("搜索 skill 失败: %v\n", err)
				os.Exit(1)
			}

			for _, skill := range resp.Data.List {
				if skill.Name == skillName {
					skillID = skill.ID
					break
				}
			}

			if skillID == 0 {
				fmt.Printf("未找到名为 '%s' 的 skill\n", skillName)
				os.Exit(1)
			}
		}

		detail, err := client.GetSkillDetail(skillID)
		if err != nil {
			fmt.Printf("获取 skill 信息失败: %v\n", err)
			os.Exit(1)
		}

		if detail.Code != 0 {
			fmt.Printf("错误: %s\n", detail.Message)
			os.Exit(1)
		}

		outputDir := downloadOutputDir
		if outputDir == "" {
			outputDir = "."
		}

		outputPath := filepath.Join(outputDir, detail.Data.Name+".zip")

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

		fmt.Printf("%s: %s\n", infoStyle.Render("下载 skill"), detail.Data.Name)
		fmt.Printf("%s: %s\n", infoStyle.Render("保存到"), outputPath)

		savedPath, err := client.DownloadSkill(skillID, outputPath)
		if err != nil {
			fmt.Printf("下载失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s: %s\n", successStyle.Render("下载完成"), savedPath)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&downloadOutputDir, "outputDir", "o", "", "输出目录 (默认: 当前目录)")
}
