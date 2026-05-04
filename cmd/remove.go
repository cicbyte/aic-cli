package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	removeGlobal bool
)

var removeCmd = &cobra.Command{
	Use:     "remove [skill-name]",
	Aliases: []string{"rm", "del"},
	Short:   "移除本地 skill",
	Long: `移除本地 .claude/skills 目录下的 skill。

默认只删除 .claude/skills 下的 skill（软连接或目录）。
使用 -g 参数同时删除全局源文件。`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		skillName := args[0]

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

		// 获取当前工作目录
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("获取当前目录失败: %v\n", err)
			os.Exit(1)
		}

		// 检查是否是 Claude Code 项目
		if !utils.IsClaudeCodeProject(cwd) {
			fmt.Println("当前目录不是 Claude Code 项目（不存在 .claude 目录）")
			os.Exit(1)
		}

		skillPath := cwd + "/.claude/skills/" + skillName

		// 检查 skill 是否存在
		if !utils.FileExists(skillPath) && !utils.DirExists(skillPath) {
			fmt.Printf("skill 不存在: %s\n", skillPath)
			os.Exit(1)
		}

		// 检查是否是软连接
		link := utils.GetLink(skillName)

		// 删除本地 skill（软连接或目录）
		fmt.Printf("%s: %s\n", infoStyle.Render("删除"), skillPath)
		if err := os.RemoveAll(skillPath); err != nil {
			fmt.Printf("删除失败: %v\n", err)
			os.Exit(1)
		}

		// 从映射记录中移除
		if link != nil {
			utils.RemoveLink(skillName)
		}

		// 处理全局源文件
		if removeGlobal {
			globalDir := utils.GetGlobalSkillsDir()
			globalPath := globalDir + "/" + skillName

			if utils.DirExists(globalPath) {
				fmt.Printf("%s: %s\n", infoStyle.Render("删除全局源文件"), globalPath)
				if err := os.RemoveAll(globalPath); err != nil {
					fmt.Printf("%s: 删除全局源文件失败: %v\n", warnStyle.Render("警告"), err)
				}
			} else {
				fmt.Printf("%s: 全局源文件不存在 %s\n", warnStyle.Render("跳过"), globalPath)
			}
		}

		fmt.Printf("%s: %s\n", successStyle.Render("移除成功"), skillName)

		// 提示
		fmt.Println()
		fmt.Printf("%s: 使用 aic-cli add %s 重新添加\n", warnStyle.Render("提示"), skillName)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVarP(&removeGlobal, "global", "g", false, "同时删除全局源文件")
}
