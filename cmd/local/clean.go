package local

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logiclocal "github.com/cicbyte/aic-cli/internal/logic/local"
	"github.com/spf13/cobra"
)

func GetCleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "清理失效的软链接",
		Long:  `扫描并清除所有失效的软链接。`,
		Run:   runClean,
	}
}

func runClean(cmd *cobra.Command, args []string) {
	processor := logiclocal.NewCleanProcessor(&logiclocal.CleanConfig{}, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("清理失败: %v\n", err)
		return
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	if len(result.Cleaned) == 0 {
		fmt.Println("没有需要清理的软链接")
		return
	}

	for _, name := range result.Cleaned {
		fmt.Printf("  %s %s\n", successStyle.Render("[已清理]"), name)
	}
	fmt.Printf("\n共清理 %d 个失效软链接\n", len(result.Cleaned))
}
