package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	searchCategoryID int
	searchPageNum    int
	searchPageSize   int
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "搜索 skills",
	Long:  `从 AIC 服务器搜索 skills，支持按名称或描述关键词搜索`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := ""
		if len(args) > 0 {
			keyword = args[0]
		}

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL)

		resp, err := client.ListSkills(searchPageNum, searchPageSize, searchCategoryID)
		if err != nil {
			fmt.Printf("搜索失败: %v\n", err)
			os.Exit(1)
		}

		if resp.Code != 0 {
			fmt.Printf("错误: %s\n", resp.Message)
			os.Exit(1)
		}

		var filtered []api.Skill
		if keyword != "" {
			keywordLower := strings.ToLower(keyword)
			for _, skill := range resp.Data.List {
				if strings.Contains(strings.ToLower(skill.Name), keywordLower) ||
					strings.Contains(strings.ToLower(skill.Description), keywordLower) {
					filtered = append(filtered, skill)
				}
			}
		} else {
			filtered = resp.Data.List
		}

		if len(filtered) == 0 {
			if keyword != "" {
				fmt.Printf("未找到匹配 '%s' 的 skills\n", keyword)
			} else {
				fmt.Println("暂无 skills")
			}
			return
		}

		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			Padding(0, 1)

		nameStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

		idStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("6"))

		if keyword != "" {
			fmt.Println(titleStyle.Render(fmt.Sprintf("搜索结果 '%s':", keyword)))
		} else {
			fmt.Println(titleStyle.Render("Skills 列表:"))
		}
		fmt.Println()

		for _, skill := range filtered {
			fmt.Printf("  %s %s\n", idStyle.Render("["+strconv.Itoa(skill.ID)+"]"), nameStyle.Render(skill.Name))
			if skill.Description != "" {
				fmt.Printf("    %s\n", descStyle.Render(skill.Description))
			}
			fmt.Printf("    %s\n", descStyle.Render(fmt.Sprintf("版本: %s | 下载: %d | 收藏: %d", skill.Version, skill.DownloadCount, skill.StarCount)))
			fmt.Println()
		}

		fmt.Printf("共 %d 个 skills\n", len(filtered))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().IntVarP(&searchCategoryID, "category", "c", 0, "分类 ID 筛选")
	searchCmd.Flags().IntVarP(&searchPageNum, "page", "p", 1, "页码")
	searchCmd.Flags().IntVarP(&searchPageSize, "size", "n", 20, "每页数量")
}
