package pkg

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ltfred/alacritty/pkg/themes"
)

type ThemeChooseModel struct {
	left  list.Model
	right OverlayModel

	themes []themes.Theme
}

func (m ThemeChooseModel) Init() tea.Cmd {
	return nil
}

func NewThemeChooseModel() ThemeChooseModel {
	items, width, themes := themeList()
	l := list.New(items, itemDelegate{}, width, 0)
	l.Title = "Select a theme"
	l.SetShowPagination(true)
	l.SetShowHelp(false)

	right := themeModel{theme: themes[0]}
	over := OverlayModel{theme: themes[0], background: right}

	return ThemeChooseModel{left: l, right: over, themes: themes}
}

func themeList() ([]list.Item, int, []themes.Theme) {
	themeData := themes.GetThemes()
	l := make([]list.Item, 0, len(themeData))
	maxWidth := 0
	for _, theme := range themeData {
		l = append(l, item(theme.Name))
		maxWidth = max(maxWidth, len(theme.Name))
	}
	return l, maxWidth + 5, themeData
}

func (m ThemeChooseModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.left.SetHeight(msg.Height)
		m.right.background.maxWidth = msg.Width - m.left.Width()
		m.right.background.maxHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.right.isConfirming {
				return m, nil
			}
			m.left.CursorDown()
		case "k", "up":
			if m.right.isConfirming {
				return m, nil
			}
			m.left.CursorUp()
		case "enter":
			if !m.right.isConfirming {
				m.right.foreground = newConfirmModel(m.themes[m.left.Cursor()].Name)
				m.right.isConfirming = true
				return m, nil
			}
			m.right.isConfirming = false
			_ = m.right.foreground.yesSelected
			m.right.foreground = confirmModel{}
			return m, nil

		case "left", "right":
			if m.right.isConfirming {
				m.right.foreground.yesSelected = !m.right.foreground.yesSelected
			}
		}
	}
	s := m.themes[m.left.Cursor()]
	m.right.background.theme = s
	m.right.leftWidth = m.left.Width()
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
