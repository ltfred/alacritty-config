package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ltfred/alacritty/pkg/ui"
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
	if _, err := tea.NewProgram(ui.NewConfigModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
