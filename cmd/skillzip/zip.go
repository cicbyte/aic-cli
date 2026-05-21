package skillzip

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskillzip "github.com/cicbyte/aic-cli/internal/logic/skillzip"
	"github.com/spf13/cobra"
)

var (
	packageOutput string
	packageFormat string
)

func GetPackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pack [skill-folder]",
		Short: "打包 skill 文件夹",
		Long: `将 skill 文件夹打包为 .zip 或 .skill 文件。

文件夹中必须包含 skill.md 文件。
默认打包为 .zip 格式，可通过 --format skill 打包为 .skill 格式。

示例:
  aic-cli skill pack ./my-skill
  aic-cli skill pack ./my-skill --format skill
  aic-cli skill pack ./my-skill -o ./output.zip`,
		Args: cobra.MaximumNArgs(1),
		Run:  runPack,
	}
	cmd.Flags().StringVarP(&packageOutput, "output", "o", "", "输出文件路径")
	cmd.Flags().StringVar(&packageFormat, "format", "zip", "输出格式: zip 或 skill")
	return cmd
}

func runPack(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("请指定要打包的 skill 文件夹路径")
		cmd.Help()
		return
	}

	if packageFormat != "zip" && packageFormat != "skill" {
		fmt.Printf("不支持的格式: %s（可选 zip 或 skill）\n", packageFormat)
		os.Exit(1)
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	config := &logicskillzip.ZipConfig{
		InputPath:  args[0],
		OutputPath: packageOutput,
		Format:     packageFormat,
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
