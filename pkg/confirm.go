package pkg

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	title       string
	yesSelected bool
	focused     bool
	height      int
	width       int
	confirmed   bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func newConfirmModel(title string, height, width int) confirmModel {
	return confirmModel{
		title:       title,
		yesSelected: true,
		focused:     true,
		height:      height,
		width:       width,
		confirmed:   false,
	}
}

func (m confirmModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}
		switch msg.String() {
		case "left", "h", "a":
			m.yesSelected = !m.yesSelected
		case "right", "l", "d":
			m.yesSelected = !m.yesSelected
		case " ", "enter":
			return NewThemeChooseModel(m.width, m.height), nil
		case "y":
			m.yesSelected = true
			m.confirmed = true
			return NewThemeChooseModel(m.width, m.height), nil
		case "n":
			m.yesSelected = false
			m.confirmed = true
			return NewThemeChooseModel(m.width, m.height), nil
		case "q", "esc", "ctrl+c":
			m.yesSelected = false
			return NewThemeChooseModel(m.width, m.height), nil
		}
	}

	return m, nil
}

func (m confirmModel) View() string {
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
		m.title,
		yesBtn,
		noBtn)

	style := lipgloss.NewStyle().
		Height(m.height).
		Width(m.width).
		Align(lipgloss.Center, lipgloss.Center).
		Background(lipgloss.Color("#000000"))

	return style.Render(content)
}

func (m confirmModel) IsConfirmed() bool {
	return m.confirmed && m.yesSelected
}
