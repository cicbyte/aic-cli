package pkg

import (
	"github.com/spf13/cobra"
)

func GetPackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg"},
		Short:   "管理技能包",
		Long: `管理技能包：查看、搜索、添加技能包中的所有技能。

示例:
  aic-cli package list
  aic-cli package list -s "前端"
  aic-cli package add 1
  aic-cli package add "前端工具"`,
	}
	cmd.AddCommand(getListCommand(), getAddCommand())
	return cmd
}
