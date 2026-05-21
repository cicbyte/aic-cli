package pkg

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	logicskill "github.com/cicbyte/aic-cli/internal/logic/skill"
	"github.com/spf13/cobra"
)

var (
	pkgAddOutputDir string
	pkgAddMode      string
)

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [package-id|package-name]",
		Short: "添加技能包中的所有技能",
		Long: `下载并安装指定技能包中的所有技能。

支持通过 ID 或名称匹配技能包。
名称匹配时会搜索远程技能包，找到第一个匹配项。

示例:
  aic-cli package add 1
  aic-cli package add "前端工具"
  aic-cli package add 1 -o ./skills
  aic-cli package add 1 --mode symlink`,
		Args: cobra.ExactArgs(1),
		Run:  runAdd,
	}
	cmd.Flags().StringVarP(&pkgAddOutputDir, "outputDir", "o", "", "输出目录 (默认: .claude/skills)")
	cmd.Flags().StringVarP(&pkgAddMode, "mode", "m", "copy", "安装模式: copy 或 symlink")
	return cmd
}

func resolvePackage(client *api.Client, arg string) (*api.SkillPackage, error) {
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

	resp, err := client.ListPackages(1, 10, arg)
	if err != nil {
		return nil, fmt.Errorf("搜索技能包失败: %w", err)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	for _, p := range resp.Data.List {
		if p.Name == arg {
			detail, err := client.GetPackageDetail(p.ID)
			if err != nil {
				return nil, err
			}
			if detail.Code != 0 {
				return nil, fmt.Errorf("%s", detail.Message)
			}
			return &detail.Data, nil
		}
	}

	if len(resp.Data.List) > 0 {
		detail, err := client.GetPackageDetail(resp.Data.List[0].ID)
		if err != nil {
			return nil, err
		}
		if detail.Code != 0 {
			return nil, fmt.Errorf("%s", detail.Message)
		}
		return &detail.Data, nil
	}

	return nil, fmt.Errorf("未找到技能包: %s", arg)
}

func runAdd(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Changed("mode") && common.AppConfigModel.Skills.DefaultMode != "" {
		pkgAddMode = common.AppConfigModel.Skills.DefaultMode
	}

	if pkgAddMode != "copy" && pkgAddMode != "symlink" {
		fmt.Printf("无效的安装模式: %s (支持: copy, symlink)\n", pkgAddMode)
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

	fmt.Printf("%s %s（%d 个技能）\n\n", infoStyle.Render("技能包:"), pkg.Name, len(pkg.Skills))

	succeeded := 0
	failed := 0

	for _, skill := range pkg.Skills {
		config := &logicskill.AddConfig{
			SkillID:   skill.ID,
			OutputDir: pkgAddOutputDir,
			Mode:      pkgAddMode,
			Overwrite: true,
		}

		processor := logicskill.NewAddProcessor(config, common.AppConfigModel)
		result, err := processor.Execute(cmd.Context())
		if err != nil {
			fmt.Printf("  %s %s: %v\n", errStyle.Render("✗"), skill.Name, err)
			failed++
			continue
		}

		fmt.Printf("  %s %s -> %s\n", successStyle.Render("✓"), result.SkillName, result.InstallPath)
		succeeded++
	}

	fmt.Println()
	fmt.Printf("完成: %s 成功, %s 失败\n",
		successStyle.Render(fmt.Sprintf("%d", succeeded)),
		errStyle.Render(fmt.Sprintf("%d", failed)),
	)

	if pkgAddMode == "symlink" {
		fmt.Println("提示: 已通过软连接安装到全局目录")
	}
}
