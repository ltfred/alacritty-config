package cmd

import (
	"github.com/ltfred/alacritty/pkg"
	"github.com/spf13/cobra"
)

// themesCmd represents the themes command
var themesCmd = &cobra.Command{
	Use:   "themes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: themes,
}

func init() {
	rootCmd.AddCommand(themesCmd)
}

func themes(cmd *cobra.Command, args []string) {
	pkg.FuzzyThemes()
}
