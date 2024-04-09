package cmd

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/ltfred/alacritty-config/pkg/config"
	"github.com/ltfred/alacritty-config/pkg/prompt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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
		sel := prompt.Select{
			Choices: []string{"Yes", "No"},
			Label: "This is " + color.New(color.FgBlue).Sprint(
				"alacritty configuration wizard.") +
				"It will ask you a few questions and configure your alacritty, continue?\n\n",
		}
		p := tea.NewProgram(sel)
		m, err := p.Run()
		if err != nil {
			color.Red(err.Error())
			return
		}
		if m, ok := m.(prompt.Select); ok && m.Result != "" {
			if m.Result == "No" {
				return
			}
		}

		cfg := config.Config{}

		// font
		fontInput := textinput.New()
		fontInput.Focus()
		fontInput.CharLimit, fontInput.Width = 156, 30
		input := prompt.Input{TextInput: fontInput, Label: "1. Input font(recommend nerd font):"}
		p = tea.NewProgram(input)
		m, err = p.Run()
		if err != nil {
			color.Red(err.Error())
			return
		}
		if m, ok := m.(prompt.Input); ok {
			cfg.Font.Normal.Family = m.TextInput.Value()
		}

		// font size
		fontSizeInput := textinput.New()
		fontSizeInput.Focus()
		fontSizeInput.CharLimit, fontInput.Width = 156, 30
		fontSizeInput.Validate = func(s string) error {
			_, err := strconv.ParseFloat(s, 32)
			return err
		}
		input = prompt.Input{TextInput: fontSizeInput, Label: "2. Input font size:"}
		p = tea.NewProgram(input)
		m, err = p.Run()
		if err != nil {
			color.Red(err.Error())
			return
		}
		if m, ok := m.(prompt.Input); ok {
			size, _ := strconv.ParseFloat(m.TextInput.Value(), 32)
			cfg.Font.Size = float32(size)
		}

		// window
	}()

	var quit bool
	select {
	case quit = <-prompt.Quit:
		if quit {
			return
		}
	}
}
