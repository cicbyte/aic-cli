package skill

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	addOutputDir string
	addMode      string
)

func GetAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [skill-id|skill-name]",
		Short: "添加 skill 到本地目录",
		Long: `从 AIC 服务器下载并安装 skill 到本地目录。

支持通过 ID 或名称匹配。
名称匹配时如有多个结果，会提示选择。
已存在的 skill 会自动跳过。

安装模式 (--mode):
  copy    - 直接复制文件到目标目录（默认）
  symlink - 下载到全局目录，创建软连接到目标目录

如果是 Claude Code 项目（存在 .claude 目录），默认保存到 .claude/skills 目录。
否则需要通过 -o 参数指定输出目录。`,
		Args: cobra.ExactArgs(1),
		Run:  runAdd,
	}
	cmd.Flags().StringVarP(&addOutputDir, "outputDir", "o", "", "输出目录 (默认: .claude/skills)")
	cmd.Flags().StringVarP(&addMode, "mode", "m", "copy", "安装模式: copy(复制文件) 或 symlink(全局存储+软连接)")
	return cmd
}

func resolveSkill(client *api.Client, arg string) (*api.Skill, error) {
	// ID 直接查询
	var id int
	if _, err := fmt.Sscanf(arg, "%d", &id); err == nil {
		resp, err := client.GetSkillDetail(id)
		if err != nil {
			return nil, err
		}
		if resp.Code != 0 {
			return nil, fmt.Errorf("%s", resp.Message)
		}
		return &resp.Data, nil
	}

	// 名称搜索
	resp, err := client.ListSkills(1, 50, 0, arg)
	if err != nil {
		return nil, fmt.Errorf("搜索 skill 失败: %w", err)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	if len(resp.Data.List) == 0 {
		return nil, fmt.Errorf("未找到 skill: %s", arg)
	}

	// 精确匹配直接使用
	for _, s := range resp.Data.List {
		if s.Name == arg {
			return &s, nil
		}
	}

	// 非精确匹配，让用户选择
	selected, err := promptSelectSkill(resp.Data.List)
	if err != nil {
		return nil, err
	}
	return selected, nil
}

func promptSelectSkill(skills []api.Skill) (*api.Skill, error) {
	options := make([]huh.Option[*api.Skill], len(skills))
	for i := range skills {
		s := &skills[i]
		label := s.Name
		if s.Description != "" {
			desc := s.Description
			if len([]rune(desc)) > 30 {
				desc = string([]rune(desc)[:30]) + "..."
			}
			label += " - " + desc
		}
		options[i] = huh.NewOption(label, s)
	}

	var selected *api.Skill
	err := huh.NewSelect[*api.Skill]().
		Title("找到如下匹配的 skill，请选择:").
		Options(options...).
		Value(&selected).
		Run()
	if err != nil {
		return nil, fmt.Errorf("选择已取消")
	}
	return selected, nil
}

func skillAlreadyInstalled(name, outputDir, mode string) bool {
	if mode == "symlink" {
		symlinkPath := filepath.Join(outputDir, name)
		if utils.FileExists(symlinkPath) || utils.DirExists(symlinkPath) {
			return true
		}
		globalDir := utils.GetGlobalSkillsDir()
		return utils.DirExists(filepath.Join(globalDir, name))
	}
	return utils.DirExists(filepath.Join(outputDir, name))
}

func runAdd(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Changed("mode") && common.AppConfigModel.Skills.DefaultMode != "" {
		addMode = common.AppConfigModel.Skills.DefaultMode
	}

	if addMode != "copy" && addMode != "symlink" {
		fmt.Printf("无效的安装模式: %s (支持: copy, symlink)\n", addMode)
		os.Exit(1)
	}

	outputDir, err := utils.GetSkillsOutputDir(addOutputDir)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
	if outputDir == "" {
		fmt.Println("当前目录不是 Claude Code 项目，请使用 -o 参数指定输出目录")
		os.Exit(1)
	}

	client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

	skill, err := resolveSkill(client, args[0])
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		os.Exit(1)
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	skipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	// 检查是否已安装
	if skillAlreadyInstalled(skill.Name, outputDir, addMode) {
		fmt.Printf("  %s %s（已存在，跳过）\n", skipStyle.Render("-"), skill.Name)
		return
	}

	config := &logicskill.AddConfig{
		SkillID:   skill.ID,
		OutputDir: addOutputDir,
		Mode:      addMode,
	}

	processor := logicskill.NewAddProcessor(config, common.AppConfigModel)
	result, err := processor.Execute(cmd.Context())

	if err != nil {
		fmt.Printf("添加失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  %s %s -> %s\n", successStyle.Render("✓"), result.SkillName, result.InstallPath)

	if addMode == "symlink" {
		fmt.Printf("%s -> 全局目录\n", result.InstallPath)
	}

	fmt.Println()
	fmt.Printf("%s: 在 Claude Code 中使用 /add-dir .claude/skills/%s 添加到上下文\n", warnStyle.Render("提示"), result.SkillName)
}
