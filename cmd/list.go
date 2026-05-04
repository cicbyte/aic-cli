package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var listGlobal bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "列出本地 skills",
	Long: `列出本地已安装的 skills。

默认显示当前项目 .claude/skills 目录下的 skills。
使用 -g 参数显示全局 skills 目录下的 skills。`,
	Run: func(cmd *cobra.Command, args []string) {
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		linkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

		if listGlobal {
			// 显示全局 skills
			globalDir := utils.GetGlobalSkillsDir()
			fmt.Printf("%s: %s\n\n", infoStyle.Render("全局 Skills 目录"), globalDir)

			if !utils.DirExists(globalDir) {
				fmt.Println("全局 skills 目录不存在")
				return
			}

			entries, err := os.ReadDir(globalDir)
			if err != nil {
				fmt.Printf("读取目录失败: %v\n", err)
				return
			}

			if len(entries) == 0 {
				fmt.Println("没有安装任何全局 skill")
				return
			}

			// 收集并排序
			var skills []string
			for _, entry := range entries {
				if entry.IsDir() {
					skills = append(skills, entry.Name())
				}
			}
			sort.Strings(skills)

			fmt.Printf("共 %d 个 skill:\n\n", len(skills))
			for _, name := range skills {
				fmt.Printf("  %s\n", nameStyle.Render(name))
			}
		} else {
			// 显示当前项目 skills
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("获取当前目录失败: %v\n", err)
				os.Exit(1)
			}

			skillsDir := filepath.Join(cwd, ".claude", "skills")
			fmt.Printf("%s: %s\n\n", infoStyle.Render("Skills 目录"), skillsDir)

			if !utils.DirExists(skillsDir) {
				fmt.Println("当前项目没有 .claude/skills 目录")
				fmt.Println("提示: 使用 -g 参数查看全局 skills")
				return
			}

			entries, err := os.ReadDir(skillsDir)
			if err != nil {
				fmt.Printf("读取目录失败: %v\n", err)
				return
			}

			if len(entries) == 0 {
				fmt.Println("没有安装任何 skill")
				return
			}

			// 收集并排序
			type skillInfo struct {
				name   string
				isLink bool
				isDir  bool
			}
			var skills []skillInfo
			for _, entry := range entries {
				info, _ := entry.Info()
				isLink := info.Mode()&os.ModeSymlink != 0
				skills = append(skills, skillInfo{
					name:   entry.Name(),
					isLink: isLink,
					isDir:  entry.IsDir(),
				})
			}
			sort.Slice(skills, func(i, j int) bool {
				return skills[i].name < skills[j].name
			})

			fmt.Printf("共 %d 个 skill:\n\n", len(skills))
			for _, s := range skills {
				marker := ""
				if s.isLink {
					marker = linkStyle.Render(" [symlink]")
				}
				fmt.Printf("  %s%s\n", nameStyle.Render(s.name), marker)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&listGlobal, "global", "g", false, "显示全局 skills 目录")
}
