package skill

import (
	"fmt"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <skill-id>",
	Short: "更新技能元数据",
	Long: `更新技能的名称、描述、分类、标签等元数据。

示例:
  aic-cli skill update 42 --name "新名称"
  aic-cli skill update 42 --desc "新描述"
  aic-cli skill update 42 --tags "go,cli,tool"
  aic-cli skill update 42 --category 2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("desc")
		tagsStr, _ := cmd.Flags().GetString("tags")
		categoryID, _ := cmd.Flags().GetInt("category")

		// 检查是否有任何更新
		if name == "" && description == "" && tagsStr == "" && categoryID == 0 {
			return fmt.Errorf("请指定要更新的字段 (--name, --desc, --tags, --category)")
		}

		req := &api.UpdateSkillRequest{
			Id: skillID,
		}

		if name != "" {
			req.Name = name
		}
		if description != "" {
			req.Description = description
		}
		if tagsStr != "" {
			req.Tags = splitTags(tagsStr)
		}
		if categoryID > 0 {
			req.CategoryId = categoryID
		}

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		if err := client.UpdateSkill(req); err != nil {
			return fmt.Errorf("更新失败: %w", err)
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		fmt.Printf("%s 技能已更新 (ID: %d)\n", successStyle.Render("✓"), skillID)

		return nil
	},
}

func splitTags(s string) []string {
	var tags []string
	for _, tag := range splitString(s, ",") {
		tag = trimSpace(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	return split(s, sep)
}

func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}

func init() {
	updateCmd.Flags().String("name", "", "技能名称")
	updateCmd.Flags().String("desc", "", "技能描述")
	updateCmd.Flags().String("tags", "", "标签 (逗号分隔)")
	updateCmd.Flags().Int("category", 0, "分类 ID")
}
