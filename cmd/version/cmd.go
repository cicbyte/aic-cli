package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 aic-cli 的版本信息，包括版本号、Git提交哈希和构建时间。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aic-cli version %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build time: %s\n", BuildTime)
	},
}
