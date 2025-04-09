package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ltfred/alacritty-config/themes"
	"github.com/pelletier/go-toml/v2"

	"github.com/manifoldco/promptui"

	"github.com/gookit/color"
	"github.com/ltfred/alacritty-config/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init alacritty configuration",
	Run:   initConfig,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initConfig(cmd *cobra.Command, args []string) {
	if !confirm("This is alacritty configuration wizard, " +
		"It will ask you a few questions and configure your alacritty, which will cover your previous configuration, " +
		"continue") {
		os.Exit(1)
	}

	cfg := &config.Config{}
	cfg.SetDefault()
	// window
	setWindowCfg(cfg)
	// font
	setFontCfg(cfg)
	//// theme
	setThemeCfg(cfg)
	//// cursor
	setCursorCfg(cfg)

	color.Infof("%v", cfg)

	cfg.WriteConfig()
}

func setWindowCfg(cfg *config.Config) {
	// window title
	if title := input("Input window title", "Alacritty"); title != "" {
		cfg.Window.Title = title
	}
	// window size
	if s := input("Input window columns", "180", func(s string) error {
		if s != "" {
			if _, err := strconv.ParseUint(s, 10, 64); err != nil {
				return fmt.Errorf("columns must be a number")
			}
		}
		return nil
	}); s != "" {
		columns, _ := strconv.ParseUint(s, 10, 64)
		cfg.Window.Dimensions.Columns = int(columns)
	}
	if s := input("Input window lines", "50", func(s string) error {
		if s != "" {
			if _, err := strconv.ParseUint(s, 10, 64); err != nil {
				return fmt.Errorf("lines must be a number")
			}
		}
		return nil
	}); s != "" {
		lines, _ := strconv.ParseUint(s, 10, 64)
		cfg.Window.Dimensions.Lines = int(lines)
	}

	// window decorations
	items := []Option{
		{Name: "Full", Desc: "Borders and title bar"},
		{Name: "None", Desc: "Neither borders nor title bar"},
		{Name: "Transparent", Desc: "Title bar, transparent background and title bar buttons(macOS only)"},
		{Name: "Buttonless", Desc: "Title bar, transparent background and no title bar buttons(macOS only)"},
	}
	i := singleSelect("Choose window decorations", items)
	cfg.Window.Decorations = items[i].Name

	// opacity
	if s := input("Input window opacity", "1.0", func(s string) error {
		if s != "" {
			opacity, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fmt.Errorf("opacity must be a floating point number from 0.0 to 1.0")
			}
			if opacity < 0 || opacity > 1 {
				return fmt.Errorf("opacity must be a floating point number from 0.0 to 1.0")
			}
		}
		return nil
	}); s != "" {
		opacity, _ := strconv.ParseFloat(s, 64)
		cfg.Window.Opacity = float32(opacity)
	}
}

func setFontCfg(cfg *config.Config) {
	if s := input("Input font(recommend nerd font)", "Menlo"); s != "" {
		cfg.Font.Normal.Family = s
	}
}

func setThemeCfg(cfg *config.Config) {
	if !confirm("Choose a theme") {
		return
	}
	i, err := displayTheme()
	if err != nil {
		os.Exit(1)
	}
	th := themes.Themes[i]
	var colors config.Config
	if err = toml.Unmarshal(themes.ThemesMap[th], &colors); err != nil {
		os.Exit(1)
	}
	cfg.Colors = colors.Colors
}

func setCursorCfg(cfg *config.Config) {
	cursor := []Option{
		{Name: "Block", Desc: "\u2588"},
		{Name: "Underline", Desc: "\u2581"},
		{Name: "Beam", Desc: "\u258F"},
	}

	i := singleSelect("Choose cursor blinking", cursor)
	cfg.Cursor.Style.Shape = cursor[i].Name

	items := []Option{
		{Name: "Never", Desc: "Prevent the cursor from ever blinking"},
		{Name: "Off", Desc: "Disable blinking by default"},
		{Name: "On", Desc: "Enable blinking by default"},
		{Name: "Always", Desc: "Force the cursor to always blink"},
	}

	i = singleSelect("Choose cursor blinking", items)
	cfg.Cursor.Style.Blinking = items[i].Name
}

func input(label, d string, validate ...promptui.ValidateFunc) string {
	t := promptui.Prompt{Label: label, Default: d}
	if len(validate) > 0 {
		t.Validate = validate[0]
	}
	s, err := t.Run()
	if err != nil {
		os.Exit(1)
	}
	return s
}

func confirm(label string) bool {
	confirmChoices := []string{"Yes", "No"}
	newSelect := promptui.Select{Label: label, Items: confirmChoices}
	_, sel, err := newSelect.Run()
	if err != nil {
		os.Exit(1)
	}
	return sel == "Yes"
}

func singleSelect(label string, items any) int {
	sel := promptui.Select{Label: label, Items: items, Templates: template}
	i, _, err := sel.Run()
	if err != nil {
		os.Exit(1)
	}
	return i
}

type Option struct {
	Name string
	Desc string
}

var template = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   fmt.Sprintf("%s {{ .Name | underline }}", promptui.IconSelect),
	Inactive: "  {{.Name}}",
	Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Name | faint }}`, promptui.IconGood),
	Details:  `Desc: {{.Desc}}`,
}
