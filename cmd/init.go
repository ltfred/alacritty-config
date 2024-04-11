package cmd

import (
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/ltfred/alacritty-config/pkg/config"
	"github.com/ltfred/alacritty-config/pkg/prompt"
	"github.com/ltfred/alacritty-config/themes"
	"github.com/pelletier/go-toml/v2"

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
		if err != nil || sel == "" || sel == "No" {
			return
		}

		cfg := config.Config{}
		// font
		newInput := prompt.NewInput("\n1. Input font(recommend nerd font):", "")
		cfg.Font.Normal.Family, err = newInput.Run()
		if err != nil {
			return
		}
		cfg.Font.Normal.Style = "Regular"
		// font size
		cfg.Font.Size = 15.0
		newInput = prompt.NewInput("2. Input font size:", "15.0", func(s string) error {
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

		// theme
		newSelect = prompt.NewSelect(confirmChoices, "3. Choose a theme:\n\n")
		sel, err = newSelect.Run()
		if err != nil || sel == "" {
			return
		}
		if sel == "Yes" {
			if i, err := displayTheme(); err != nil {
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
		// window title
		cfg.Window.Title = "Alacritty"
		newInput = prompt.NewInput("\n3. Input window title:", "")
		title, err := newInput.Run()
		if err != nil {
			return
		}
		if title != "" {
			cfg.Window.Title = title
		}
		// window size
		newInput = prompt.NewInput("4. Input window size:", "180 * 50", func(s string) error {
			split := strings.Split(s, " * ")
			if len(split) != 2 {
				color.Warn.Prompt("The format should be like 180 * 50")
				return err
			}
			if _, err = strconv.ParseInt(split[0], 10, 64); err != nil {
				color.Warn.Prompt("The format should be like 180 * 50")
				return err
			}
			if _, err = strconv.ParseInt(split[1], 10, 64); err != nil {
				color.Warn.Prompt("The format should be like 180 * 50")
				return err
			}
			return nil
		})
		window, err := newInput.Run()
		if err != nil {
			return
		}
		if window != "" {
			split := strings.Split(window, " * ")
			columns, _ := strconv.ParseInt(split[0], 10, 64)
			rows, _ := strconv.ParseInt(split[1], 10, 64)
			cfg.Window.Dimensions.Columns, cfg.Window.Dimensions.Lines = int(columns), int(rows)
		} else {
			cfg.Window.Dimensions.Columns, cfg.Window.Dimensions.Lines = 180, 50
		}
		// window decorations
		newSelect = prompt.NewSelect([]string{"Full", "None", "Transparent", "Buttonless"},
			"5. Choose window decorations:\n\n")
		sel, err = newSelect.Run()
		if err != nil || sel == "" {
			return
		}
		cfg.Window.Decorations = sel

		// cursor

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
