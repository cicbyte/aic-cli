package server

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看服务器连接状态",
	Long: `检查 AIC 服务器的连接状态和认证信息。

示例:
  aic-cli server status`,
	Run: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	appConfig := common.AppConfigModel
	baseURL := appConfig.AIC.BaseURL
	token := appConfig.AIC.Token

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	fmt.Printf("  %s %s\n", labelStyle.Render("服务器:"), baseURL)

	// 连接检查
	client := api.NewClient(baseURL, "")
	health, err := client.HealthCheck()
	if err != nil {
		fmt.Printf("  %s %s\n", labelStyle.Render("连接:"), errStyle.Render("失败"))
		fmt.Printf("  %s %v\n", labelStyle.Render("错误:"), err)
		return
	}
	fmt.Printf("  %s %s\n", labelStyle.Render("连接:"), okStyle.Render("正常"))
	if health.Status != "" {
		fmt.Printf("  %s %s\n", labelStyle.Render("状态:"), health.Status)
	}

	// 认证检查
	if token == "" {
		fmt.Printf("  %s %s\n", labelStyle.Render("认证:"), errStyle.Render("未登录"))
		return
	}

	clientWithToken := api.NewClient(baseURL, token)
	_, err = clientWithToken.ListSkills(1, 1, 0, "")
	if err != nil {
		fmt.Printf("  %s %s\n", labelStyle.Render("认证:"), errStyle.Render("Token 无效"))
		return
	}
	fmt.Printf("  %s %s\n", labelStyle.Render("认证:"), okStyle.Render("已登录"))
}
