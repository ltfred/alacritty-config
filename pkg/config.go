package pkg

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ConfigModel struct {
	form *huh.Form
}

func NewConfigModel() ConfigModel {
	return ConfigModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().Title("1. WINDOW"),
				huh.NewSelect[string]().
					Title("1.1 Choose decorations").
					Options(
						huh.NewOption("Full", "Full"),
						huh.NewOption("None", "None"),
						huh.NewOption("Transparent(macOS only)", "Transparent"),
						huh.NewOption("Buttonless(macOS only)", "Buttonless"),
					).Key("decorations"),

				huh.NewSelect[string]().
					Title("1.2 Choose startup mode").
					Options(
						huh.NewOption("Windowed", "Windowed"),
						huh.NewOption("Maximized", "Maximized"),
						huh.NewOption("Fullscreen", "Fullscreen"),
						huh.NewOption("SimpleFullscreen(macOS only)", "SimpleFullscreen"),
					).Key("startupMode"),
				huh.NewInput().Title("1.3 Input columns").Placeholder("180").Key("columns").Validate(func(
					s string) error {
					if s != "" {
						_, err := strconv.Atoi(s)
						return err
					}
					return nil
				}),
				huh.NewInput().Title("1.4 Input lines").Placeholder("50").Key("lines").Validate(func(s string) error {
					if s != "" {
						_, err := strconv.Atoi(s)
						return err
					}
					return nil
				}),
				huh.NewInput().Title("1.5 Input opacity").Placeholder("1.0").Key("opacity").Validate(func(
					s string) error {
					if s != "" {
						_, err := strconv.ParseFloat(s, 64)
						return err
					}
					return nil
				}),
			),

			huh.NewGroup(
				huh.NewNote().Title("2. FONT"),
				huh.NewInput().Title("2.1 Input font").Key("font"),
				huh.NewInput().Title("2.2 Input font size").Key("fontSize"),
			),

			huh.NewGroup(
				huh.NewNote().Title("3. CURSOR"),
				huh.NewSelect[string]().
					Title("3.1 Choose shape").
					Options(
						huh.NewOption("Block", "Block"),
						huh.NewOption("Underline", "Underline"),
						huh.NewOption("Beam", "Beam"),
					).Key("shape"),
				huh.NewSelect[string]().
					Title("3.2 Choose blinking").
					Options(
						huh.NewOption("Never", "Never"),
						huh.NewOption("Off", "Off"),
						huh.NewOption("On", "On"),
						huh.NewOption("Always", "Always"),
					).Key("blinking"),
			),
		),
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m ConfigModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(message)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m ConfigModel) View() string {
	if m.form.State == huh.StateCompleted {
		class := m.form.GetString("class")
		level := m.form.GetInt("level")
		return fmt.Sprintf("You selected: %s, Lvl. %d", class, level)
	}
	return m.form.View()
}
