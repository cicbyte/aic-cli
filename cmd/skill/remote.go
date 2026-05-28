package skill

import (
	"github.com/spf13/cobra"
)

func GetRemoteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "远程技能管理（AI Agent 专用）",
		Long: `远程技能文件的细粒度操作，主要供 AI Agent 自动化使用。

文件操作:
  cat        查看远程文件内容
  tree       显示远程技能文件树
  edit       编辑远程技能文件
  patch      增量编辑技能文件

校验与发布:
  validate   校验技能文件
  publish    发布技能
  unpublish  取消发布

元数据:
  update     更新技能元数据`,
	}

	// 文件操作
	cmd.AddCommand(catCmd)
	cmd.AddCommand(treeCmd)
	cmd.AddCommand(editCmd)
	cmd.AddCommand(patchCmd)

	// 校验与发布
	cmd.AddCommand(validateCmd)
	cmd.AddCommand(publishCmd)
	cmd.AddCommand(unpublishCmd)

	// 元数据
	cmd.AddCommand(updateCmd)

	return cmd
}
