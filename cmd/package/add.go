package pkg

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
	pkgAddOutputDir string
	pkgAddMode      string
	pkgAddAgentName string
)

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [package-id|package-name]",
		Short: "添加技能包中的所有技能",
		Long: `下载并安装指定技能包中的所有技能。

支持通过 ID 或名称匹配技能包。
名称匹配时如有多个结果，会提示选择。
已存在的技能会自动跳过。

示例:
  aic-cli package add 1
  aic-cli package add "前端工具"
  aic-cli package add 1 -o ./skills
  aic-cli package add 1 --mode symlink`,
		Args: cobra.ExactArgs(1),
		Run:  runAdd,
	}
	cmd.Flags().StringVarP(&pkgAddOutputDir, "outputDir", "o", "", "输出目录 (覆盖 Agent 默认路径)")
	cmd.Flags().StringVarP(&pkgAddMode, "mode", "m", "", "安装模式: symlink 或 copy")
	cmd.Flags().StringVar(&pkgAddAgentName, "agent", "", "目标 Agent (claude, cursor, continue, amazonq, copilot, windsurf, cline)")
	return cmd
}

func resolvePackage(client *api.Client, arg string) (*api.SkillPackage, error) {
	// ID 直接查询
	var id int
	if _, err := fmt.Sscanf(arg, "%d", &id); err == nil {
		resp, err := client.GetPackageDetail(id)
		if err != nil {
			return nil, err
		}
		if resp.Code != 0 {
			return nil, fmt.Errorf("%s", resp.Message)
		}
		return &resp.Data, nil
	}

	// 名称搜索
	resp, err := client.ListPackages(1, 50, arg)
	if err != nil {
		return nil, fmt.Errorf("搜索技能包失败: %w", err)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	if len(resp.Data.List) == 0 {
		return nil, fmt.Errorf("未找到技能包: %s", arg)
	}

	// 精确匹配直接使用
	for _, p := range resp.Data.List {
		if p.Name == arg {
			return getPackageDetail(client, p.ID)
		}
	}

	// 非精确匹配，都需要用户确认
	selected, err := promptSelectPackage(resp.Data.List)
	if err != nil {
		return nil, err
	}
	return getPackageDetail(client, selected.ID)
}

func getPackageDetail(client *api.Client, id int) (*api.SkillPackage, error) {
	resp, err := client.GetPackageDetail(id)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%s", resp.Message)
	}
	return &resp.Data, nil
}

func promptSelectPackage(packages []api.SkillPackage) (*api.SkillPackage, error) {
	options := make([]huh.Option[*api.SkillPackage], len(packages))
	for i := range packages {
		p := &packages[i]
		label := fmt.Sprintf("%s (%d 个技能)", p.Name, p.SkillCount)
		if p.Description != "" {
			label += " - " + p.Description
		}
		options[i] = huh.NewOption(label, p)
	}

	var selected *api.SkillPackage
	err := huh.NewSelect[*api.SkillPackage]().
		Title("找到如下匹配的技能包，请选择:").
		Options(options...).
		Value(&selected).
		Run()
	if err != nil {
		return nil, fmt.Errorf("选择已取消")
	}
	return selected, nil
}

func getOutputDirAndAgent() (string, agent.AgentProfile, error) {
	if pkgAddOutputDir != "" {
		dir, err := utils.GetSkillsOutputDir(pkgAddOutputDir)
		return dir, nil, err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}
	a, err := utils.SelectAgent(cwd, pkgAddAgentName)
	if err != nil {
		return "", nil, err
	}
	return a.SkillsDir(cwd), a, nil
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
		pkgAddMode = common.AppConfigModel.Skills.DefaultMode
	}

	if pkgAddMode != "copy" && pkgAddMode != "symlink" {
		fmt.Printf("无效的安装模式: %s (支持: copy, symlink)\n", pkgAddMode)
		os.Exit(1)
	}

	outputDir, a, err := getOutputDirAndAgent()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

	pkg, err := resolvePackage(client, args[0])
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		os.Exit(1)
	}

	if len(pkg.Skills) == 0 {
		fmt.Printf("技能包 \"%s\" 中没有技能\n", pkg.Name)
		return
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	skipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	// 过滤已安装的技能
	type skillAction struct {
		skill    api.Skill
		installed bool
	}
	var actions []skillAction
	skipped := 0
	for _, skill := range pkg.Skills {
		if exists, _ := skillAlreadyInstalled(skill.Name, a, outputDir); exists {
			actions = append(actions, skillAction{skill, true})
			skipped++
		} else {
			actions = append(actions, skillAction{skill, false})
		}
	}

	fmt.Printf("%s %s（%d 个技能", infoStyle.Render("技能包:"), pkg.Name, len(pkg.Skills))
	if skipped > 0 {
		fmt.Printf("，%d 个已存在", skipped)
	}
	fmt.Println(")")
	fmt.Println()

	// 显示已跳过的
	for _, a := range actions {
		if a.installed {
			fmt.Printf("  %s %s（已存在，跳过）\n", skipStyle.Render("-"), a.skill.Name)
		}
	}

	// 安装未存在的
	succeeded := 0
	failed := 0
	for _, a := range actions {
		if a.installed {
			continue
		}

		config := &logicskill.AddConfig{
			SkillID:   a.skill.ID,
			OutputDir: outputDir,
			Mode:      pkgAddMode,
			Overwrite: false,
		}

		processor := logicskill.NewAddProcessor(config, common.AppConfigModel)
		result, err := processor.Execute(cmd.Context())
		if err != nil {
			fmt.Printf("  %s %s: %v\n", errStyle.Render("✗"), a.skill.Name, err)
			failed++
			continue
		}

		fmt.Printf("  %s %s -> %s\n", successStyle.Render("✓"), result.SkillName, result.InstallPath)
		succeeded++
	}

	fmt.Println()
	fmt.Printf("完成: %s 新增, %s 已存在, %s 失败\n",
		successStyle.Render(fmt.Sprintf("%d", succeeded)),
		skipStyle.Render(fmt.Sprintf("%d", skipped)),
		errStyle.Render(fmt.Sprintf("%d", failed)),
	)
}
