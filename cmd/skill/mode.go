package skill

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var modeCmd = &cobra.Command{
	Use:   "mode [copy|symlink]",
	Short: "查看或切换 skill 安装模式",
	Long: `查看当前安装模式，或切换到指定模式。

不带参数时显示当前模式及说明。
带参数时直接切换到指定模式。

安装模式:
  symlink - 下载到全局目录，创建软连接到目标目录（默认）
  copy    - 直接复制文件到目标目录`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMode,
}

func runMode(cmd *cobra.Command, args []string) {
	current := common.AppConfigModel.Skills.DefaultMode
	if current == "" {
		current = "symlink"
	}

	if len(args) == 0 {
		showMode(current)
		return
	}

	target := args[0]
	if target != "copy" && target != "symlink" {
		fmt.Printf("无效的模式: %s（支持: copy, symlink）\n", target)
		os.Exit(1)
	}

	if target == current {
		fmt.Printf("当前已是 %s 模式\n", current)
		return
	}

	common.AppConfigModel.Skills.DefaultMode = target
	utils.ConfigInstance.SaveConfig(common.AppConfigModel)

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	fmt.Printf("%s 已切换到 %s 模式\n", successStyle.Render("✓"), target)
	showMode(target)
}

func showMode(mode string) {
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	fmt.Printf("当前安装模式: %s\n\n", activeStyle.Render(mode))

	modes := []struct {
		name string
		desc string
	}{
		{"symlink", "下载到全局目录（~/.ciclebyte/aic-cli/skills/），在目标位置创建软连接。节省磁盘空间，多项目共享同一份文件。Windows 使用 Junction，无需管理员权限。"},
		{"copy", "直接复制文件到目标目录。各项目独立，互不影响。"},
	}

	for _, m := range modes {
		marker := "  "
		if m.name == mode {
			marker = infoStyle.Render("▸ ")
		}
		fmt.Printf("%s%s\n", marker, infoStyle.Render(m.name))
		fmt.Printf("    %s\n\n", descStyle.Render(m.desc))
	}

	fmt.Printf("切换模式: %s\n", descStyle.Render("aic-cli skill mode <copy|symlink>"))
}
