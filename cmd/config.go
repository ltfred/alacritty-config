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
	Short: "Launch an interactive configuration interface for Alacritty.",
	Run:   config,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func config(cmd *cobra.Command, args []string) {
	if _, err := tea.NewProgram(ui.NewConfigModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
