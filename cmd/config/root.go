package config

import (
	"github.com/spf13/cobra"
)

func GetConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "管理应用配置",
		Long: `管理 aic-cli 应用配置（服务器地址、认证 Token、技能安装模式、日志等参数）。

示例:
  aic-cli config list
  aic-cli config get aic.base_url
  aic-cli config set aic.base_url http://192.168.1.100:8000
  aic-cli config set aic.token                    # 交互式输入（不回显）`,
	}
	cmd.AddCommand(listCmd, getCmd, setCmd)
	return cmd
}
