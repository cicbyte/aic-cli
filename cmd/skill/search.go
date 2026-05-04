package skill

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/spf13/cobra"
)

var (
	searchCategoryID int
	searchPageNum    int
	searchPageSize   int
)

func GetSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [keyword]",
		Short: "搜索 skills",
		Long:  `从 AIC 服务器搜索 skills，支持按名称或描述关键词搜索`,
		Args:  cobra.MaximumNArgs(1),
		Run:   runSearch,
	}
	cmd.Flags().IntVarP(&searchCategoryID, "category", "c", 0, "分类 ID 筛选")
	cmd.Flags().IntVarP(&searchPageNum, "page", "p", 1, "页码")
	cmd.Flags().IntVarP(&searchPageSize, "size", "n", 20, "每页数量")
	return cmd
}

func runSearch(cmd *cobra.Command, args []string) {
	keyword := ""
	if len(args) > 0 {
		keyword = args[0]
	}

	config := &logicskill.SearchConfig{
		Keyword:    keyword,
		CategoryID: searchCategoryID,
		PageNum:    searchPageNum,
		PageSize:   searchPageSize,
	}

	processor := logicskill.NewSearchProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
		os.Exit(1)
	}

	if len(result.Skills) == 0 {
		if keyword != "" {
			fmt.Printf("未找到匹配 '%s' 的 skills\n", keyword)
		} else {
			fmt.Println("暂无 skills")
		}
		return
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Padding(0, 1)
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	if keyword != "" {
		fmt.Println(titleStyle.Render(fmt.Sprintf("搜索结果 '%s':", keyword)))
	} else {
		fmt.Println(titleStyle.Render("Skills 列表:"))
	}
	fmt.Println()

	for _, s := range result.Skills {
		fmt.Printf("  %s %s\n", idStyle.Render("["+strconv.Itoa(s.ID)+"]"), nameStyle.Render(s.Name))
		if s.Description != "" {
			fmt.Printf("    %s\n", descStyle.Render(s.Description))
		}
		fmt.Printf("    %s\n", descStyle.Render(fmt.Sprintf("版本: %s | 下载: %d | 收藏: %d", s.Version, s.DownloadCount, s.StarCount)))
		fmt.Println()
	}

	fmt.Printf("共 %d 个 skills\n", result.Total)
}
