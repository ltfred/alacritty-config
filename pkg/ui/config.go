package ui

import (
	"runtime"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ltfred/alacritty/pkg/config"
)

type ConfigModel struct {
	form *huh.Form
}

var (
	decorations    string
	startupMode    string
	title          string
	column         string
	line           string
	opacity        string
	normalFont     string
	boldFont       string
	italicFont     string
	boldItalicFont string
	fontSize       string
	shape          string
	blinking       string
)

var descMap = map[string]string{
	"Full":        "Borders and title bar",
	"None":        "Neither borders nor title bar",
	"Transparent": "Title bar, transparent background and title bar buttons",
	"Buttonless":  "Title bar, transparent background and no title bar buttons",

	"Windowed":         "Regular window",
	"Maximized":        "The window will be maximized on startup",
	"Fullscreen":       "The window will be fullscreened on startup",
	"SimpleFullscreen": "Same as Fullscreen, but you can stack windows on top",

	"Never":  "Prevent the cursor from ever blinking",
	"Off":    "Disable blinking by default",
	"On":     "Enable blinking by default",
	"Always": "Force the cursor to always blink",

	"Block":     "▌",
	"Underline": "_",
	"Beam":      "|",
}

func NewConfigModel() ConfigModel {
	cfg := config.GetConfigStruct()

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
					).Value(&decorations).DescriptionFunc(func() string {
					return descMap[decorations]
				}, &decorations),

				huh.NewSelect[string]().
					Title("1.2 Choose startup mode").
					Options(
						huh.NewOption("Windowed", "Windowed"),
						huh.NewOption("Maximized", "Maximized"),
						huh.NewOption("Fullscreen", "Fullscreen"),
						huh.NewOption("SimpleFullscreen(macOS only)", "SimpleFullscreen"),
					).Value(&startupMode).DescriptionFunc(func() string {
					return descMap[startupMode]
				}, &startupMode),

				huh.NewInput().Title("1.3 window title").Value(&title).DescriptionFunc(func() string {
					if title == "" {
						return "Default: \"Alacritty\""
					}
					return title
				}, &title),

				huh.NewInput().Title("1.4 Input columns").Value(&column).Validate(validateIntF()).DescriptionFunc(
					func() string {
						if column == "" {
							return "Default: 180"
						}
						return column
					}, &column).Placeholder(strconv.Itoa(cfg.Window.Dimensions.Columns)),

				huh.NewInput().Title("1.5 Input lines").Value(&line).Validate(validateIntF()).DescriptionFunc(func() string {
					if line == "" {
						return "Default: 40"
					}
					return line
				}, &line).Placeholder(strconv.Itoa(cfg.Window.Dimensions.Lines)),

				huh.NewInput().Title("1.6 Input opacity").Value(&opacity).Validate(validateFloatF()).DescriptionFunc(
					func() string {
						if opacity == "" {
							return "Default: 1.0"
						}
						return opacity
					}, &opacity).Placeholder(strconv.FormatFloat(cfg.Window.Opacity, 'f', -1, 64)),
			),

			huh.NewGroup(
				huh.NewNote().Title("2. FONT"),
				huh.NewInput().Title("2.1 Input normal font").Value(&normalFont).DescriptionFunc(func() string {
					if normalFont == "" {
						switch runtime.GOOS {
						case "darwin":
							return "Default: Menlo"
						case "linux":
							return "Default: monospace"
						case "windows":
							return "Default: Consolas"
						default:
							return ""
						}
					}
					return normalFont
				}, &normalFont),
				huh.NewInput().Title("2.2 Input bold font").Value(&boldFont).DescriptionFunc(func() string {
					if boldFont == "" {
						return "If the family is not specified, it will fall back to the value specified for the normal font"
					}
					return boldFont
				}, &boldFont),
				huh.NewInput().Title("2.3 Input italic font").Value(&italicFont).DescriptionFunc(func() string {
					if italicFont == "" {
						return "If the family is not specified, it will fall back to the value specified for the normal font"
					}
					return italicFont
				}, &italicFont),
				huh.NewInput().Title("2.4 Input bold italic font").Value(&boldItalicFont).DescriptionFunc(func() string {
					if boldItalicFont == "" {
						return "If the family is not specified, it will fall back to the value specified for the normal font"
					}
					return boldItalicFont
				}, &boldItalicFont),
				huh.NewInput().Title("2.5 Input font size").Value(&fontSize).Validate(validateFloatF()).
					DescriptionFunc(func() string {
						if fontSize == "" {
							return "Default: 11.25"
						}
						return fontSize
					}, &fontSize),
			),

			huh.NewGroup(
				huh.NewNote().Title("3. CURSOR"),
				huh.NewSelect[string]().
					Title("3.1 Choose shape").
					Options(
						huh.NewOption("Block", "Block"),
						huh.NewOption("Underline", "Underline"),
						huh.NewOption("Beam", "Beam"),
					).Value(&shape).DescriptionFunc(func() string {
					return descMap[shape]
				}, &shape),
				huh.NewSelect[string]().
					Title("3.2 Choose blinking").
					Options(
						huh.NewOption("Never", "Never"),
						huh.NewOption("Off", "Off"),
						huh.NewOption("On", "On"),
						huh.NewOption("Always", "Always"),
					).Value(&blinking).DescriptionFunc(func() string {
					return descMap[blinking]
				}, &blinking),
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
	}
	return m.form.View()
}

func validateIntF() func(s string) error {
	return func(s string) error {
		if s != "" {
			_, err := strconv.Atoi(s)
			return err
		}
		return nil
	}
}

func validateFloatF() func(s string) error {
	return func(s string) error {
		if s != "" {
			_, err := strconv.ParseFloat(s, 32)
			return err
		}
		return nil
	}
}
