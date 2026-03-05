package ui

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/ltfred/alacritty/pkg/config"
)

type ConfigModel struct {
	form           *huh.Form
	oldCfg         config.Config
	oldCfgMap      map[string]any
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
}

const (
	defaultColumns  = 180
	defaultLines    = 40
	defaultOpacity  = 1.0
	defaultFontSize = 11.25
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
	cfg, err := config.GetConfigStruct()
	if err != nil {
		color.New(color.FgRed).PrintFunc()(err)
		os.Exit(1)
	}
	configMap, err := config.GetConfigMap()
	if err != nil {
		color.New(color.FgRed).PrintFunc()(err)
		os.Exit(1)
	}

	model := ConfigModel{
		oldCfg:    cfg,
		oldCfgMap: configMap,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("1. WINDOW"),
			createDecorationsSelect(cfg.Window.Decorations, &model.decorations),
			createStartupModeSelect(cfg.Window.StartMode, &model.startupMode),
			createTitleInput(cfg.Window.Title, &model.title),
			createColumnsInput(cfg.Window.Dimensions.Columns, &model.column),
			createLinesInput(cfg.Window.Dimensions.Lines, &model.line),
			createOpacityInput(cfg.Window.Opacity, &model.opacity),
		),

		huh.NewGroup(
			huh.NewNote().Title("2. FONT"),
			createNormalFontInput(cfg.Font.Normal.Family, &model.normalFont),
			createBoldFontInput(cfg.Font.Bold.Family, &model.boldFont),
			createItalicFontInput(cfg.Font.Italic.Family, &model.italicFont),
			createBoldItalicFontInput(cfg.Font.BoldItalic.Family, &model.boldItalicFont),
			createFontSizeInput(cfg.Font.Size, &model.fontSize),
		),

		huh.NewGroup(
			huh.NewNote().Title("3. CURSOR"),
			createShapeSelect(cfg.Cursor.Style.Shape, &model.shape),
			createBlinkingSelect(cfg.Cursor.Style.BlinkIng, &model.blinking),
		),
	)

	model.form = form

	return model
}

func createDecorationsSelect(defaultValue string, value *string) *huh.Select[string] {
	return huh.NewSelect[string]().
		Title("1.1 Choose decorations").
		Options(
			huh.NewOption(optionSelected("Full", defaultValue)),
			huh.NewOption(optionSelected("None", defaultValue)),
			huh.NewOption(optionSelected("Transparent", defaultValue)),
			huh.NewOption(optionSelected("Buttonless", defaultValue)),
		).
		Value(value).
		DescriptionFunc(func() string {
			return descMap[*value]
		}, value)
}

func createStartupModeSelect(defaultValue string, value *string) *huh.Select[string] {
	return huh.NewSelect[string]().
		Title("1.2 Choose startup mode").
		Options(
			huh.NewOption(optionSelected("Windowed", defaultValue)),
			huh.NewOption(optionSelected("Maximized", defaultValue)),
			huh.NewOption(optionSelected("Fullscreen", defaultValue)),
			huh.NewOption(optionSelected("SimpleFullscreen", defaultValue)),
		).
		Value(value).
		DescriptionFunc(func() string {
			return descMap[*value]
		}, value)
}

func createTitleInput(defaultValue string, value *string) *huh.Input {
	return huh.NewInput().
		Title("1.3 window title").
		Placeholder(defaultValue).
		Value(value).
		DescriptionFunc(func() string {
			if *value == "" {
				return "Default: \"Alacritty\""
			}
			return *value
		}, value)
}

func createColumnsInput(defaultValue int, value *string) *huh.Input {
	return huh.NewInput().
		Title("1.4 Input columns").
		Value(value).
		Validate(validateIntF()).
		Placeholder(strconv.Itoa(defaultValue)).
		DescriptionFunc(func() string {
			if *value == "" {
				return "Default: 180"
			}
			return *value
		}, value)
}

func createLinesInput(defaultValue int, value *string) *huh.Input {
	return huh.NewInput().
		Title("1.5 Input lines").
		Value(value).
		Validate(validateIntF()).
		Placeholder(strconv.Itoa(defaultValue)).
		DescriptionFunc(func() string {
			if *value == "" {
				return "Default: 40"
			}
			return *value
		}, value)
}

func createOpacityInput(defaultValue float64, value *string) *huh.Input {
	return huh.NewInput().
		Title("1.6 Input opacity").
		Value(value).
		Validate(validateFloatF()).
		DescriptionFunc(func() string {
			if *value == "" {
				return "Default: 1.0"
			}
			return *value
		}, value).
		Placeholder(strconv.FormatFloat(defaultValue, 'f', -1, 64))
}

func createNormalFontInput(defaultValue string, value *string) *huh.Input {
	return huh.NewInput().
		Title("2.1 Input normal font").
		Value(value).
		DescriptionFunc(func() string {
			if *value == "" {
				return "Default: " + defaultFont()
			}
			return *value
		}, value).
		Placeholder(defaultValue)
}

func createBoldFontInput(defaultValue string, value *string) *huh.Input {
	return huh.NewInput().
		Title("2.2 Input bold font").
		Value(value).
		DescriptionFunc(func() string {
			if *value == "" {
				return "If the family is not specified, it will fall back to the value specified for the normal font"
			}
			return *value
		}, value).
		Placeholder(defaultValue)
}

func createItalicFontInput(defaultValue string, value *string) *huh.Input {
	return huh.NewInput().
		Title("2.3 Input italic font").
		Value(value).
		DescriptionFunc(func() string {
			if *value == "" {
				return "If the family is not specified, it will fall back to the value specified for the normal font"
			}
			return *value
		}, value).
		Placeholder(defaultValue)
}

func createBoldItalicFontInput(defaultValue string, value *string) *huh.Input {
	return huh.NewInput().
		Title("2.4 Input bold italic font").
		Value(value).
		DescriptionFunc(func() string {
			if *value == "" {
				return "If the family is not specified, it will fall back to the value specified for the normal font"
			}
			return *value
		}, value).
		Placeholder(defaultValue)
}

func createFontSizeInput(defaultValue float64, value *string) *huh.Input {
	return huh.NewInput().
		Title("2.5 Input font size").
		Value(value).
		Validate(validateFloatF()).
		DescriptionFunc(func() string {
			if *value == "" {
				return "Default: 11.25"
			}
			return *value
		}, value).
		Placeholder(strconv.FormatFloat(defaultValue, 'f', -1, 64))
}

func createShapeSelect(defaultValue string, value *string) *huh.Select[string] {
	return huh.NewSelect[string]().
		Title("3.1 Choose shape").
		Options(
			huh.NewOption(optionSelected("Block", defaultValue)),
			huh.NewOption(optionSelected("Underline", defaultValue)),
			huh.NewOption(optionSelected("Beam", defaultValue)),
		).
		Value(value).
		DescriptionFunc(func() string {
			return descMap[*value]
		}, value)
}

func createBlinkingSelect(defaultValue string, value *string) *huh.Select[string] {
	return huh.NewSelect[string]().
		Title("3.2 Choose blinking").
		Options(
			huh.NewOption(optionSelected("Never", defaultValue)),
			huh.NewOption(optionSelected("Off", defaultValue)),
			huh.NewOption(optionSelected("On", defaultValue)),
			huh.NewOption(optionSelected("Always", defaultValue)),
		).
		Value(value).
		DescriptionFunc(func() string {
			return descMap[*value]
		}, value)
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
				config.WriteConfig(combineCfg(m.oldCfgMap, m.newCfgMap()))
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
		return fmt.Sprintf(`1. WINDOW

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

Confirm?[Y/n]
`,
			resultColor(m.oldCfg.Window.Decorations, m.decorations),
			resultColor(m.oldCfg.Window.StartMode, m.startupMode),
			resultColor(m.oldCfg.Window.Title, m.title),
			resultColor(strconv.Itoa(m.oldCfg.Window.Dimensions.Columns), m.column),
			resultColor(strconv.Itoa(m.oldCfg.Window.Dimensions.Lines), m.line),
			resultColor(strconv.FormatFloat(m.oldCfg.Window.Opacity, 'f', -1, 64), m.opacity),

			resultColor(m.oldCfg.Font.Normal.Family, m.normalFont),
			resultColor(m.oldCfg.Font.Bold.Family, m.boldFont),
			resultColor(m.oldCfg.Font.Italic.Family, m.italicFont),
			resultColor(m.oldCfg.Font.BoldItalic.Family, m.boldItalicFont),
			resultColor(strconv.FormatFloat(m.oldCfg.Font.Size, 'f', -1, 64), m.fontSize),

			resultColor(m.oldCfg.Cursor.Style.Shape, m.shape),
			resultColor(m.oldCfg.Cursor.Style.BlinkIng, m.blinking),
		)
	}
	return m.form.View()
}

func defaultFont() string {
	switch runtime.GOOS {
	case "darwin":
		return "Menlo"
	case "linux":
		return "monospace"
	case "windows":
		return "Consolas"
	default:
		return ""
	}
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

func mustCovert[T ~int | ~float64](s string, defaultValue T) T {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}

	var v T
	switch any(v).(type) {
	case int:
		return T(int(f))
	default:
		return T(f)
	}
}

func defaultValue(s, v string) string {
	if s == "" {
		return v
	}
	return s
}

func (m ConfigModel) newCfgMap() map[string]any {
	window := map[string]any{
		"decorations":  m.decorations,
		"startup_mode": m.startupMode,
		"title":        defaultValue(m.title, "Alacritty"),
		"dimensions": map[string]any{
			"columns": mustCovert(m.column, defaultColumns),
			"lines":   mustCovert(m.line, defaultLines),
		},
		"opacity": mustCovert(m.opacity, defaultOpacity),
	}
	f := defaultValue(m.normalFont, defaultFont())
	font := map[string]any{
		"size": mustCovert(m.fontSize, defaultFontSize),
		"normal": map[string]any{
			"family": f,
		},
		"bold": map[string]any{
			"family": defaultValue(m.boldFont, f),
		},
		"italic": map[string]any{
			"family": defaultValue(m.italicFont, f),
		},
		"bold_italic": map[string]any{
			"family": defaultValue(m.boldItalicFont, f),
		},
	}

	cursor := map[string]any{
		"style": map[string]any{
			"shape":    m.shape,
			"blinking": m.blinking,
		},
	}

	return map[string]any{
		"window": window,
		"font":   font,
		"cursor": cursor,
	}
}

func combineCfg(oldCfg, newCfg map[string]any) map[string]any {
	for k := range oldCfg {
		c, ok := newCfg[k]
		if ok {
			oldCfg[k] = c
		}
	}
	return oldCfg
}
