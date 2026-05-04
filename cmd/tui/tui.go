package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cicbyte/aic-cli/internal/api"
	"github.com/cicbyte/aic-cli/internal/common"
	"github.com/cicbyte/aic-cli/internal/utils"
)

func Run() {
	client := api.NewClient(common.AppConfigModel.AIC.BaseURL, common.AppConfigModel.AIC.Token)
	p := tea.NewProgram(initialModel(client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("启动失败: %v\n", err)
		os.Exit(1)
	}
}

type state int

const (
	stateLoading state = iota
	stateList
	stateDetail
	stateAdding
	stateDone
)

type skillItem struct {
	skill api.Skill
}

func (i skillItem) FilterValue() string { return i.skill.Name }
func (i skillItem) Title() string       { return i.skill.Name }
func (i skillItem) Description() string {
	desc := i.skill.Description
	if len(desc) > 50 {
		desc = desc[:50] + "..."
	}
	return fmt.Sprintf("ID: %d | %s", i.skill.ID, desc)
}

type model struct {
	state    state
	client   *api.Client
	list     list.Model
	spinner  spinner.Model
	skills   []api.Skill
	selected *api.Skill
	detail   *api.Skill
	files    []api.FileNode
	message  string
	err      error
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)

func initialModel(client *api.Client) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(2)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Skills 列表"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return model{
		state:   stateLoading,
		client:  client,
		list:    l,
		spinner: s,
	}
}

type loadSkillsMsg struct {
	skills []api.Skill
	err    error
}

type loadDetailMsg struct {
	detail *api.Skill
	files  []api.FileNode
	err    error
}

type addSkillMsg struct {
	path string
	err  error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			resp, err := m.client.ListSkills(1, 100, 0, "")
			if err != nil {
				return loadSkillsMsg{err: err}
			}
			return loadSkillsMsg{skills: resp.Data.List}
		},
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

		if m.state == stateList {
			switch msg.String() {
			case "enter":
				if i, ok := m.list.SelectedItem().(skillItem); ok {
					m.selected = &i.skill
					m.state = stateLoading
					return m, tea.Batch(
						m.spinner.Tick,
						func() tea.Msg {
							detail, err := m.client.GetSkillDetail(i.skill.ID)
							if err != nil {
								return loadDetailMsg{err: err}
							}
							filesResp, err := m.client.GetSkillFiles(i.skill.ID)
							if err != nil {
								return loadDetailMsg{detail: &detail.Data, err: err}
							}
							return loadDetailMsg{detail: &detail.Data, files: filesResp.Data.Files}
						},
					)
				}
			}
		}

		if m.state == stateDetail {
			switch msg.String() {
			case "a":
				m.state = stateAdding
				return m, tea.Batch(
					m.spinner.Tick,
					func() tea.Msg {
						outputDir, err := utils.GetSkillsOutputDir("")
						if err != nil {
							return addSkillMsg{err: err}
						}
						if outputDir == "" {
							return addSkillMsg{err: fmt.Errorf("不是 Claude Code 项目，请使用 CLI 模式指定输出目录")}
						}

						skillDir := filepath.Join(outputDir, m.selected.Name)
						if utils.DirExists(skillDir) {
							os.RemoveAll(skillDir)
						}

						tmpFile, err := os.CreateTemp("", "skill-*.zip")
						if err != nil {
							return addSkillMsg{err: err}
						}
						tmpPath := tmpFile.Name()
						tmpFile.Close()
						defer os.Remove(tmpPath)

						_, err = m.client.DownloadSkill(m.selected.ID, tmpPath)
						if err != nil {
							return addSkillMsg{err: err}
						}

						if err := utils.Unzip(tmpPath, skillDir); err != nil {
							return addSkillMsg{err: err}
						}

						return addSkillMsg{path: skillDir}
					},
				)
			case "esc", "b":
				m.state = stateList
				m.selected = nil
				m.detail = nil
				m.files = nil
				return m, nil
			}
		}

		if m.state == stateDone {
			if msg.String() == "enter" || msg.String() == "esc" {
				m.state = stateList
				m.selected = nil
				m.detail = nil
				m.files = nil
				m.message = ""
				return m, nil
			}
		}

	case spinner.TickMsg:
		if m.state == stateLoading || m.state == stateAdding {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case loadSkillsMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateDone
			m.message = fmt.Sprintf("加载失败: %v", msg.err)
			return m, nil
		}
		m.skills = msg.skills
		items := make([]list.Item, len(msg.skills))
		for i, s := range msg.skills {
			items[i] = skillItem{skill: s}
		}
		m.list.SetItems(items)
		m.state = stateList
		return m, nil

	case loadDetailMsg:
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("加载详情失败: %v", msg.err)
			m.state = stateList
			return m, nil
		}
		m.detail = msg.detail
		m.files = msg.files
		m.state = stateDetail
		return m, nil

	case addSkillMsg:
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("添加失败: %v", msg.err)
		} else {
			m.message = fmt.Sprintf("添加成功: %s\n\n使用 /add-dir .claude/skills/%s 添加到上下文", msg.path, m.selected.Name)
		}
		m.state = stateDone
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil
	}

	if m.state == stateList {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func flattenFiles(nodes []api.FileNode, b *strings.Builder) {
	for _, f := range nodes {
		b.WriteString(fmt.Sprintf("    %s\n", f.Path))
		if len(f.Children) > 0 {
			flattenFiles(f.Children, b)
		}
	}
}

func (m model) View() string {
	var b strings.Builder

	switch m.state {
	case stateLoading:
		b.WriteString(fmt.Sprintf("\n  %s 加载中...\n", m.spinner.View()))

	case stateList:
		b.WriteString(m.list.View())
		b.WriteString("\n" + helpStyle.Render("  ↑/↓ 导航 | enter 查看 | q 退出"))

	case stateDetail:
		if m.detail != nil {
			b.WriteString(titleStyle.Render("Skill 详情") + "\n\n")
			b.WriteString(fmt.Sprintf("  名称: %s\n", m.detail.Name))
			b.WriteString(fmt.Sprintf("  版本: %s\n", m.detail.Version))
			b.WriteString(fmt.Sprintf("  描述: %s\n", m.detail.Description))
			b.WriteString(fmt.Sprintf("  下载: %d | 收藏: %d\n", m.detail.DownloadCount, m.detail.StarCount))
			if len(m.files) > 0 {
				b.WriteString("\n  文件:\n")
				flattenFiles(m.files, &b)
			}
			b.WriteString("\n" + helpStyle.Render("  a 添加到本地 | esc 返回 | q 退出"))
		}

	case stateAdding:
		b.WriteString(fmt.Sprintf("\n  %s 添加中...\n", m.spinner.View()))

	case stateDone:
		if m.err != nil {
			b.WriteString("\n  " + errorStyle.Render(m.message) + "\n")
		} else {
			b.WriteString("\n  " + selectedStyle.Render(m.message) + "\n")
		}
		b.WriteString("\n" + helpStyle.Render("  enter 返回列表 | q 退出"))
	}

	return b.String()
}
