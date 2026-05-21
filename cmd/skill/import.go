package skill

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	importDescription string
	importCategory    int
	importOverwrite   bool
)

func GetImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [path]",
		Short: "导入本地技能到 AIC 服务器",
		Long: `将本地技能导入到 AIC 服务器，支持以下格式：
  - .zip 文件
  - .skill 文件（等同于 .zip）
  - 包含 skill.md 的文件夹（自动压缩后上传）

示例:
  aic-cli import ./my-skill.zip
  aic-cli import ./my-skill.skill
  aic-cli import ./my-skill-folder/
  aic-cli import ./my-skill.zip -d "技能描述"
  aic-cli import ./my-skill.zip -c 1`,
		Args: cobra.MaximumNArgs(1),
		Run:  runImport,
	}
	cmd.Flags().StringVarP(&importDescription, "description", "d", "", "技能描述")
	cmd.Flags().IntVarP(&importCategory, "category", "c", 0, "分类 ID")
	cmd.Flags().BoolVarP(&importOverwrite, "overwrite", "o", false, "覆盖已存在的技能")
	return cmd
}

func runImport(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("请指定要导入的路径（.zip、.skill 文件或文件夹）")
		cmd.Help()
		return
	}

	inputPath := args[0]
	if !utils.FileExists(inputPath) {
		fmt.Printf("路径不存在: %s\n", inputPath)
		os.Exit(1)
	}

	config := &logicskill.ImportConfig{
		InputPath:   inputPath,
		Description: importDescription,
		CategoryID:  importCategory,
		Overwrite:   importOverwrite,
	}

	processor := logicskill.NewImportProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())
	if err != nil {
		fmt.Printf("导入失败: %v\n", err)
		os.Exit(1)
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	fmt.Printf("%s: %s (ID: %d)\n", successStyle.Render("导入成功"), result.SkillName, result.SkillID)
}
