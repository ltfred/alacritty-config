package prompt

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Select struct {
	cursor  int
	Result  string
	Label   string
	Choices []string
}

func (m Select) Init() tea.Cmd {
	return nil
}

func (m Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			Quit <- true
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.Result = m.Choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.Choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.Choices) - 1
			}
		}
	}

	return m, nil
}

func (m Select) View() string {
	s := strings.Builder{}
	s.WriteString(m.Label)

	for i := 0; i < len(m.Choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.Choices[i])
		s.WriteString("\n")
	}

	return s.String()
}
