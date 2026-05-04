package skill

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/spf13/cobra"
)

var (
	addOutputDir string
	addSkillID   int
	addSkillName string
	addMode      string
)

func GetAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [skill-id|skill-name]",
		Short: "添加 skill 到本地目录",
		Long: `从 AIC 服务器下载并解压 skill 到本地目录。

安装模式 (--mode):
  copy    - 直接复制文件到目标目录（默认）
  symlink - 下载到全局目录，创建软连接到目标目录

如果是 Claude Code 项目（存在 .claude 目录），默认保存到 .claude/skills 目录。
否则需要通过 -o 参数指定输出目录。`,
		Args: cobra.MaximumNArgs(1),
		Run:  runAdd,
	}
	cmd.Flags().StringVarP(&addOutputDir, "outputDir", "o", "", "输出目录 (默认: .claude/skills)")
	cmd.Flags().StringVarP(&addMode, "mode", "m", "copy", "安装模式: copy(复制文件) 或 symlink(全局存储+软连接)")
	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		if _, err := fmt.Sscanf(args[0], "%d", &addSkillID); err != nil {
			addSkillName = args[0]
		}
	}

	if addSkillID == 0 && addSkillName == "" {
		fmt.Println("请指定要添加的 skill ID 或名称")
		os.Exit(1)
	}

	if !cmd.Flags().Changed("mode") && common.AppConfigModel.Skills.DefaultMode != "" {
		addMode = common.AppConfigModel.Skills.DefaultMode
	}

	if addMode != "copy" && addMode != "symlink" {
		fmt.Printf("无效的安装模式: %s (支持: copy, symlink)\n", addMode)
		os.Exit(1)
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	config := &logicskill.AddConfig{
		SkillID:   addSkillID,
		SkillName: addSkillName,
		OutputDir: addOutputDir,
		Mode:      addMode,
	}

	processor := logicskill.NewAddProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())

	if err != nil && strings.Contains(err.Error(), "需要确认覆盖") {
		fmt.Printf("%s\n", err.Error())
		fmt.Print("是否覆盖? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			fmt.Println("已取消")
			os.Exit(0)
		}

		config.Overwrite = true
		result, err = processor.Execute(cmd.Context())
	}

	if err != nil {
		fmt.Printf("添加失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s: %s\n", infoStyle.Render("下载 skill"), result.SkillName)
	fmt.Printf("%s: %s\n", successStyle.Render("添加成功"), result.InstallPath)

	if addMode == "symlink" {
		fmt.Printf("%s -> 全局目录\n", result.InstallPath)
	}

	fmt.Println()
	fmt.Printf("%s: 在 Claude Code 中使用 /add-dir .claude/skills/%s 添加到上下文\n", warnStyle.Render("提示"), result.SkillName)
}
