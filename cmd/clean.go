package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清除失效的软连接",
	Long: `扫描所有通过 symlink 模式创建的软连接，
自动清除全局源文件已删除或软连接失效的记录。`,
	Run: func(cmd *cobra.Command, args []string) {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

		fmt.Printf("%s...\n", infoStyle.Render("扫描失效软连接"))

		cleaned := utils.CleanBrokenLinks()

		if len(cleaned) == 0 {
			fmt.Println("没有发现失效的软连接")
			return
		}

		fmt.Printf("%s %d 个失效软连接: %v\n", successStyle.Render("已清除"), len(cleaned), cleaned)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
