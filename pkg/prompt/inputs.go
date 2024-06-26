package prompt

import (
	"strings"

	"github.com/gookit/color"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

type Inputs struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	options    []InputsOptions
	prompt     string
}

type InputsOptions struct {
	Label       string
	Placeholder string
	CharLimit   int
	Validate    textinput.ValidateFunc
}

func NewInputs(options []InputsOptions, prompt string) Inputs {
	inputs := make([]textinput.Model, 0, len(options))
	for i, option := range options {
		t := textinput.New()
		t.Placeholder, t.CharLimit, t.Validate = option.Placeholder, option.CharLimit, option.Validate
		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}
		inputs = append(inputs, t)
	}

	return Inputs{inputs: inputs, prompt: prompt, options: options}
}

func (m Inputs) Run() ([]string, error) {
	program := tea.NewProgram(m)
	model, err := program.Run()
	if err != nil {
		color.Error.Prompt(err.Error())
		return nil, err
	}
	if m, ok := model.(Inputs); ok {
		outputs := make([]string, len(m.inputs))
		for i := range m.inputs {
			outputs[i] = m.inputs[i].Value()
		}
		return outputs, nil
	}

	return nil, nil
}

func (m Inputs) Init() tea.Cmd {
	return textinput.Blink
}

func (m Inputs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			Quit <- true
			return m, tea.Quit
		// Set focus to next input
		case tea.KeyTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			if msg.Type == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}
			// Cycle indexes
			if msg.Type == tea.KeyUp {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		default:
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m Inputs) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Inputs) View() string {
	var b strings.Builder
	b.WriteString(m.prompt)
	b.WriteString("\n\n")
	for i := range m.inputs {
		b.WriteString(m.options[i].Label)
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
