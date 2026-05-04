package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var (
	importZipPath   string
	importDesc      string
	importCategory  int
	importOverwrite bool
)

var importCmd = &cobra.Command{
	Use:   "import [zip-file]",
	Short: "导入本地 ZIP 文件到 AIC 服务器",
	Long: `将本地 skill ZIP 文件导入到 AIC 服务器。

ZIP 文件应包含 skill 的完整文件结构，必须包含 skill.md 文件。
技能名称从 skill.md 的 YAML front matter 中提取。

示例:
  aic-cli import ./my-skill.zip
  aic-cli import ./my-skill.zip -d "技能描述"
  aic-cli import ./my-skill.zip -c 1`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			importZipPath = args[0]
		}

		if importZipPath == "" {
			fmt.Println("请指定要导入的 ZIP 文件路径")
			os.Exit(1)
		}

		if _, err := os.Stat(importZipPath); os.IsNotExist(err) {
			fmt.Printf("ZIP 文件不存在: %s\n", importZipPath)
			os.Exit(1)
		}

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL)

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

		fmt.Printf("%s: %s\n", infoStyle.Render("导入 ZIP 文件"), importZipPath)

		resp, err := client.ImportZip(importZipPath, importDesc, importCategory, importOverwrite)
		if err != nil {
			fmt.Printf("%s: %v\n", errorStyle.Render("导入失败"), err)
			os.Exit(1)
		}

		if resp.Code != 0 {
			fmt.Printf("%s: %s\n", errorStyle.Render("错误"), resp.Message)
			if strings.Contains(resp.Message, "已存在") {
				fmt.Println("使用 --overwrite 参数覆盖")
			}
			os.Exit(1)
		}

		fmt.Printf("%s: %s (ID: %d)\n", successStyle.Render("导入成功"), resp.Data.Name, resp.Data.SkillID)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&importDesc, "description", "d", "", "技能描述")
	importCmd.Flags().IntVarP(&importCategory, "category", "c", 0, "分类 ID")
	importCmd.Flags().BoolVarP(&importOverwrite, "overwrite", "o", false, "覆盖已存在的技能")
}
