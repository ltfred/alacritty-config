package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	title       string
	yesSelected bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func newConfirmModel(title string) confirmModel {
	return confirmModel{
		title:       title,
		yesSelected: true,
	}
}

func (m confirmModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m confirmModel) View() string {
	if m.title == "" {
		return ""
	}
	yesStyle := lipgloss.NewStyle().Padding(0, 1)
	noStyle := lipgloss.NewStyle().Padding(0, 1)
	if m.yesSelected {
		yesStyle = yesStyle.Background(lipgloss.Color("170"))
		noStyle = noStyle.Background(lipgloss.Color("#808080"))
	} else {
		yesStyle = yesStyle.Background(lipgloss.Color("#808080"))
		noStyle = noStyle.Background(lipgloss.Color("170"))
	}

	yesBtn := yesStyle.Render(" Yes! ")
	noBtn := noStyle.Render(" No. ")

	content := fmt.Sprintf("%s\n\n%s %s",
		fmt.Sprintf("switch theme to %s", m.title),
		yesBtn,
		noBtn)
	width, height := lipgloss.Size(content)
	style := lipgloss.NewStyle().
		Height(height+2).
		Width(width+4).
		Align(lipgloss.Center, lipgloss.Center).
		Background(lipgloss.Color("#000000"))

	return style.Render(content)
}
