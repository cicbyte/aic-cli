package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/spf13/cobra"
)

var patchCmd = &cobra.Command{
	Use:   "patch <skill-id>",
	Short: "增量编辑技能文件",
	Long: `非交互式增量编辑技能文件，适合 AI Agent 自动化。

模式:
  1. 命令行参数（简单场景）
  2. 管道 JSON（自动检测 stdin）
  3. 批量文件 (--batch)

示例:
  # 命令行参数 - 纯内容匹配
  aic-cli skill patch 42 --path SKILL.md \
    --old "description: 旧描述" \
    --new "description: 新描述"

  # 命令行参数 - 行号 + 内容校验（推荐）
  aic-cli skill patch 42 --path SKILL.md \
    --line 5 \
    --old "description: 旧描述" \
    --new "description: 新描述"

  # 管道 JSON - 单文件编辑
  echo '{"path":"SKILL.md","edits":[{"old_string":"旧","new_string":"新"}]}' | aic-cli skill patch 42

  # 批量文件 - 多文件操作
  aic-cli skill patch 42 --batch edits.json

查看 JSON 格式:
  aic-cli skill patch --schema`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 输出 schema
		showSchema, _ := cmd.Flags().GetBool("schema")
		if showSchema {
			printPatchSchema()
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("请指定 skill ID")
		}

		var skillID int
		if _, err := fmt.Sscanf(args[0], "%d", &skillID); err != nil {
			return fmt.Errorf("无效的 skill ID: %s", args[0])
		}

		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		// 自动检测管道输入
		if hasStdinData() {
			return patchFromStdin(client, skillID)
		}

		batchFile, _ := cmd.Flags().GetString("batch")
		if batchFile != "" {
			return patchFromBatchFile(client, skillID, batchFile)
		}

		return patchFromFlags(client, skillID, cmd)
	},
}

func printPatchSchema() {
	fmt.Println(`# 单文件编辑 JSON 格式:
{
  "path": "SKILL.md",           // 文件路径（必填）
  "expected_sha256": "...",      // 可选，乐观锁
  "edits": [                     // 编辑操作列表（必填）
    {
      "old_string": "旧文本",    // 要替换的文本（与 line_start 至少一个）
      "new_string": "新文本",    // 替换后的文本（必填）
      "replace_all": false,      // 可选，替换所有匹配
      "line_start": 5,           // 可选，起始行号（1-based）
      "line_end": 10             // 可选，结束行号（默认=line_start）
    }
  ]
}

# 三种编辑模式:
# 1. 纯内容匹配: 只需 old_string + new_string
# 2. 纯行号替换: 只需 line_start + new_string
# 3. 行号+内容校验: line_start + old_string + new_string（最安全）

# 批量操作 JSON 格式 (--batch):
{
  "operations": [
    {"op": "patch", "path": "SKILL.md", "edits": [...]},
    {"op": "write", "path": "new.md", "content": "..."},
    {"op": "delete", "path": "old.md"}
  ]
}`)
}

// hasStdinData 检测 stdin 是否有管道数据
func hasStdinData() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	// 检查是否是管道或重定向
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func patchFromStdin(client *api.Client, skillID int) error {
	var req api.PatchFileRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		return fmt.Errorf("解析 stdin JSON 失败: %w", err)
	}

	resp, err := client.PatchFile(skillID, &req)
	if err != nil {
		return err
	}

	fmt.Printf("✓ 已更新 %s (替换 %d 处, sha256: %s)\n",
		resp.Data.Path, resp.Data.TotalReplacements, resp.Data.Sha256)
	return nil
}

func patchFromBatchFile(client *api.Client, skillID int, batchFile string) error {
	data, err := os.ReadFile(batchFile)
	if err != nil {
		return fmt.Errorf("读取批量文件失败: %w", err)
	}

	var req api.BatchRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("解析批量文件失败: %w", err)
	}

	resp, err := client.BatchFiles(skillID, &req)
	if err != nil {
		return err
	}

	for _, result := range resp.Data.Results {
		switch result.Op {
		case "write", "patch":
			fmt.Printf("✓ %s %s (sha256: %s)\n", result.Op, result.Path, result.Sha256)
		case "delete":
			fmt.Printf("✓ 删除 %s\n", result.Path)
		case "rename":
			fmt.Printf("✓ 重命名 %s -> %s\n", result.Path, result.Path)
		}
	}
	return nil
}

func patchFromFlags(client *api.Client, skillID int, cmd *cobra.Command) error {
	path, _ := cmd.Flags().GetString("path")
	if path == "" {
		return fmt.Errorf("必须指定 --path 参数")
	}

	oldString, _ := cmd.Flags().GetString("old")
	newString, _ := cmd.Flags().GetString("new")
	lineStr, _ := cmd.Flags().GetString("line")
	replaceAll, _ := cmd.Flags().GetBool("replace-all")

	if newString == "" {
		return fmt.Errorf("必须指定 --new 参数")
	}

	edit := api.PatchEdit{
		NewString:  newString,
		ReplaceAll: replaceAll,
	}

	// 解析行号
	if lineStr != "" {
		parts := strings.SplitN(lineStr, "-", 2)
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("无效的行号: %s", parts[0])
		}
		edit.LineStart = &start

		if len(parts) > 1 {
			end, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("无效的行号: %s", parts[1])
			}
			edit.LineEnd = &end
		}
	}

	if oldString != "" {
		edit.OldString = oldString
	}

	if edit.OldString == "" && edit.LineStart == nil {
		return fmt.Errorf("必须指定 --old 或 --line 参数")
	}

	req := &api.PatchFileRequest{
		Path:  path,
		Edits: []api.PatchEdit{edit},
	}

	resp, err := client.PatchFile(skillID, req)
	if err != nil {
		return err
	}

	fmt.Printf("✓ 已更新 %s (替换 %d 处, sha256: %s)\n",
		resp.Data.Path, resp.Data.TotalReplacements, resp.Data.Sha256)
	return nil
}

func init() {
	patchCmd.Flags().String("path", "", "文件路径")
	patchCmd.Flags().String("old", "", "要替换的原始文本")
	patchCmd.Flags().String("new", "", "替换后的文本")
	patchCmd.Flags().String("line", "", "行号范围 (如: 5 或 20-25)")
	patchCmd.Flags().Bool("replace-all", false, "替换所有匹配")
	patchCmd.Flags().String("batch", "", "批量操作 JSON 文件")
	patchCmd.Flags().Bool("schema", false, "输出 JSON 格式说明")
}
