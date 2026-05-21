package skill

import (
	"github.com/cicbyte/aic-cli/cmd/local"
	"github.com/cicbyte/aic-cli/cmd/skillzip"
	"github.com/spf13/cobra"
)

func GetSkillCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "管理 Claude Skills",
		Long: `管理 Claude Skills 的搜索、下载、安装、导入、打包等操作。

远程操作:
  search     搜索远程技能
  add        下载并安装技能
  download   下载技能 ZIP 包
  import     导入本地技能到服务器
  categories 列出所有分类

本地管理:
  list       列出已安装的技能
  remove     移除已安装的技能
  clean      清理失效的软链接

打包:
  pack       将技能文件夹打包为 .zip 或 .skill`,
	}

	// 远程操作
	cmd.AddCommand(GetSearchCommand())
	cmd.AddCommand(GetAddCommand())
	cmd.AddCommand(GetDownloadCommand())
	cmd.AddCommand(GetImportCommand())
	cmd.AddCommand(GetCategoriesCommand())

	// 本地管理
	cmd.AddCommand(local.GetListCommand())
	cmd.AddCommand(local.GetRemoveCommand())
	cmd.AddCommand(local.GetCleanCommand())

	// 打包
	cmd.AddCommand(skillzip.GetPackCommand())

	return cmd
}
