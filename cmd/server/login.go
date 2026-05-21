package server

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login [token]",
	Short: "登录 AIC 服务器",
	Long: `配置服务器地址并使用 Token 登录 AIC 服务器。

首次登录会引导输入服务器地址和 Token，后续登录可直接传入 Token。
Token 可在 AIC 服务器的设置页面获取。

示例:
  aic-cli server login
  aic-cli server login my-secret-token`,
	Args: cobra.MaximumNArgs(1),
	Run:  runLogin,
}

var defaultBaseURL = "http://localhost:8000"

func runLogin(cmd *cobra.Command, args []string) {
	baseURL := common.AppConfigModel.AIC.BaseURL
	token := ""

	if len(args) > 0 {
		token = args[0]
	}

	// 首次或使用默认地址时，引导输入服务器地址
	if baseURL == "" || baseURL == defaultBaseURL {
		baseURL = defaultBaseURL
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("服务器地址").
					Value(&baseURL).
					Placeholder(defaultBaseURL),
			),
		).Run(); err != nil {
			fmt.Println("已取消")
			os.Exit(1)
		}
		if baseURL == "" {
			baseURL = defaultBaseURL
		}
	}

	// 交互式输入 Token
	if token == "" {
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Token").
					Value(&token).
					EchoMode(huh.EchoModePassword),
			),
		).Run(); err != nil {
			fmt.Println("已取消")
			os.Exit(1)
		}
	}

	if token == "" {
		fmt.Println("错误: Token 不能为空")
		os.Exit(1)
	}

	// 保存地址
	common.AppConfigModel.AIC.BaseURL = baseURL
	utils.ConfigInstance.SaveConfig(common.AppConfigModel)

	// 验证登录
	client := api.NewClient(baseURL, "")
	resp, err := client.Login(token)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		os.Exit(1)
	}

	if resp.Code != 0 {
		fmt.Printf("登录失败: %s\n", resp.Message)
		os.Exit(1)
	}

	common.AppConfigModel.AIC.Token = token
	utils.ConfigInstance.SaveConfig(common.AppConfigModel)

	fmt.Printf("登录成功! Token: %s\n", resp.Data.Token)
}
