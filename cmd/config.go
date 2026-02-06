package cmd

import (
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: config,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func config(cmd *cobra.Command, args []string) {
	var (
		decorations string
		startupMode string
		columns     string
		lines       string
		opacity     string

		font     string
		fontSize string
	)

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose decorations").
				Options(
					huh.NewOption("Full", "Full"),
					huh.NewOption("None", "None"),
					huh.NewOption("Transparent(macOS only)", "Transparent"),
					huh.NewOption("Buttonless(macOS only)", "Buttonless"),
				).Value(&decorations),

			huh.NewSelect[string]().
				Title("Choose startup_mode").
				Options(
					huh.NewOption("Windowed", "Windowed"),
					huh.NewOption("Maximized", "Maximized"),
					huh.NewOption("Fullscreen", "Fullscreen"),
					huh.NewOption("SimpleFullscreen(macOS only)", "SimpleFullscreen"),
				).Value(&startupMode),
		),

		huh.NewGroup(
			huh.NewInput().Title("Input columns").Placeholder("180").Value(&columns).Validate(func(s string) error {
				if s != "" {
					_, err := strconv.Atoi(s)
					return err
				}
				return nil
			}),
			huh.NewInput().Title("Input lines").Placeholder("50").Value(&lines).Validate(func(s string) error {
				if s != "" {
					_, err := strconv.Atoi(s)
					return err
				}
				return nil
			}),
			huh.NewInput().Title("Input opacity").Placeholder("1.0").Value(&opacity).Validate(func(s string) error {
				if s != "" {
					_, err := strconv.ParseFloat(s, 64)
					return err
				}
				return nil
			}),
		),

		huh.NewGroup(
			huh.NewInput().Title("Input font").Value(&font),
			huh.NewInput().Title("Input font size").Value(&fontSize),
		),
	)

	if err := f.Run(); err != nil {
		os.Exit(1)
	}
}
