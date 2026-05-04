package skill

import (
	"fmt"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

func GetLoginCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "login [token]",
		Short: "登录 AIC 服务器",
		Long: `使用 Token 登录 AIC 服务器并保存到配置文件。

Token 可在 AIC 服务器的设置页面获取。
登录成功后，Token 会保存到配置文件，后续操作自动使用。

示例:
  aic-cli login my-secret-token
  aic-cli login`,
		Args: cobra.MaximumNArgs(1),
		Run:  runLogin,
	}
}

func runLogin(cmd *cobra.Command, args []string) {
	token := ""
	if len(args) > 0 {
		token = args[0]
	}

	if token == "" {
		fmt.Print("请输入 Token: ")
		fmt.Scanln(&token)
	}

	if token == "" {
		fmt.Println("Token 不能为空")
		return
	}

	client := api.NewClient(common.AppConfigModel.AIC.BaseURL, "")
	resp, err := client.Login(token)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}

	if resp.Code != 0 {
		fmt.Printf("登录失败: %s\n", resp.Message)
		return
	}

	common.AppConfigModel.AIC.Token = token
	utils.ConfigInstance.SaveConfig(common.AppConfigModel)

	fmt.Printf("登录成功! Token: %s\n", resp.Data.Token)
}
