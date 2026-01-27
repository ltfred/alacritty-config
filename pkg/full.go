package pkg

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ThemeChooseModel struct {
	left  list.Model
	right OverlayModel
}

func (m ThemeChooseModel) Init() tea.Cmd {
	return nil
}

func NewThemeChooseModel(maxWidth int, maxHeight int) ThemeChooseModel {
	items, width := themeList()
	l := list.New(items, itemDelegate{}, width, 0)
	l.Title = "Select a theme"
	l.SetShowPagination(true)
	l.SetShowHelp(false)
	right := liveMode{theme: themes[0]}
	if maxWidth > 0 {
		right.maxWidth = maxWidth - width
	}
	if maxHeight > 0 {
		right.maxHeight = maxHeight
		l.SetHeight(maxHeight)
	}

	return ThemeChooseModel{left: l, right: right}
}

func (m ThemeChooseModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.left.SetHeight(msg.Height)
		m.right.maxWidth = msg.Width - m.left.Width()
		m.right.maxHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.left.CursorDown()
		case "k", "up":
			m.left.CursorUp()
		case "enter":
			confirm := newConfirmModel(themes[m.left.Cursor()].name, m.left.Height(), m.left.Width()+m.right.maxWidth)
			return confirm, tea.EnterAltScreen
		}
	}
	s := themes[m.left.Cursor()]
	m.right.theme = s
	m.right.Update(message)

	return m, nil
}

func (m ThemeChooseModel) View() string {
	var views []string
	views = append(views, m.left.View())

	height := m.left.Height()
	dividerLineStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		Width(1).Height(height).Bold(true)

	divider := dividerLineStyle.Render("")

	views = append(views, divider)
	views = append(views, m.right.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}
