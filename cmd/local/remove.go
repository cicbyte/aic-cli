package local

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logiclocal "github.com/cicbyte/aic-cli/internal/logic/local"
	"github.com/spf13/cobra"
)

var removeGlobal bool

func GetRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [skill-name]",
		Aliases: []string{"rm", "del"},
		Short:   "删除 skill",
		Long:    `删除本地已安装的 skill。`,
		Args:    cobra.MaximumNArgs(1),
		Run:     runRemove,
	}
	cmd.Flags().BoolVarP(&removeGlobal, "global", "g", false, "同时删除全局源文件")
	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("请指定要删除的 skill 名称")
		cmd.Help()
		return
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	config := &logiclocal.RemoveConfig{
		SkillName: args[0],
		Global:    removeGlobal,
	}

	processor := logiclocal.NewRemoveProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("删除失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s: %s\n", successStyle.Render("已删除"), result.SkillName)

	if result.RemovedGlobal {
		fmt.Printf("%s: 全局源文件已删除\n", warnStyle.Render("注意"))
	}

	fmt.Printf("%s: 使用 aic-cli add %s 重新添加\n", warnStyle.Render("提示"), result.SkillName)
}
