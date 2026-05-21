package server

import (
	"fmt"
	"os"

	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出 AIC 服务器",
	Long: `清除保存的 Token，退出当前登录状态。

示例:
  aic-cli server logout`,
	Run: runLogout,
}

func runLogout(cmd *cobra.Command, args []string) {
	if common.AppConfigModel.AIC.Token == "" {
		fmt.Println("当前未登录")
		os.Exit(0)
	}

	common.AppConfigModel.AIC.Token = ""
	utils.ConfigInstance.SaveConfig(common.AppConfigModel)

	fmt.Println("已退出登录")
}
