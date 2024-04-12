package prompt

import (
	"fmt"

	"github.com/gookit/color"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type Input struct {
	TextInput textinput.Model
	Label     string
	err       error
}

func (m Input) Init() tea.Cmd {
	return textinput.Blink
}

func (m Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			Quit <- true
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Input) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n",
		m.Label, m.TextInput.View(),
	) + "\n"
}

func NewInput(label, placeholder string, validateFunc ...textinput.ValidateFunc) Input {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit, ti.Width = 156, 30
	ti.Placeholder = placeholder
	if len(validateFunc) != 0 {
		ti.Validate = validateFunc[0]
	}

	return Input{TextInput: ti, Label: label}
}

func (m Input) Run() (string, error) {
	program := tea.NewProgram(m)
	model, err := program.Run()
	if err != nil {
		color.Error.Prompt(err.Error())
		return "", err
	}
	if m, ok := model.(Input); ok {
		return m.TextInput.Value(), err
	}

	return "", nil
}
