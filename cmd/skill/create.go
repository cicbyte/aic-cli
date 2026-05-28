package skill

import (
	"fmt"
	"strings"

	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "创建新技能",
	Long: `创建新技能并上传到服务器。

示例:
  # 交互式创建
  aic-cli skill create

  # 指定名称创建
  aic-cli skill create my-skill

  # 指定分类和标签
  aic-cli skill create my-skill --category 1 --tags "go,cli"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)

		// 获取技能名称
		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			// 交互式输入
			input := huh.NewInput().
				Title("技能名称").
				Placeholder("my-skill").
				Value(&name)
			if err := input.Run(); err != nil {
				return fmt.Errorf("输入已取消")
			}
		}

		if name == "" {
			return fmt.Errorf("技能名称不能为空")
		}

		// 获取分类 ID
		categoryID, _ := cmd.Flags().GetInt("category")
		if categoryID == 0 {
			// 交互式选择分类
			categories, err := client.ListCategories()
			if err == nil && len(categories) > 0 {
				options := make([]huh.Option[int], len(categories))
				for i, cat := range categories {
					options[i] = huh.NewOption(fmt.Sprintf("%s (%s)", cat.Name, cat.Description), cat.ID)
				}
				selectCat := huh.NewSelect[int]().
					Title("选择分类").
					Options(options...).
					Value(&categoryID)
				if err := selectCat.Run(); err != nil {
					categoryID = 0
				}
			}
		}

		// 获取标签
		tagsStr, _ := cmd.Flags().GetString("tags")
		var tags []string
		if tagsStr != "" {
			tags = splitTags(tagsStr)
		}
		// 只有在没有指定任何参数时才交互式输入标签
		// 否则使用命令行参数

		// 获取描述
		description, _ := cmd.Flags().GetString("desc")
		if description == "" {
			input := huh.NewInput().
				Title("技能描述").
				Placeholder("一个很有用的技能").
				Value(&description)
			if err := input.Run(); err != nil {
				description = ""
			}
		}

		// 生成模板内容
		skillTitle := titleCase(name)
		skillMd := generateSkillMd(name, skillTitle, description)

		// 创建技能
		fmt.Printf("正在创建技能 '%s'...\n", name)

		createReq := &api.CreateSkillRequest{
			Name:        name,
			Description: description,
			CategoryId:  categoryID,
			Tags:        tags,
			Files: []api.FileNode{
				{
					Path:    "SKILL.md",
					Content: skillMd,
				},
			},
		}

		resp, err := client.CreateSkill(createReq)
		if err != nil {
			return fmt.Errorf("创建技能失败: %w", err)
		}

		fmt.Printf("✓ 技能已创建 (ID: %d)\n", resp.Data.SkillId)
		fmt.Println()
		fmt.Println("下一步:")
		fmt.Printf("  1. 编辑 SKILL.md: aic-cli skill remote cat %d SKILL.md\n", resp.Data.SkillId)
		fmt.Printf("  2. 校验技能: aic-cli skill remote validate %d\n", resp.Data.SkillId)
		fmt.Printf("  3. 发布技能: aic-cli skill remote publish %d\n", resp.Data.SkillId)

		return nil
	},
}

func titleCase(name string) string {
	words := strings.Split(name, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func generateSkillMd(name, title, description string) string {
	desc := description
	if desc == "" {
		desc = "[TODO: 填写技能描述，说明这个技能做什么以及何时使用]"
	}

	return fmt.Sprintf(`---
name: %s
description: %s
---

# %s

## 概述

[TODO: 1-2 句话说明这个技能能做什么]

## 使用方法

[TODO: 添加使用说明和示例]

## 资源

### scripts/

可执行脚本目录。

### references/

参考文档目录。

### assets/

资源文件目录（模板、图片等）。

---

**不需要的目录可以删除。** 并非每个技能都需要所有类型的资源。
`, name, desc, title)
}

func init() {
	createCmd.Flags().IntP("category", "c", 0, "分类 ID")
	createCmd.Flags().String("tags", "", "标签 (逗号分隔)")
	createCmd.Flags().StringP("desc", "d", "", "技能描述")
}
