package skill

import (
	"fmt"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <skill-id>",
	Short: "校验技能文件",
	Long: `校验远程技能文件的完整性。

示例:
  aic-cli skill validate 42                    # 标准校验
  aic-cli skill validate 42 --strict          # 严格模式`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}

		strict, _ := cmd.Flags().GetBool("strict")
		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		resp, err := client.Validate(skillID, strict)
		if err != nil {
			return err
		}

		// 样式
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

		if resp.Data.Valid {
			fmt.Printf("%s 校验通过\n", successStyle.Render("✓"))
		} else {
			fmt.Printf("%s 校验失败\n", errorStyle.Render("✗"))
		}

		if len(resp.Data.Errors) > 0 {
			fmt.Printf("\n%s\n", errorStyle.Render("错误:"))
			for _, err := range resp.Data.Errors {
				fmt.Printf("  • %s\n", err)
			}
		}

		if len(resp.Data.Warnings) > 0 {
			fmt.Printf("\n%s\n", warningStyle.Render("警告:"))
			for _, warn := range resp.Data.Warnings {
				fmt.Printf("  • %s\n", warn)
			}
		}

		if !resp.Data.Valid {
			return fmt.Errorf("校验失败")
		}

		return nil
	},
}

func init() {
	validateCmd.Flags().Bool("strict", false, "严格模式")
}
