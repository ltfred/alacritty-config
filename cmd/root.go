package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/color"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/ltfred/alacritty-config/pkg/config"
	"github.com/ltfred/alacritty-config/pkg/prompt"
	"github.com/ltfred/alacritty-config/themes"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alacritty-config",
	Short: "Init alacritty config",
	Run:   initConfig,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig(cmd *cobra.Command, args []string) {
	go func() {
		defer func() {
			prompt.Quit <- true
		}()
		// confirm
		isConfirm, err := confirm()
		if err != nil || !isConfirm {
			return
		}

		cfg := config.Config{}

		// font
		fontFamily, err := font()
		if err != nil {
			return
		}
		cfg.Font.Normal.Family = fontFamily

		// font size
		fontSize, err := fonSize()
		if err != nil {
			return
		}
		cfg.Font.Size = fontSize

		// theme
		i, err := displayTheme()
		if err != nil {
			color.Error.Prompt(err.Error())
			return
		}
		th := themes.Themes[i]
		var colors config.Config
		err = toml.Unmarshal(themes.ThemesMap[th], &colors)
		if err != nil {
			color.Error.Prompt(err.Error())
			return
		}
		cfg.Colors = colors.Colors

	}()

	var quit bool
	select {
	case quit = <-prompt.Quit:
		if quit {
			return
		}
	}
}

func confirm() (bool, error) {
	sel := prompt.Select{
		Choices: []string{"Yes", "No"},
		Label: "This is " + color.Blue.Sprintf(
			"alacritty configuration wizard.") +
			"It will ask you a few questions and configure your alacritty, " +
			"which will cover your previous configuration, continue?\n\n",
	}
	model, err := runPrompt(sel)
	if err != nil {
		return false, err
	}
	if m, ok := model.(prompt.Select); ok && m.Result != "" {
		if m.Result == "No" {
			return false, nil
		}
	}

	return true, nil
}

func font() (string, error) {
	fontInput := textinput.New()
	fontInput.Focus()
	fontInput.CharLimit, fontInput.Width = 156, 30
	input := prompt.Input{TextInput: fontInput, Label: "1. Input font(recommend nerd font):"}
	model, err := runPrompt(input)
	if err != nil {
		return "", err
	}
	if m, ok := model.(prompt.Input); ok {
		return m.TextInput.Value(), nil
	}

	return "", nil
}

func fonSize() (float32, error) {
	fontSizeInput := textinput.New()
	fontSizeInput.Focus()
	fontSizeInput.CharLimit, fontSizeInput.Width = 156, 30
	fontSizeInput.Validate = func(s string) error {
		_, err := strconv.ParseFloat(s, 32)
		return err
	}
	input := prompt.Input{TextInput: fontSizeInput, Label: "2. Input font size:"}
	model, err := runPrompt(input)
	if err != nil {
		return 0, err
	}
	if m, ok := model.(prompt.Input); ok {
		size, _ := strconv.ParseFloat(m.TextInput.Value(), 32)
		return float32(size), nil
	}

	return 0, nil
}

func runPrompt(model tea.Model) (tea.Model, error) {
	p := tea.NewProgram(model)
	m, err := p.Run()
	if err != nil {
		color.Error.Prompt(err.Error())
		return nil, err
	}
	return m, nil
}

func colorPrint(col string) string {
	if len(col) > 2 {
		col = col[1:]
	}

	return color.Sprintf("<bg=%s>llo</>", col)
}

func displayTheme() (int, error) {
	idx, err := fuzzyfinder.FindMulti(
		themes.Themes,
		func(i int) string {
			return themes.Themes[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			th := themes.Themes[i]
			var colors config.Config
			err := toml.Unmarshal(themes.ThemesMap[th], &colors)
			if err != nil {
				color.Error.Prompt(err.Error())
				return ""
			}

			return fmt.Sprintf("        Normal Bright\n  Black %s     %s\n    Red %s     %s\n  Green %s     %s\n"+
				" Yellow %s     %s\n   Blue %s     %s\nMagenta %s     %s\n   Cyan %s     %s\n  White %s     %s",
				colorPrint(colors.Colors.Normal.Black), colorPrint(colors.Colors.Bright.Black),
				colorPrint(colors.Colors.Normal.Red), colorPrint(colors.Colors.Bright.Red),
				colorPrint(colors.Colors.Normal.Green), colorPrint(colors.Colors.Bright.Green),
				colorPrint(colors.Colors.Normal.Yellow), colorPrint(colors.Colors.Bright.Yellow),
				colorPrint(colors.Colors.Normal.Blue), colorPrint(colors.Colors.Bright.Blue),
				colorPrint(colors.Colors.Normal.Magenta), colorPrint(colors.Colors.Bright.Magenta),
				colorPrint(colors.Colors.Normal.Cyan), colorPrint(colors.Colors.Bright.Cyan),
				colorPrint(colors.Colors.Normal.White), colorPrint(colors.Colors.Bright.White),
			)
		}),
	)
	if err != nil {
		return 0, err
	}

	return idx[0], nil
}
