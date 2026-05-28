package skill

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <skill-id> [path]",
	Short: "编辑远程技能文件",
	Long: `交互式编辑远程技能文件。下载到临时目录，编辑后上传变更。

示例:
  aic-cli skill edit 42                        # 交互式：列出文件供选择
  aic-cli skill edit 42 SKILL.md              # 直接编辑指定文件
  aic-cli skill edit 42 SKILL.md --dry-run    # 只显示差异，不提交`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		// 获取技能详情
		detail, err := client.GetSkillDetail(skillID)
		if err != nil {
			return fmt.Errorf("获取技能详情失败: %w", err)
		}

		// 确定要编辑的文件
		var filePath string
		if len(args) > 1 {
			filePath = args[1]
		} else {
			// 交互式选择文件
			filePath, err = selectFile(client, skillID)
			if err != nil {
				return err
			}
		}

		// 获取文件内容
		fileResp, err := client.GetFile(skillID, filePath)
		if err != nil {
			return fmt.Errorf("获取文件内容失败: %w", err)
		}

		originalContent := fileResp.Data.Content
		originalSha256 := fileResp.Data.Sha256

		// 创建临时目录
		tmpDir, err := os.MkdirTemp("", "aic-edit-*")
		if err != nil {
			return fmt.Errorf("创建临时目录失败: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		// 写入临时文件
		tmpFile := filepath.Join(tmpDir, filepath.Base(filePath))
		if err := os.WriteFile(tmpFile, []byte(originalContent), 0644); err != nil {
			return fmt.Errorf("写入临时文件失败: %w", err)
		}

		// 打开编辑器
		editor := getEditor()
		fmt.Printf("正在打开编辑器: %s\n", editor)
		fmt.Printf("文件: %s\n", tmpFile)

		editCmd := exec.Command(editor, tmpFile)
		editCmd.Stdin = os.Stdin
		editCmd.Stdout = os.Stdout
		editCmd.Stderr = os.Stderr

		if err := editCmd.Run(); err != nil {
			return fmt.Errorf("编辑器退出失败: %w", err)
		}

		// 读取编辑后的内容
		newContent, err := os.ReadFile(tmpFile)
		if err != nil {
			return fmt.Errorf("读取编辑后的文件失败: %w", err)
		}

		newContentStr := string(newContent)

		// 检查是否有变更
		if newContentStr == originalContent {
			fmt.Println("没有变更")
			return nil
		}

		// 计算差异
		originalLines := strings.Split(originalContent, "\n")
		newLines := strings.Split(newContentStr, "\n")
		added := len(newLines) - len(originalLines)

		fmt.Printf("\n变更检测:\n")
		if added > 0 {
			fmt.Printf("  + 新增 %d 行\n", added)
		} else if added < 0 {
			fmt.Printf("  - 删除 %d 行\n", -added)
		} else {
			fmt.Printf("  ~ 修改内容\n")
		}

		if dryRun {
			fmt.Println("\n[DRY RUN] 不提交变更")
			return nil
		}

		// 确认提交
		var confirm bool
		confirmPrompt := huh.NewConfirm().
			Title("确认提交变更?").
			Value(&confirm)
		if err := confirmPrompt.Run(); err != nil {
			return fmt.Errorf("确认已取消")
		}

		if !confirm {
			fmt.Println("已取消")
			return nil
		}

		// 上传变更
		putReq := &api.PutFileRequest{
			Path:           filePath,
			Content:        newContentStr,
			ExpectedSha256: originalSha256,
		}

		putResp, err := client.PutFile(skillID, putReq)
		if err != nil {
			return fmt.Errorf("上传变更失败: %w", err)
		}

		fmt.Printf("✓ 已更新 %s/%s (sha256: %s)\n", detail.Data.Name, filePath, putResp.Data.Sha256)
		return nil
	},
}

func selectFile(client *api.Client, skillID int) (string, error) {
	resp, err := client.GetSkillFiles(skillID)
	if err != nil {
		return "", fmt.Errorf("获取文件列表失败: %w", err)
	}

	// 收集所有文件
	var files []api.FileNode
	var collectFiles func(nodes []api.FileNode)
	collectFiles = func(nodes []api.FileNode) {
		for _, node := range nodes {
			if node.Type == "file" {
				files = append(files, node)
			} else if node.Type == "folder" {
				collectFiles(node.Children)
			}
		}
	}
	collectFiles(resp.Data.Files)

	if len(files) == 0 {
		return "", fmt.Errorf("技能没有文件")
	}

	// 创建选项
	options := make([]huh.Option[string], len(files))
	for i, f := range files {
		label := f.Path
		if f.Size > 0 {
			label += fmt.Sprintf(" (%s)", formatSize(f.Size))
		}
		options[i] = huh.NewOption(label, f.Path)
	}

	var selected string
	selectPrompt := huh.NewSelect[string]().
		Title("选择要编辑的文件:").
		Options(options...).
		Value(&selected)

	if err := selectPrompt.Run(); err != nil {
		return "", fmt.Errorf("选择已取消")
	}

	return selected, nil
}

func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	if runtime.GOOS == "windows" {
		return "notepad"
	}

	return "vi"
}

func init() {
	editCmd.Flags().Bool("dry-run", false, "只显示差异，不提交")
}
