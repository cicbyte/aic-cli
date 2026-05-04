package skill

import (
	"fmt"

	"github.com/cicbyte/aic-cli/internal/api"
)

func ResolveSkillID(client *api.Client, skillID int, skillName string) (int, error) {
	if skillID != 0 {
		return skillID, nil
	}
	if skillName == "" {
		return 0, fmt.Errorf("必须指定 skill ID 或名称")
	}
	resp, err := client.ListSkills(1, 100, 0)
	if err != nil {
		return 0, fmt.Errorf("搜索 skill 失败: %w", err)
	}
	for _, skill := range resp.Data.List {
		if skill.Name == skillName {
			return skill.ID, nil
		}
	}
	return 0, fmt.Errorf("未找到名为 '%s' 的 skill", skillName)
}
