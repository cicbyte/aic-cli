package pkg

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var listSearch string

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出技能包",
		Long: `列出所有技能包，支持关键词搜索。

示例:
  aic-cli package list
  aic-cli package list -s "前端"`,
		Run: runList,
	}
	cmd.Flags().StringVarP(&listSearch, "search", "s", "", "搜索关键词")
	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

	resp, err := client.ListPackages(1, 100, listSearch)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		os.Exit(1)
	}
	if resp.Code != 0 {
		fmt.Printf("查询失败: %s\n", resp.Message)
		os.Exit(1)
	}

	if len(resp.Data.List) == 0 {
		fmt.Println("没有找到技能包")
		return
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	fmt.Printf("  %s  %-20s  %-6s  %s\n",
		headerStyle.Render("ID"),
		headerStyle.Render("名称"),
		headerStyle.Render("技能数"),
		headerStyle.Render("描述"),
	)

	for _, p := range resp.Data.List {
		desc := p.Description
		if len([]rune(desc)) > 40 {
			desc = string([]rune(desc)[:40]) + "..."
		}
		fmt.Printf("  %s  %-20s  %-6d  %s\n",
			idStyle.Render(fmt.Sprintf("%d", p.ID)),
			p.Name,
			p.SkillCount,
			desc,
		)
	}
}
