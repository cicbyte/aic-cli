package skill

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/agent"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	addOutputDir  string
	addMode       string
	addAgentName  string
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
  symlink - 下载到全局目录，创建软连接到目标目录（默认）
  copy    - 直接复制文件到目标目录

如果是 Claude Code 项目（存在 .claude 目录），默认保存到 .claude/skills 目录。
否则需要通过 -o 参数指定输出目录。`,
		Args: cobra.ExactArgs(1),
		Run:  runAdd,
	}
	cmd.Flags().StringVarP(&addOutputDir, "outputDir", "o", "", "输出目录 (覆盖 Agent 默认路径)")
	cmd.Flags().StringVarP(&addMode, "mode", "m", "", "安装模式: symlink(全局存储+软连接) 或 copy(复制文件)")
	cmd.Flags().StringVar(&addAgentName, "agent", "", "目标 Agent (claude, cursor, continue, amazonq, copilot, windsurf, cline)")
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

func skillAlreadyInstalled(name string, a agent.AgentProfile, outputDir string) (bool, string) {
	// 检查项目级目录
	localPath := filepath.Join(outputDir, name)
	if utils.FileExists(localPath) || utils.DirExists(localPath) {
		return true, localPath
	}
	// 检查 Agent 全局 skill 目录（如 ~/.claude/skills/）
	if a != nil && a.HasProjectSkills() {
		globalPath := filepath.Join(a.GlobalSkillsDir(), name)
		if utils.FileExists(globalPath) || utils.DirExists(globalPath) {
			return true, globalPath
		}
	}
	return false, ""
}

func runAdd(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Changed("mode") && common.AppConfigModel.Skills.DefaultMode != "" {
		addMode = common.AppConfigModel.Skills.DefaultMode
	}

	if addMode != "copy" && addMode != "symlink" {
		fmt.Printf("无效的安装模式: %s (支持: copy, symlink)\n", addMode)
		os.Exit(1)
	}

	var (
		outputDir string
		a         agent.AgentProfile
		err       error
	)
	if addOutputDir != "" {
		outputDir, err = utils.GetSkillsOutputDir(addOutputDir)
	} else {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			fmt.Printf("错误: %v\n", cwdErr)
			os.Exit(1)
		}
		a, err = utils.SelectAgent(cwd, addAgentName)
		if err == nil {
			outputDir = a.SkillsDir(cwd)
		}
	}
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

	skill, err := resolveSkill(client, args[0])
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		os.Exit(1)
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	skipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	// 检查是否已安装
	if exists, existingPath := skillAlreadyInstalled(skill.Name, a, outputDir); exists {
		fmt.Printf("  %s %s（已存在，跳过）\n", skipStyle.Render("-"), skill.Name)
		fmt.Printf("    路径: %s\n", existingPath)
		return
	}

	config := &logicskill.AddConfig{
		SkillID:   skill.ID,
		OutputDir: outputDir,
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
		fmt.Printf("  （软链接，实际文件存储在全局缓存目录）\n")
	}
}
