package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	addOutputDir string
	addSkillID   int
	addSkillName string
	addMode      string // copy 或 symlink
)

var addCmd = &cobra.Command{
	Use:   "add [skill-id|skill-name]",
	Short: "添加 skill 到本地目录",
	Long: `从 AIC 服务器下载并解压 skill 到本地目录。

安装模式 (--mode):
  copy    - 直接复制文件到目标目录（默认）
  symlink - 下载到全局目录，创建软连接到目标目录

如果是 Claude Code 项目（存在 .claude 目录），默认保存到 .claude/skills 目录。
否则需要通过 -o 参数指定输出目录。`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if _, err := fmt.Sscanf(args[0], "%d", &addSkillID); err != nil {
				addSkillName = args[0]
			}
		}

		if addSkillID == 0 && addSkillName == "" {
			fmt.Println("请指定要添加的 skill ID 或名称")
			os.Exit(1)
		}

		// 如果用户未指定 mode，使用配置文件中的默认值
		if !cmd.Flags().Changed("mode") && common.AppConfigModel.Skills.DefaultMode != "" {
			addMode = common.AppConfigModel.Skills.DefaultMode
		}

		// 验证 mode 参数
		if addMode != "copy" && addMode != "symlink" {
			fmt.Printf("无效的安装模式: %s (支持: copy, symlink)\n", addMode)
			os.Exit(1)
		}

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL)

		if addSkillID == 0 && addSkillName != "" {
			resp, err := client.ListSkills(1, 100, 0)
			if err != nil {
				fmt.Printf("搜索 skill 失败: %v\n", err)
				os.Exit(1)
			}

			for _, skill := range resp.Data.List {
				if skill.Name == addSkillName {
					addSkillID = skill.ID
					break
				}
			}

			if addSkillID == 0 {
				fmt.Printf("未找到名为 '%s' 的 skill\n", addSkillName)
				os.Exit(1)
			}
		}

		detail, err := client.GetSkillDetail(addSkillID)
		if err != nil {
			fmt.Printf("获取 skill 信息失败: %v\n", err)
			os.Exit(1)
		}

		if detail.Code != 0 {
			fmt.Printf("错误: %s\n", detail.Message)
			os.Exit(1)
		}

		outputDir, err := utils.GetSkillsOutputDir(addOutputDir)
		if err != nil {
			fmt.Printf("获取输出目录失败: %v\n", err)
			os.Exit(1)
		}

		if outputDir == "" {
			cwd, _ := os.Getwd()
			fmt.Println("当前目录不是 Claude Code 项目（不存在 .claude 目录）")
			fmt.Printf("当前目录: %s\n", cwd)
			fmt.Println("请使用 -o 参数指定输出目录")
			os.Exit(1)
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

		// 创建临时文件下载
		tmpFile, err := os.CreateTemp("", "skill-*.zip")
		if err != nil {
			fmt.Printf("创建临时文件失败: %v\n", err)
			os.Exit(1)
		}
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		fmt.Printf("%s: %s\n", infoStyle.Render("下载 skill"), detail.Data.Name)
		_, err = client.DownloadSkill(addSkillID, tmpPath)
		if err != nil {
			fmt.Printf("下载失败: %v\n", err)
			os.Exit(1)
		}

		if addMode == "symlink" {
			// symlink 模式：解压到全局目录，创建软连接
			globalSkillsDir := utils.GetGlobalSkillsDir()
			globalSkillDir := filepath.Join(globalSkillsDir, detail.Data.Name)

			fmt.Printf("%s: %s\n", infoStyle.Render("全局目录"), globalSkillDir)

			if utils.DirExists(globalSkillDir) {
				fmt.Printf("%s: 全局目录已存在，跳过下载\n", infoStyle.Render("复用"))
			} else {
				fmt.Printf("%s...\n", infoStyle.Render("解压到全局目录"))
				if err := unzip(tmpPath, globalSkillDir); err != nil {
					fmt.Printf("解压失败: %v\n", err)
					os.Exit(1)
				}
			}

			// 创建软连接
			symlinkPath := filepath.Join(outputDir, detail.Data.Name)
			if utils.FileExists(symlinkPath) || utils.DirExists(symlinkPath) {
				fmt.Printf("目标目录已存在: %s\n", symlinkPath)
				fmt.Print("是否覆盖? (y/N): ")
				var confirm string
				fmt.Scanln(&confirm)
				if strings.ToLower(confirm) != "y" {
					fmt.Println("已取消")
					os.Exit(0)
				}
				os.Remove(symlinkPath)
			}

			// 确保目标父目录存在
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Printf("创建目标目录失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s...\n", infoStyle.Render("创建连接"))
			if err := createLink(globalSkillDir, symlinkPath); err != nil {
				fmt.Printf("创建连接失败: %v\n", err)
				os.Exit(1)
			}

			// 记录映射关系
			if err := utils.AddLink(detail.Data.Name, globalSkillDir, symlinkPath); err != nil {
				fmt.Printf("警告: 记录映射关系失败: %v\n", err)
			}

			fmt.Printf("%s: %s -> %s\n", successStyle.Render("添加成功"), symlinkPath, globalSkillDir)
		} else {
			// copy 模式：直接解压到目标目录
			skillDir := filepath.Join(outputDir, detail.Data.Name)

			fmt.Printf("%s: %s\n", infoStyle.Render("保存到"), skillDir)

			if utils.DirExists(skillDir) {
				fmt.Printf("skill 目录已存在: %s\n", skillDir)
				fmt.Print("是否覆盖? (y/N): ")
				var confirm string
				fmt.Scanln(&confirm)
				if strings.ToLower(confirm) != "y" {
					fmt.Println("已取消")
					os.Exit(0)
				}
				os.RemoveAll(skillDir)
			}

			fmt.Printf("%s...\n", infoStyle.Render("解压中"))
			if err := unzip(tmpPath, skillDir); err != nil {
				fmt.Printf("解压失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s: %s\n", successStyle.Render("添加成功"), skillDir)
		}

		fmt.Println()
		fmt.Printf("%s: 在 Claude Code 中使用 /add-dir .claude/skills/%s 添加到上下文\n", warnStyle.Render("提示"), detail.Data.Name)
	},
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// createLink 创建目录连接（Windows 使用 junction，其他系统使用 symlink）
func createLink(oldname, newname string) error {
	if runtime.GOOS == "windows" {
		// Windows 使用 junction（目录连接点），不需要管理员权限
		cmd := exec.Command("cmd", "/c", "mklink", "/J", newname, oldname)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("创建 junction 失败: %w\n输出: %s", err, string(output))
		}
		return nil
	}
	// Linux/macOS 使用 symlink
	return os.Symlink(oldname, newname)
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addOutputDir, "outputDir", "o", "", "输出目录 (默认: .claude/skills)")
	addCmd.Flags().StringVarP(&addMode, "mode", "m", "copy", "安装模式: copy(复制文件) 或 symlink(全局存储+软连接)")
}
