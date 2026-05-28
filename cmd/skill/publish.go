package skill

import (
	"fmt"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish <skill-id>",
	Short: "发布技能",
	Long: `将技能从 draft 变为 public。

示例:
  aic-cli skill publish 42                     # 发布
  aic-cli skill publish 42 --version 1.0.0    # 指定版本
  aic-cli skill publish 42 --changelog "新增功能"  # 带说明`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}

		version, _ := cmd.Flags().GetString("version")
		changelog, _ := cmd.Flags().GetString("changelog")

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		// 先校验
		fmt.Println("正在校验...")
		validateResp, err := client.Validate(skillID, false)
		if err != nil {
			return fmt.Errorf("校验失败: %w", err)
		}

		if !validateResp.Data.Valid {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			fmt.Printf("%s 校验失败，无法发布:\n", errorStyle.Render("✗"))
			for _, e := range validateResp.Data.Errors {
				fmt.Printf("  • %s\n", e)
			}
			return fmt.Errorf("校验失败")
		}

		// 发布
		fmt.Println("正在发布...")
		resp, err := client.Publish(skillID, &api.PublishRequest{
			Version:   version,
			Changelog: changelog,
		})
		if err != nil {
			return fmt.Errorf("发布失败: %w", err)
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		fmt.Printf("%s 技能已发布 (ID: %d, 版本: %s, 状态: %s)\n",
			successStyle.Render("✓"),
			resp.Data.SkillId,
			resp.Data.Version,
			resp.Data.Status)

		return nil
	},
}

var unpublishCmd = &cobra.Command{
	Use:   "unpublish <skill-id>",
	Short: "取消发布",
	Long: `将技能从 public 变为 draft。

示例:
  aic-cli skill unpublish 42`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		if err := client.Unpublish(skillID); err != nil {
			return fmt.Errorf("取消发布失败: %w", err)
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		fmt.Printf("%s 已取消发布 (ID: %d)\n", successStyle.Render("✓"), skillID)

		return nil
	},
}

func init() {
	publishCmd.Flags().String("version", "", "语义化版本号")
	publishCmd.Flags().String("changelog", "", "版本说明")
}
