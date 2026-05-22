package models

type AppConfig struct {
	Version string `yaml:"version"` // 版本号，用于升级时判断
	AIC struct {
		BaseURL string `yaml:"base_url"` // AIC 服务器地址
		Token   string `yaml:"token"`    // 认证 Token
	} `yaml:"aic"`
	Skills struct {
		GlobalDir    string `yaml:"global_dir"`    // 全局skills存储目录
		DefaultMode  string `yaml:"default_mode"`  // 默认安装模式: copy 或 symlink
		DefaultAgent string `yaml:"default_agent"` // 默认目标 Agent: claude, cursor, windsurf 等
	} `yaml:"skills"`
	Log struct {
		Level      string `yaml:"level"`
		MaxSize    int    `yaml:"maxSize"`
		MaxBackups int    `yaml:"maxBackups"`
		MaxAge     int    `yaml:"maxAge"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"log"`
}
