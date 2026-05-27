package cmd

import (
	"fmt"
	"os"

	"github.com/cicbyte/aic-cli/cmd/config"
	pkg "github.com/cicbyte/aic-cli/cmd/package"
	"github.com/cicbyte/aic-cli/cmd/server"
	"github.com/cicbyte/aic-cli/cmd/skill"
	"github.com/cicbyte/aic-cli/cmd/version"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/log"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aic-cli",
	Short: "AIC CLI - 管理 Coding Agent Skills 的命令行工具",
	Long: `AIC CLI 是一个用于管理 Coding Agent Skills 的命令行工具。

支持功能:
  - 搜索、下载、安装 skills
  - 技能包管理
  - 多 Agent 支持（Claude Code、Cursor、Windsurf 等）

使用 aic-cli --help 查看所有命令。`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initApp)

	// server 模块（连接管理）
	rootCmd.AddCommand(server.GetServerCommand())

	// skill 模块（技能管理）
	rootCmd.AddCommand(skill.GetSkillCommand())

	// package 模块（技能包管理）
	rootCmd.AddCommand(pkg.GetPackageCommand())

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
