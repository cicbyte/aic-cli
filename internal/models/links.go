package models

// SkillLink 表示一个 skill 的软连接映射
type SkillLink struct {
	Name       string `json:"name"`        // skill 名称
	GlobalPath string `json:"global_path"` // 全局目录路径
	LinkPath   string `json:"link_path"`   // 软连接路径
	AddedAt    string `json:"added_at"`    // 添加时间
}

// LinksConfig 存储所有软连接映射
type LinksConfig struct {
	Links []SkillLink `json:"links"`
}
