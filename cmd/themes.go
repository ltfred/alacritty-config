package cmd

import (
	"github.com/ltfred/alacritty/pkg"
	"github.com/spf13/cobra"
)

// themesCmd represents the themes command
var themesCmd = &cobra.Command{
	Use:   "themes",
	Short: "The `themes` command is used to preview or list all the available themes for Alacritty.",
	Run:   themes,
}

func init() {
	rootCmd.AddCommand(themesCmd)
}

func themes(cmd *cobra.Command, args []string) {
	pkg.FuzzyThemes()
}
