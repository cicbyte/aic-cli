package skill

import (
	"fmt"
	"os"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat <skill-id> <path>",
	Short: "查看远程技能文件内容",
	Long: `读取远程技能的单个文件内容并输出到 stdout。

示例:
  aic-cli skill cat 42 SKILL.md
  aic-cli skill cat 42 prompts/review.md -o local.md`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}
		path := args[1]

		outputFile, _ := cmd.Flags().GetString("output")

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)
		resp, err := client.GetFile(skillID, path)
		if err != nil {
			return err
		}

		if outputFile != "" {
			if err := os.WriteFile(outputFile, []byte(resp.Data.Content), 0644); err != nil {
				return fmt.Errorf("写入文件失败: %w", err)
			}
			fmt.Printf("已保存到 %s (sha256: %s)\n", outputFile, resp.Data.Sha256)
		} else {
			fmt.Print(resp.Data.Content)
		}

		return nil
	},
}

func init() {
	catCmd.Flags().StringP("output", "o", "", "保存到本地文件")
}
