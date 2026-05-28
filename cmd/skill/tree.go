package skill

import (
	"fmt"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:   "tree <skill-id>",
	Short: "显示远程技能文件树",
	Long: `显示远程技能的文件结构。

示例:
  aic-cli skill tree 42`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		// 获取技能详情
		detail, err := client.GetSkillDetail(skillID)
		if err != nil {
			return fmt.Errorf("获取技能详情失败: %w", err)
		}

		// 获取文件树
		resp, err := client.GetSkillFiles(skillID)
		if err != nil {
			return fmt.Errorf("获取文件树失败: %w", err)
		}

		// 样式
		nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
		idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

		fmt.Printf("%s %s\n", nameStyle.Render(detail.Data.Name), idStyle.Render(fmt.Sprintf("(ID: %d)", skillID)))

		for i, node := range resp.Data.Files {
			printFileNode(node, "", i == len(resp.Data.Files)-1)
		}

		return nil
	},
}

func printFileNode(node api.FileNode, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	dirStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	sizeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	if node.Type == "folder" {
		fmt.Printf("%s%s%s\n", prefix, dirStyle.Render(node.Name+"/"), sizeStyle.Render(formatSize(node.Size)))
		childPrefix := prefix + "│   "
		if isLast {
			childPrefix = prefix + "    "
		}
		for i, child := range node.Children {
			printFileNode(child, childPrefix, i == len(node.Children)-1)
		}
	} else {
		sizeStr := ""
		if node.Size > 0 {
			sizeStr = " " + sizeStyle.Render(formatSize(node.Size))
		}
		fmt.Printf("%s%s%s%s\n", prefix, connector, nameStyle.Render(node.Name), sizeStr)
	}
}

func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	} else {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	}
}
