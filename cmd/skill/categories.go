package skill

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/spf13/cobra"
)

func GetCategoriesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "列出所有分类",
		Long:  `从 AIC 服务器获取所有分类列表`,
		Run:   runCategories,
	}
}

func runCategories(cmd *cobra.Command, args []string) {
	processor := logicskill.NewCategoriesProcessor(&logicskill.CategoriesConfig{}, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("获取分类失败: %v\n", err)
		os.Exit(1)
	}

	if len(result.Categories) == 0 {
		fmt.Println("暂无分类")
		return
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Padding(0, 1)
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	fmt.Println(titleStyle.Render("分类列表:"))
	fmt.Println()

	for _, cat := range result.Categories {
		fmt.Printf("  %s %s\n", idStyle.Render(fmt.Sprintf("[%d]", cat.ID)), nameStyle.Render(cat.Name))
		if cat.Description != "" {
			fmt.Printf("    %s\n", descStyle.Render(cat.Description))
		}
		fmt.Println()
	}

	fmt.Printf("共 %d 个分类\n", len(result.Categories))
}
