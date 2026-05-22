package local

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logiclocal "github.com/cicbyte/aic-cli/internal/logic/local"
	"github.com/spf13/cobra"
)

var (
	listGlobal    bool
	listAgentName string
)

func GetListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "列出本地 skills",
		Long: `列出本地已安装的 skills。

默认显示当前项目 Agent 目录下的 skills。
使用 -g 参数显示全局 skills 目录下的 skills。
使用 --agent 指定目标 Agent。`,
		Run: runList,
	}
	cmd.Flags().BoolVarP(&listGlobal, "global", "g", false, "显示全局 skills 目录")
	cmd.Flags().StringVar(&listAgentName, "agent", "", "目标 Agent (claude, cursor, continue, amazonq, copilot, windsurf, cline)")
	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	cwd, _ := os.Getwd()
	config := &logiclocal.ListConfig{Global: listGlobal, WorkingDir: cwd, AgentName: listAgentName}

	processor := logiclocal.NewListProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("读取目录失败: %v\n", err)
		os.Exit(1)
	}

	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	linkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	dirLabel := "Skills 目录"
	if listGlobal {
		dirLabel = "全局 Skills 目录"
	}
	fmt.Printf("%s: %s\n\n", infoStyle.Render(dirLabel), result.Dir)

	if len(result.Skills) == 0 {
		if listGlobal {
			fmt.Println("没有安装任何全局 skill")
		} else {
			fmt.Println("没有安装任何 skill")
			fmt.Println("提示: 使用 -g 参数查看全局 skills")
		}
		return
	}

	fmt.Printf("共 %d 个 skill:\n\n", len(result.Skills))
	for _, s := range result.Skills {
		marker := ""
		if s.IsLink {
			marker = linkStyle.Render(" [symlink]")
		}
		fmt.Printf("  %s%s\n", nameStyle.Render(s.Name), marker)
	}
}
