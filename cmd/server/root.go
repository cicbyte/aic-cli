package server

import (
	"github.com/spf13/cobra"
)

func GetServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "管理 AIC 服务器连接",
		Long: `管理与 AIC 服务器的连接（登录、退出等）。

示例:
  aic-cli server login <token>
  aic-cli server login
  aic-cli server logout`,
	}
	cmd.AddCommand(loginCmd, logoutCmd)
	return cmd
}
