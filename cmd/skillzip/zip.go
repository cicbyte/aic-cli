package skillzip

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskillzip "github.com/cicbyte/aic-cli/internal/logic/skillzip"
	"github.com/spf13/cobra"
)

var zipOutput string

func GetZipCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zip [skill-folder]",
		Short: "打包 skill 为 ZIP 文件",
		Long: `将 skill 文件夹打包为 ZIP 文件。

文件夹中必须包含 skill.md 文件。

示例:
  aic-cli zip ./my-skill
  aic-cli zip ./my-skill -o ./output.zip`,
		Args: cobra.MaximumNArgs(1),
		Run:  runZip,
	}
	cmd.Flags().StringVarP(&zipOutput, "output", "o", "", "输出 ZIP 文件路径")
	return cmd
}

func runZip(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("请指定要打包的 skill 文件夹路径")
		cmd.Help()
		return
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	config := &logicskillzip.ZipConfig{
		InputPath:  args[0],
		OutputPath: zipOutput,
	}

	processor := logicskillzip.NewZipProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("打包失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s: %s\n", infoStyle.Render("打包完成"), result.OutputPath)
	fmt.Printf("%s: %d 个文件\n", successStyle.Render("文件数量"), result.FileCount)
}
