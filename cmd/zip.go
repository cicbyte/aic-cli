package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	zipInputPath  string
	zipOutputPath string
)

var zipCmd = &cobra.Command{
	Use:   "zip [skill-folder]",
	Short: "将 skill 文件夹打包为 ZIP 文件",
	Long: `将 skill 文件夹打包为 ZIP 文件。

要求:
  - 文件夹第一层级必须包含 skill.md 文件
  - 打包后的 ZIP 解压时不会有多余的嵌套层级

示例:
  aic-cli zip ./my-skill
  aic-cli zip ./my-skill -o ./output.zip`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			zipInputPath = args[0]
		}

		if zipInputPath == "" {
			fmt.Println("请指定要打包的 skill 文件夹路径")
			os.Exit(1)
		}

		info, err := os.Stat(zipInputPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("文件夹不存在: %s\n", zipInputPath)
			} else {
				fmt.Printf("读取文件夹失败: %v\n", err)
			}
			os.Exit(1)
		}

		if !info.IsDir() {
			fmt.Printf("指定的路径不是文件夹: %s\n", zipInputPath)
			os.Exit(1)
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

		skillMdPath := filepath.Join(zipInputPath, "skill.md")
		if _, err := os.Stat(skillMdPath); os.IsNotExist(err) {
			fmt.Printf("%s: 文件夹第一层级必须包含 skill.md 文件\n", errorStyle.Render("错误"))
			fmt.Printf("期望路径: %s\n", skillMdPath)
			os.Exit(1)
		}

		if zipOutputPath == "" {
			zipOutputPath = strings.TrimSuffix(filepath.Base(zipInputPath), string(filepath.Separator)) + ".zip"
		}

		fmt.Printf("%s: %s\n", infoStyle.Render("打包文件夹"), zipInputPath)
		fmt.Printf("%s: %s\n", infoStyle.Render("输出文件"), zipOutputPath)

		fileCount, err := createZip(zipInputPath, zipOutputPath)
		if err != nil {
			fmt.Printf("%s: %v\n", errorStyle.Render("打包失败"), err)
			os.Exit(1)
		}

		fmt.Printf("%s: 共 %d 个文件\n", successStyle.Render("打包成功"), fileCount)
	},
}

func createZip(srcDir, outputPath string) (int, error) {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("创建 ZIP 文件失败: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileCount := 0
	srcDir = strings.TrimSuffix(strings.TrimSuffix(srcDir, "/"), "\\")
	baseDir := filepath.Base(srcDir)

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		zipPath := baseDir + "/" + strings.ReplaceAll(relPath, "\\", "/")

		if info.IsDir() {
			_, err = zipWriter.Create(zipPath + "/")
			return err
		}

		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		fileCount++
		return nil
	})

	if err != nil {
		return 0, err
	}

	return fileCount, nil
}

func init() {
	rootCmd.AddCommand(zipCmd)

	zipCmd.Flags().StringVarP(&zipOutputPath, "output", "o", "", "输出 ZIP 文件路径 (默认: 文件夹名.zip)")
}
