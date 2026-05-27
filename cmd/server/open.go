package server

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "在浏览器中打开 AIC 服务页面",
	Long:  `使用默认浏览器打开 AIC 服务器的 Web 页面。`,
	Run: func(cmd *cobra.Command, args []string) {
		url := common.AppConfigModel.AIC.BaseURL
		if url == "" {
			url = "http://localhost:8000"
		}

		fmt.Printf("正在打开: %s\n", url)

		var err error
		switch runtime.GOOS {
		case "windows":
			err = exec.Command("cmd", "/c", "start", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			err = exec.Command("xdg-open", url).Start()
		}

		if err != nil {
			fmt.Printf("打开浏览器失败: %v\n", err)
			os.Exit(1)
		}
	},
}
