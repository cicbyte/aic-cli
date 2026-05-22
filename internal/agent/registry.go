package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

// RegisteredAgents 是所有已注册的 Agent 列表
var RegisteredAgents = []AgentProfile{
	&ClaudeAgent{},
	&CursorAgent{},
	&ContinueAgent{},
	&AmazonQAgent{},
	&CopilotAgent{},
	&WindsurfAgent{},
	&ClineAgent{},
	&OpenCodeAgent{},
	&CodexAgent{},
	&GeminiAgent{},
	&RooAgent{},
	&AmpAgent{},
	&TraeAgent{},
}

// DetectAgents 检测当前项目中已安装的 Agent
func DetectAgents(projectDir string) []AgentProfile {
	var detected []AgentProfile
	for _, agent := range RegisteredAgents {
		if agent.Detect(projectDir) {
			detected = append(detected, agent)
		}
	}
	return detected
}

// GetAgent 按名称查找 Agent
func GetAgent(name string) (AgentProfile, error) {
	for _, agent := range RegisteredAgents {
		if agent.Name() == name {
			return agent, nil
		}
	}
	return nil, fmt.Errorf("未知的 Agent: %s", name)
}

// detectDir 检测目录是否存在
func detectDir(projectDir, relPath string) bool {
	info, err := os.Stat(filepath.Join(projectDir, relPath))
	return err == nil && info.IsDir()
}

// detectFile 检测文件是否存在
func detectFile(projectDir, relPath string) bool {
	info, err := os.Stat(filepath.Join(projectDir, relPath))
	return err == nil && !info.IsDir()
}

// detectDirOrFile 检测目录或文件是否存在
func detectDirOrFile(projectDir, relPath string) bool {
	_, err := os.Stat(filepath.Join(projectDir, relPath))
	return err == nil
}

