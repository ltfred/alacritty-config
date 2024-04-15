package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ltfred/alacritty-config/themes"
	"github.com/pelletier/go-toml/v2"

	"github.com/gookit/color"
	"github.com/ltfred/alacritty-config/pkg/config"
	"github.com/ltfred/alacritty-config/pkg/prompt"
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
	go func() {
		defer func() {
			prompt.Quit <- true
		}()
		// confirm
		confirmChoices := []string{"Yes", "No"}
		newSelect := prompt.NewSelect(confirmChoices, "This is "+color.Blue.Sprintf(
			"alacritty configuration wizard.")+
			"It will ask you a few questions and configure your alacritty, "+
			"which will cover your previous configuration, continue?\n\n")
		sel, err := newSelect.Run()
		if err != nil || sel == "No" {
			return
		}
		cfg := &config.Config{}
		cfg.SetDefault()
		// window
		setWindowCfg(cfg)
		// font
		setFontCfg(cfg)
		// theme
		setThemeCfg(cfg)
		// cursor
		setCursorCfg(cfg)

		color.Infof("%v", cfg)

	}()
	// quit channel
	var quit bool
	select {
	case quit = <-prompt.Quit:
		if quit {
			return
		}
	}
}

func setWindowCfg(cfg *config.Config) {
	// window title
	newInput := prompt.NewInput("\nInput window title:", "Alacritty")
	title, err := newInput.Run()
	if err != nil {
		return
	}
	if title != "" {
		cfg.Window.Title = title
	}
	// window size
	inputs := prompt.NewInputs([]prompt.InputsOptions{
		{
			Label:       "Columns",
			Placeholder: "180",
			CharLimit:   3,
			Validate: func(s string) error {
				if _, err = strconv.ParseInt(s, 10, 64); err != nil {
					color.Warn.Prompt("Columns must be a number")
					return err
				}
				return nil
			},
		},
		{
			Label:       "Lines",
			Placeholder: "50",
			CharLimit:   3,
			Validate: func(s string) error {
				if _, err = strconv.ParseInt(s, 10, 64); err != nil {
					color.Warn.Prompt("Lines must be a number")
					return err
				}
				return nil
			},
		},
	}, "Input window size:")
	window, err := inputs.Run()
	if err != nil {
		return
	}
	cfg.Window.Dimensions.Columns, cfg.Window.Dimensions.Lines = 180, 50
	if window[0] != "" {
		columns, _ := strconv.ParseInt(window[0], 10, 64)
		cfg.Window.Dimensions.Columns = int(columns)
	}
	if window[1] != "" {
		rows, _ := strconv.ParseInt(window[1], 10, 64)
		cfg.Window.Dimensions.Lines = int(rows)
	}
	// window decorations
	decorations := []string{
		fmt.Sprintf("%s--Borders and title bar", color.Blue.Sprint("Full")),
		fmt.Sprintf("%s--Neither borders nor title bar", color.Blue.Sprint("None")),
		fmt.Sprintf("%s--Title bar, transparent background and title bar buttons(macOS only)", color.Blue.Sprint("Transparent")),
		fmt.Sprintf("%s--Title bar, transparent background and no title bar buttons(macOS only)", color.Blue.Sprint("Buttonless")),
	}
	newSelect := prompt.NewSelect(decorations, "5. Choose window decorations:\n\n")
	sel, err := newSelect.Run()
	if err != nil {
		return
	}
	cfg.Window.Decorations = color.ClearCode(strings.Split(sel, "--")[0])

	startupModes := []string{
		fmt.Sprintf("%s--Regular window", color.Blue.Sprint("Windowed")),
		fmt.Sprintf("%s--The window will be maximized on startup", color.Blue.Sprint("Maximized")),
		fmt.Sprintf("%s--The window will be fullscreened on startup", color.Blue.Sprint("Fullscreen")),
		fmt.Sprintf("%s--Same as Fullscreen, but you can stack windows on top", color.Blue.Sprint("SimpleFullscreen")),
	}
	newSelect = prompt.NewSelect(startupModes, "\nChoose Startup mode (changes require restart):\n\n")
	sel, err = newSelect.Run()
	if err != nil {
		return
	}
	cfg.Window.StartupMode = color.ClearCode(strings.Split(sel, "--")[0])
}

func setFontCfg(cfg *config.Config) {
	// font
	newInput := prompt.NewInput("\nInput font(recommend nerd font):", "")
	font, err := newInput.Run()
	if err != nil {
		return
	}
	cfg.Font.Normal.Family = font
	// font size
	newInput = prompt.NewInput("Input font size:", "15.0", func(s string) error {
		if _, err = strconv.ParseFloat(s, 32); err != nil {
			color.Warn.Prompt("Please enter a number")
			return err
		}
		return nil
	})
	fontSize, err := newInput.Run()
	if err != nil {
		return
	}
	size, _ := strconv.ParseFloat(fontSize, 32)
	if size != 0 {
		cfg.Font.Size = float32(size)
	}
}

func setThemeCfg(cfg *config.Config) {
	confirmChoices := []string{"Yes", "No"}
	newSelect := prompt.NewSelect(confirmChoices, "Choose a theme:\n\n")
	sel, err := newSelect.Run()
	if err != nil {
		return
	}
	if sel == "Yes" {
		i, err := displayTheme()
		if err != nil {
			color.Error.Prompt(err.Error())
			return
		}
		th := themes.Themes[i]
		var colors config.Config
		if err = toml.Unmarshal(themes.ThemesMap[th], &colors); err != nil {
			color.Error.Prompt(err.Error())
			return
		}
		cfg.Colors = colors.Colors
	}
}

func setCursorCfg(cfg *config.Config) {
	newSelect := prompt.NewSelect([]string{"Block", "Underline", "Beam"}, "\nChoose cursor style:\n\n")
	cursor, err := newSelect.Run()
	if err != nil {
		return
	}
	if cursor != "" {
		cfg.Cursor.Style.Shape = cursor
	}
	blinkings := []string{
		fmt.Sprintf("%s--Prevent the cursor from ever blinking", color.Blue.Sprint("Never")),
		fmt.Sprintf("%s--Disable blinking by default", color.Blue.Sprint("Off")),
		fmt.Sprintf("%s--Enable blinking by default", color.Blue.Sprint("Never")),
		fmt.Sprintf("%s--Force the cursor to always blink", color.Blue.Sprint("Always")),
	}
	newSelect = prompt.NewSelect(blinkings, "\n Choose cursor blinking:\n\n")
	blinking, err := newSelect.Run()
	if err != nil {
		return
	}
	if cursor != "" {
		cfg.Cursor.Style.Blinking = color.ClearCode(strings.Split(blinking, "--")[0])

	}
	if cfg.Cursor.Style.Blinking == "On" || cfg.Cursor.Style.Blinking == "Always" {
		input := prompt.NewInput("\nInput Cursor blinking interval(milliseconds):", "750",
			func(s string) error {
				_, err = strconv.ParseInt(s, 10, 64)
				if err != nil {
					color.Warn.Prompt("Please enter a number")
					return err
				}
				return nil
			})
		interval, err := input.Run()
		if err != nil {
			return
		}
		if interval != "" {
			inter, _ := strconv.ParseInt(interval, 10, 64)
			cfg.Cursor.BlinkInterval = int(inter)
		}
	}
}
