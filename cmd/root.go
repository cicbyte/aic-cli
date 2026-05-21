package cmd

import (
	"fmt"
	"os"

	"github.com/cicbyte/aic-cli/cmd/config"
	"github.com/cicbyte/aic-cli/cmd/local"
	"github.com/cicbyte/aic-cli/cmd/skill"
	"github.com/cicbyte/aic-cli/cmd/skillzip"
	"github.com/cicbyte/aic-cli/cmd/tui"
	"github.com/cicbyte/aic-cli/cmd/version"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/log"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aic-cli",
	Short: "AIC CLI - 管理 Claude Skills 的命令行工具",
	Long: `AIC CLI 是一个用于管理 Claude Skills 的命令行工具。

支持功能:
  - 搜索 skills (search)
  - 下载 skill zip 包
  - 添加 skill 到本地目录
  - 交互式 TUI 界面

如果是 Claude Code 项目（存在 .claude 目录），
skill 会自动保存到 .claude/skills 目录。

不带任何参数运行将进入交互式 TUI 界面。`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.Run()
	},
}

func Execute() {
	if len(os.Args) == 1 {
		initApp()
		tui.Run()
		return
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initApp)

	// skill 模块（远程操作）
	rootCmd.AddCommand(skill.GetLoginCommand())
	rootCmd.AddCommand(skill.GetSearchCommand())
	rootCmd.AddCommand(skill.GetAddCommand())
	rootCmd.AddCommand(skill.GetDownloadCommand())
	rootCmd.AddCommand(skill.GetImportCommand())
	rootCmd.AddCommand(skill.GetCategoriesCommand())

	// local 模块（本地管理）
	rootCmd.AddCommand(local.GetListCommand())
	rootCmd.AddCommand(local.GetRemoveCommand())
	rootCmd.AddCommand(local.GetCleanCommand())

	// zip 模块（打包）
	rootCmd.AddCommand(skillzip.GetZipCommand())

	// config 模块
	rootCmd.AddCommand(config.GetConfigCommand())

	// version 模块
	rootCmd.AddCommand(version.Cmd)
}

func initApp() {
	if err := utils.InitAppDirs(); err != nil {
		fmt.Printf("初始化目录失败: %v\n", err)
		os.Exit(1)
	}
	common.AppConfigModel = utils.ConfigInstance.LoadConfig()
	if err := log.Init(utils.ConfigInstance.GetLogPath()); err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		os.Exit(1)
	}
}
