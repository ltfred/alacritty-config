package ui

import (
	"fmt"
	"runtime"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/ltfred/alacritty/pkg/config"
)

type ConfigModel struct {
	form      *huh.Form
	originCfg config.Config
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
	"Transparent": "Title bar, transparent background and title bar buttons(macOS only)",
	"Buttonless":  "Title bar, transparent background and no title bar buttons(macOS only)",

	"Windowed":         "Regular window",
	"Maximized":        "The window will be maximized on startup",
	"Fullscreen":       "The window will be fullscreened on startup",
	"SimpleFullscreen": "Same as Fullscreen, but you can stack windows on top(macOS only)",

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

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("1. WINDOW"),
			huh.NewSelect[string]().
				Title("1.1 Choose decorations").
				Options(
					huh.NewOption(optionSelected("Full", cfg.Window.Decorations)),
					huh.NewOption(optionSelected("None", cfg.Window.Decorations)),
					huh.NewOption(optionSelected("Transparent", cfg.Window.Decorations)),
					huh.NewOption(optionSelected("Buttonless", cfg.Window.Decorations)),
				).
				Value(&decorations).
				DescriptionFunc(func() string {
					return descMap[decorations]
				}, &decorations),

			huh.NewSelect[string]().
				Title("1.2 Choose startup mode").
				Options(
					huh.NewOption(optionSelected("Windowed", cfg.Window.StartMode)),
					huh.NewOption(optionSelected("Maximized", cfg.Window.StartMode)),
					huh.NewOption(optionSelected("Fullscreen", cfg.Window.StartMode)),
					huh.NewOption(optionSelected("SimpleFullscreen", cfg.Window.StartMode)),
				).
				Value(&startupMode).
				DescriptionFunc(func() string {
					return descMap[startupMode]
				}, &startupMode),

			huh.NewInput().Title("1.3 window title").Placeholder(cfg.Window.Title).Value(&title).DescriptionFunc(func() string {
				if title == "" {
					return "Default: \"Alacritty\""
				}
				return title
			}, &title),

			huh.NewInput().Title("1.4 Input columns").Value(&column).Validate(validateIntF()).
				Placeholder(strconv.Itoa(cfg.Window.Dimensions.Columns)).DescriptionFunc(func() string {
				if column == "" {
					return "Default: 180"
				}
				return column
			}, &column),

			huh.NewInput().Title("1.5 Input lines").Value(&line).Validate(validateIntF()).
				Placeholder(strconv.Itoa(cfg.Window.Dimensions.Lines)).DescriptionFunc(func() string {
				if line == "" {
					return "Default: 40"
				}
				return line
			}, &line),

			huh.NewInput().Title("1.6 Input opacity").Value(&opacity).Validate(validateFloatF()).DescriptionFunc(func() string {
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
			}, &normalFont).Placeholder(cfg.Font.Normal.Family),
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
			huh.NewInput().Title("2.5 Input font size").Value(&fontSize).Validate(validateFloatF()).DescriptionFunc(func() string {
				if fontSize == "" {
					return "Default: 11.25"
				}
				return fontSize
			}, &fontSize).Placeholder(strconv.FormatFloat(cfg.Font.Size, 'f', -1, 64)),
		),

		huh.NewGroup(
			huh.NewNote().Title("3. CURSOR"),
			huh.NewSelect[string]().
				Title("3.1 Choose shape").
				Options(
					huh.NewOption(optionSelected("Block", cfg.Cursor.Style.Shape)),
					huh.NewOption(optionSelected("Underline", cfg.Cursor.Style.Shape)),
					huh.NewOption(optionSelected("Beam", cfg.Cursor.Style.Shape)),
				).Value(&shape).DescriptionFunc(func() string {
				return descMap[shape]
			}, &shape),
			huh.NewSelect[string]().
				Title("3.2 Choose blinking").
				Options(
					huh.NewOption(optionSelected("Never", cfg.Cursor.Style.BlinkIng)),
					huh.NewOption(optionSelected("Off", cfg.Cursor.Style.BlinkIng)),
					huh.NewOption(optionSelected("On", cfg.Cursor.Style.BlinkIng)),
					huh.NewOption(optionSelected("Always", cfg.Cursor.Style.BlinkIng)),
				).Value(&blinking).DescriptionFunc(func() string {
				return descMap[blinking]
			}, &blinking),
		),
	)

	return ConfigModel{
		form:      form,
		originCfg: cfg,
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m ConfigModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "y", "Y":
			if m.form.State == huh.StateCompleted {
				return m, tea.Quit
			}
		case "n", "N":
			if m.form.State == huh.StateCompleted {
				return m, tea.Quit
			}
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
		return fmt.Sprintf(`
1. WINDOW

decorations: %s
startup mode: %s
title: %s
columns: %s
lines: %s
opacity: %s

2. FONT

normal font: %s
bold font: %s
italic font: %s
bold italic font: %s
font size: %s

3. CURSOR

shape: %s
blinking: %s

Confirm[Y/n]
`,
			resultColor(m.originCfg.Window.Decorations, decorations),
			resultColor(m.originCfg.Window.StartMode, startupMode),
			resultColor(m.originCfg.Window.Title, title),
			resultColor(strconv.Itoa(m.originCfg.Window.Dimensions.Columns), column),
			resultColor(strconv.Itoa(m.originCfg.Window.Dimensions.Lines), line),
			resultColor(strconv.FormatFloat(m.originCfg.Window.Opacity, 'f', -1, 64), opacity),

			resultColor(m.originCfg.Font.Normal.Family, normalFont),
			resultColor(m.originCfg.Font.Bold.Family, boldFont),
			resultColor(m.originCfg.Font.Italic.Family, italicFont),
			resultColor(m.originCfg.Font.BoldItalic.Family, boldItalicFont),
			resultColor(strconv.FormatFloat(m.originCfg.Font.Size, 'f', -1, 64), fontSize),

			resultColor(m.originCfg.Cursor.Style.Shape, shape),
			resultColor(m.originCfg.Cursor.Style.BlinkIng, blinking),
		)
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

func optionSelected(option, selected string) (string, string) {
	if option == selected {
		return option + " (selected)", option
	}
	return option, option
}

func resultColor(origin, new string) string {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	if origin == new {
		return origin + " -> " + new
	}

	if origin == "" && new != "" {
		return green(origin + " -> " + new + " (added)")
	}
	if origin != "" && new != "" && origin != new {
		return yellow(origin + " -> " + new + " (modified)")
	}
	return red(origin + " -> " + new + " (deleted)")
}
