package prompt

import (
	"strings"

	"github.com/gookit/color"

	tea "github.com/charmbracelet/bubbletea"
)

type Select struct {
	cursor  int
	result  string
	label   string
	choices []string
}

func (m Select) Init() tea.Cmd {
	return nil
}

func (m Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			Quit <- true
			return m, tea.Quit
		case tea.KeyEnter:
			// Send the choice on the channel and exit.
			m.result = m.choices[m.cursor]
			return m, tea.Quit
		case tea.KeyDown:
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
		case tea.KeyUp:
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		default:
		}
	}

	return m, nil
}

func (m Select) View() string {
	s := strings.Builder{}
	s.WriteString(m.label)

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}

	return s.String()
}

func NewSelect(choices []string, label string) Select {
	return Select{
		label:   label,
		choices: choices,
	}
}

func (m Select) Run() (string, error) {
	program := tea.NewProgram(m)
	model, err := program.Run()
	if err != nil {
		color.Error.Prompt(err.Error())
		return "", err
	}
	if m, ok := model.(Select); ok {
		return m.result, nil
	}

	return "", nil
}
